package utils

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewMetric(t *testing.T) {
	expected := Metric{
		Type:  MetricType("counter"),
		Name:  "RandomValue",
		Value: "1",
	}
	metric := NewMetric("counter", "RandomValue", "1")
	assert.Equal(t, expected, metric)
}

func TestMetric_IsValidType(t *testing.T) {
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
			m := NewMetric(tt.mType, "name", "value")
			assert.Equal(t, tt.want, m.IsValidType())
		})
	}
}

func TestMetric_IsValidValue(t *testing.T) {
	tests := []struct {
		name   string
		mType  string
		mName  string
		mValue string
		want   bool
	}{
		{name: "gauge valid value 1", mType: "gauge", mName: "RandomValue", mValue: "0", want: true},
		{name: "gauge valid value 2", mType: "gauge", mName: "RandomValue", mValue: "123.45", want: true},
		{name: "gauge invalid value", mType: "gauge", mName: "Alloc", mValue: "foo", want: false},
		{name: "counter valid value", mType: "counter", mName: "PollCount", mValue: "5", want: true},
		{name: "counter invalid value", mType: "counter", mName: "Wrong", mValue: "12s", want: false},
		{name: "wrong type", mType: "wrong", mName: "name", mValue: "0", want: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := NewMetric(tt.mType, tt.mName, tt.mValue)
			assert.Equal(t, tt.want, m.IsValidValue())
		})
	}
}
