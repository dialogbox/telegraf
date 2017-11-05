package protobuf

import (
	"testing"

	"github.com/influxdata/telegraf/testutil"
	"github.com/stretchr/testify/assert"
)

func TestMarshalMetricFloat(t *testing.T) {
	m := testutil.TestMetric(float64(91.5))

	buf, err := MarshalMetric(m)
	assert.NoError(t, err)

	tm, err := UnmarshalMetric(buf)
	assert.NoError(t, err)

	testutil.RequireMetricEqual(t, m, tm)
}

func TestMarshalMetricInt(t *testing.T) {
	m := testutil.TestMetric(int64(90))

	buf, err := MarshalMetric(m)
	assert.NoError(t, err)

	tm, err := UnmarshalMetric(buf)
	assert.NoError(t, err)

	testutil.RequireMetricEqual(t, m, tm)
}

func TestMarshalMetricString(t *testing.T) {
	m := testutil.TestMetric("foobar")

	buf, err := MarshalMetric(m)
	assert.NoError(t, err)

	tm, err := UnmarshalMetric(buf)
	assert.NoError(t, err)

	testutil.RequireMetricEqual(t, m, tm)
}

func TestMarshalMultiFields(t *testing.T) {
	m := testutil.TestMetric("foobar")

	m.AddField("int_field", int64(90))
	m.AddField("float_field", float64(8559615))
	m.AddField("string_field", "string_value")
	m.AddField("bool_field", true)

	buf, err := MarshalMetric(m)
	assert.NoError(t, err)

	tm, err := UnmarshalMetric(buf)
	assert.NoError(t, err)

	testutil.RequireMetricEqual(t, m, tm)
}

func TestMarshalMetricWithEscapes(t *testing.T) {
	m := testutil.TestMetric("foobar")

	m.AddField("U,age=Idle", int64(90))

	buf, err := MarshalMetric(m)
	assert.NoError(t, err)

	tm, err := UnmarshalMetric(buf)
	assert.NoError(t, err)

	testutil.RequireMetricEqual(t, m, tm)
}
