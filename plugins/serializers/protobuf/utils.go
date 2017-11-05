package protobuf

import (
	"bytes"
	"encoding/binary"
	"fmt"

	"github.com/influxdata/telegraf"
	"github.com/influxdata/telegraf/metric"
	"google.golang.org/protobuf/proto"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
)

// MarshalMetric marshals a metric into protobuf binary
func MarshalMetric(tm telegraf.Metric) ([]byte, error) {
	pm, err := telegrafMetricToProtoMetric(tm)
	if err != nil {
		return nil, err
	}

	result, err := proto.Marshal(pm)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// MarshalMetricLengthPrefix marshals a metric into protobuf binary with length prefix
func MarshalMetricLengthPrefix(tm telegraf.Metric) ([]byte, error) {
	protoMessage, err := MarshalMetric(tm)
	if err != nil {
		return nil, err
	}

	buf := &bytes.Buffer{}
	binary.Write(buf, binary.LittleEndian, uint32(len(protoMessage)))
	buf.Write(protoMessage)

	return buf.Bytes(), nil
}

// MarshalMetrics marshals metrics into a stream of length prefixed protobuf binary
func MarshalMetrics(tms []telegraf.Metric) ([]byte, error) {
	buf := &bytes.Buffer{}

	for i := range tms {
		protoMessage, err := MarshalMetricLengthPrefix(tms[i])
		if err != nil {
			return nil, err
		}

		buf.Write(protoMessage)
	}

	return buf.Bytes(), nil
}

// UnmarshalMetric unmarshals a protobuf binary data of a metric
func UnmarshalMetric(buf []byte) (telegraf.Metric, error) {
	pm := &Metric{}

	err := proto.Unmarshal(buf, pm)
	if err != nil {
		return nil, err
	}

	tm, err := protoMetricToTelegrafMetric(pm)
	if err != nil {
		return nil, err
	}

	return tm, nil
}

// UnmarshalMetrics unmarshals a length prefixed protobuf binary stream
func UnmarshalMetrics(buf []byte) ([]telegraf.Metric, error) {
	tms := make([]telegraf.Metric, 0)

	var total uint32 = uint32(len(buf))
	var cur uint32
	var msgLen uint32

	for cur < total {
		if cur+4 > total {
			return nil, fmt.Errorf("Unexpected end of input data at the index(%d)", cur)
		}

		msgLen = binary.LittleEndian.Uint32(buf[cur : cur+4])
		if cur+msgLen > total {
			return nil, fmt.Errorf("Number of bytes is not matched to the prefixed length at the index(%d): %d > %d", cur, msgLen, total-cur-4)
		}
		cur += 4 // advancing 32 bits

		tm, err := UnmarshalMetric(buf[cur : cur+msgLen])
		if err != nil {
			return nil, err
		}
		cur += msgLen

		tms = append(tms, tm)
	}

	return tms, nil
}

// protoToTelegrafMetric converts a telegraf.Metric to a protobuf structure
func telegrafMetricToProtoMetric(tm telegraf.Metric) (*Metric, error) {
	fields := make(map[string]*FieldValue)

	for k, v := range tm.Fields() {
		switch v := v.(type) {
		case string:
			fields[k] = &FieldValue{Value: &FieldValue_StringValue{StringValue: v}}
		case int64:
			fields[k] = &FieldValue{Value: &FieldValue_IntValue{IntValue: v}}
		case float64:
			fields[k] = &FieldValue{Value: &FieldValue_FloatValue{FloatValue: v}}
		case bool:
			fields[k] = &FieldValue{Value: &FieldValue_BoolValue{BoolValue: v}}
		default:
			return nil, fmt.Errorf("Unsupported field value data type: %T", v)
		}
	}

	return &Metric{
		Name:      tm.Name(),
		Timestamp: timestamppb.New(tm.Time()),
		Tags:      tm.Tags(),
		Fields:    fields,
	}, nil
}

// protoMetricToTelegrafMetric converts a protobuf structure to a telegraf.Metric
func protoMetricToTelegrafMetric(m *Metric) (telegraf.Metric, error) {
	fields := make(map[string]interface{})

	for k, v := range m.Fields {
		switch v := v.Value.(type) {
		case *FieldValue_BoolValue:
			fields[k] = v.BoolValue
		case *FieldValue_IntValue:
			fields[k] = v.IntValue
		case *FieldValue_StringValue:
			fields[k] = v.StringValue
		case *FieldValue_FloatValue:
			fields[k] = v.FloatValue
		default:
			return nil, fmt.Errorf("Unsupported field value type %T", v)
		}
	}

	metric, err := metric.New(m.Name, m.Tags, fields, m.Timestamp.AsTime())
	if err != nil {
		return nil, err
	}

	return metric, nil
}
