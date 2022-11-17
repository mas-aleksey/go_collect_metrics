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

func TestIsValidCounterMetricName(t *testing.T) {
	tests := []struct {
		name  string
		mName string
		want  bool
	}{
		{name: "counter PollCount name", mName: "PollCount", want: true},
		{name: "counter Wrong name", mName: "Wrong", want: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := NewMetric("counter", tt.mName, "value")
			assert.Equal(t, tt.want, m.IsValidName())
		})
	}
}

func TestIsValidGaugeMetricName(t *testing.T) {
	tests := []struct {
		name  string
		mName string
		want  bool
	}{
		{name: "gauge RandomValue name", mName: "RandomValue", want: true},
		{name: "gauge Alloc name", mName: "Alloc", want: true},
		{name: "gauge HeapObjects name", mName: "HeapObjects", want: true},
		{name: "gauge MSpanSys name", mName: "MSpanSys", want: true},
		{name: "gauge TotalAlloc name", mName: "TotalAlloc", want: true},
		{name: "gauge Wrong name", mName: "Wrong", want: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := NewMetric("gauge", tt.mName, "value")
			assert.Equal(t, tt.want, m.IsValidName())
		})
	}
}

func TestMetric_IsValidName(t *testing.T) {
	tests := []struct {
		name  string
		mType string
		mName string
		want  bool
	}{
		{name: "gauge RandomValue name", mType: "gauge", mName: "RandomValue", want: true},
		{name: "gauge Alloc name", mType: "gauge", mName: "Alloc", want: true},
		{name: "gauge HeapObjects name", mType: "gauge", mName: "HeapObjects", want: true},
		{name: "gauge MSpanSys name", mType: "gauge", mName: "MSpanSys", want: true},
		{name: "gauge TotalAlloc name", mType: "gauge", mName: "TotalAlloc", want: true},
		{name: "gauge Wrong name", mType: "gauge", mName: "Wrong", want: false},
		{name: "counter PollCount name", mType: "counter", mName: "PollCount", want: true},
		{name: "counter Wrong name", mType: "counter", mName: "Wrong", want: false},
		{name: "wrong type", mType: "wrong", mName: "name", want: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := NewMetric(tt.mType, tt.mName, "value")
			assert.Equal(t, tt.want, m.IsValidName())
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

func Test_contains(t *testing.T) {
	type args struct {
		value string
		array []string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{name: "contains in arr", args: args{"Alloc", RuntimeMetricNames}, want: true},
		{name: "not contains in arr", args: args{"Wrong", RuntimeMetricNames}, want: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, contains(tt.args.value, tt.args.array))
		})
	}
}
