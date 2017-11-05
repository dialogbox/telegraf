package protobuf

import (
	"fmt"

	"github.com/influxdata/telegraf"
)

// Serializer encodes metrics into Protocol Buffer format
type Serializer struct {
	PrependLength bool
}

// Serialize implements serializers.Serializer.Serialize
// github.com/influxdata/telegraf/plugins/serializers/Serializer
func (s *Serializer) Serialize(m telegraf.Metric) ([]byte, error) {
	if s.PrependLength {
		return MarshalMetricLengthPrefix(m)
	}

	return MarshalMetric(m)
}

// SerializeBatch implements serializers.Serializer.SerializeBatch
// github.com/influxdata/telegraf/plugins/serializers/Serializer
func (s *Serializer) SerializeBatch(metrics []telegraf.Metric) ([]byte, error) {
	if s.PrependLength == false {
		return nil, fmt.Errorf("PrependLength must be enabled to use batch serialization")
	}

	return MarshalMetrics(metrics)
}
