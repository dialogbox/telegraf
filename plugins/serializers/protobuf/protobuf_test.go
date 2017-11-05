package protobuf

import (
	"testing"

	"github.com/influxdata/telegraf"
	"github.com/influxdata/telegraf/testutil"

	"github.com/stretchr/testify/assert"
)

func TestSerialize(t *testing.T) {
	m := testutil.TestMetric("foobar")

	m.AddField("int_field", int64(90))
	m.AddField("float_field", float64(8559615))
	m.AddField("string_field", "string_value")
	m.AddField("bool_field", true)

	s := Serializer{}
	buf, err := s.Serialize(m)
	assert.NoError(t, err)

	tm, err := UnmarshalMetric(buf)
	assert.NoError(t, err)

	testutil.RequireMetricEqual(t, m, tm)
}

func TestSerializeMultipleMetric(t *testing.T) {
	m := testutil.TestMetric(int(90))

	s := Serializer{PrependLength: true}

	encoded, err := s.Serialize(m)
	assert.NoError(t, err)

	// Multiple metrics in continous bytes stream
	var buf []byte
	buf = append(buf, encoded...)
	buf = append(buf, encoded...)
	buf = append(buf, encoded...)
	buf = append(buf, encoded...)

	ms, err := UnmarshalMetrics(buf)
	assert.NoError(t, err)

	for i := range ms {
		testutil.RequireMetricEqual(t, m, ms[i])
	}
}

func TestSerializeBatch(t *testing.T) {
	m := testutil.TestMetric(int(90))

	metrics := []telegraf.Metric{m, m, m, m}

	s := Serializer{PrependLength: true}
	buf, err := s.SerializeBatch(metrics)
	assert.NoError(t, err)

	ms, err := UnmarshalMetrics(buf)
	assert.NoError(t, err)

	for i := range ms {
		testutil.RequireMetricEqual(t, m, ms[i])
	}
}
