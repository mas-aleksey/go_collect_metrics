package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewCounterJSONMetric(t *testing.T) {
	metric := NewCounterJSONMetric("name", 123)
	assert.Equal(t, "counter", metric.MType)
	assert.Equal(t, "name", metric.ID)
	assert.Equal(t, int64(123), *metric.Delta)

	metric.Hash = CalcHash(metric.String(), "foo")
	assert.Equal(t, "fe55643dfd050289512090577d279326e6e27f23f82b6b2e61232e965b8e6d5a", *metric.Hash)
	assert.Nil(t, metric.Value)
}

func TestNewGaugeJSONMetric(t *testing.T) {
	metric := NewGaugeJSONMetric("name", 123.4)
	assert.Equal(t, "gauge", metric.MType)
	assert.Equal(t, "name", metric.ID)
	assert.Equal(t, 123.4, *metric.Value)

	metric.Hash = CalcHash(metric.String(), "bar")
	assert.Equal(t, "e1106794864f57461d121d23149ed87659726933abaa3f8019d82f2fa022052c", *metric.Hash)
	assert.Nil(t, metric.Delta)
}

func TestLoadJsonMetric(t *testing.T) {
	errorMsg := func(err error) string {
		if err == nil {
			return ""
		}
		return err.Error()
	}
	goodDelta := int64(123)
	badDelta := int64(0)
	goodValue := 123.456
	badValue := float64(0)
	hash := "any_hash"
	tests := []struct {
		name   string
		body   []byte
		want   JSONMetric
		errMsg string
	}{
		{
			name:   "success",
			body:   []byte(`{"ID":"PollCount","type":"counter","Delta":123,"Value":123.456}`),
			want:   JSONMetric{"PollCount", "counter", &goodDelta, &goodValue, nil},
			errMsg: "",
		},
		{
			name:   "bad JSON input",
			body:   []byte(`{"ID":"PollCount","type":"counter","Delta":123,"Value":123.456`),
			want:   JSONMetric{},
			errMsg: "unexpected end of JSON input",
		},
		{
			name:   "bad JSON input 2",
			body:   []byte(`"ID":"PollCount","type":"counter","Delta":123,"Value":123.456}`),
			want:   JSONMetric{},
			errMsg: "invalid character ':' after top-level value",
		},
		{
			name:   "bad JSON input 3",
			body:   []byte(`false`),
			want:   JSONMetric{},
			errMsg: "json: cannot unmarshal bool into Go value of type utils.JSONMetric",
		},
		{
			name:   "bad int64",
			body:   []byte(`{"ID":"PollCount","type":"counter","Delta":"123","Value":123.456}`),
			want:   JSONMetric{"PollCount", "counter", &badDelta, &goodValue, nil},
			errMsg: "json: cannot unmarshal string into Go struct field JSONMetric.delta of type int64",
		},
		{
			name:   "bad float64",
			body:   []byte(`{"ID":"PollCount","type":"counter","Delta":123,"Value":"123.456"}`),
			want:   JSONMetric{"PollCount", "counter", &goodDelta, &badValue, nil},
			errMsg: "json: cannot unmarshal string into Go struct field JSONMetric.value of type float64",
		},
		{
			name:   "bad int64 and float64 and hash",
			body:   []byte(`{"ID":"PollCount","type":"counter","Delta":"123","Value":"123.456", "hash": "any_hash"}`),
			want:   JSONMetric{"PollCount", "counter", &badDelta, &badValue, &hash},
			errMsg: "json: cannot unmarshal string into Go struct field JSONMetric.delta of type int64",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := LoadJSONMetric(tt.body)
			assert.Equal(t, tt.errMsg, errorMsg(err))
			assert.Equalf(t, tt.want, got, "LoadJSONMetric(%v)", tt.body)
		})
	}
}

func TestJsonMetric_IsValidType(t *testing.T) {
	tests := []struct {
		name  string
		mType string
		want  bool
	}{
		{name: "counter type", mType: "counter", want: true},
		{name: "gauge type", mType: "gauge", want: true},
		{name: "wrong type", mType: "wrong", want: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := JSONMetric{ID: "name", MType: tt.mType}
			assert.Equal(t, tt.want, m.IsValidType())
		})
	}
}

func TestJsonMetric_IsValidValue(t *testing.T) {
	goodDelta := int64(123)
	goodValue := 123.456

	tests := []struct {
		name  string
		mType string
		delta *int64
		value *float64
		want  bool
	}{
		{name: "counter valid value", mType: "counter", delta: &goodDelta, value: nil, want: true},
		{name: "counter bad value", mType: "counter", delta: nil, value: &goodValue, want: false},
		{name: "gauge valid value", mType: "gauge", delta: nil, value: &goodValue, want: true},
		{name: "gauge bad value", mType: "gauge", delta: &goodDelta, value: nil, want: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := JSONMetric{ID: "name", MType: tt.mType, Delta: tt.delta, Value: tt.value}
			assert.Equal(t, tt.want, m.IsValidValue())
		})
	}
}
