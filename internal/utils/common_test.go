package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestToStr(t *testing.T) {
	t.Run("uint64 to str", func(t *testing.T) {
		var input uint64 = 111
		out := ToStr(input)
		assert.Equal(t, "111", out)
	})
	t.Run("int64 to str", func(t *testing.T) {
		var input int64 = 222
		out := ToStr(input)
		assert.Equal(t, "222", out)
	})
	t.Run("float64 to str", func(t *testing.T) {
		var input = 333.333
		out := ToStr(input)
		assert.Equal(t, "333.333", out)
	})
	t.Run("uint32 to str", func(t *testing.T) {
		var input uint32 = 444
		out := ToStr(input)
		assert.Equal(t, "444", out)
	})
	t.Run("something else to str", func(t *testing.T) {
		input := "foo555"
		out := ToStr(input)
		assert.Equal(t, "0", out)
	})
}

func TestToFloat64(t *testing.T) {
	t.Run("uint64 to float64", func(t *testing.T) {
		var input uint64 = 111
		out := ToFloat64(input)
		assert.Equal(t, float64(111), out)
	})
	t.Run("int64 to float64", func(t *testing.T) {
		var input int64 = 222
		out := ToFloat64(input)
		assert.Equal(t, float64(222), out)
	})
	t.Run("float64 to float64", func(t *testing.T) {
		var input = 333.333
		out := ToFloat64(input)
		assert.Equal(t, 333.333, out)
	})
	t.Run("uint32 to float64", func(t *testing.T) {
		var input uint32 = 444
		out := ToFloat64(input)
		assert.Equal(t, float64(444), out)
	})
	t.Run("something else to float64", func(t *testing.T) {
		input := "foo555"
		out := ToFloat64(input)
		assert.Equal(t, float64(0), out)
	})
}

func TestCalcHash(t *testing.T) {
	t.Run("equal hash", func(t *testing.T) {
		assert.Equal(t, CalcHash("some_data", "123"), CalcHash("some_data", "123"))
	})
	t.Run("not equal keys", func(t *testing.T) {
		assert.NotEqual(t, CalcHash("some_data", "123"), CalcHash("some_data", "321"))
	})
	t.Run("not equal data", func(t *testing.T) {
		assert.NotEqual(t, CalcHash("some_data1", "123"), CalcHash("some_data2", "123"))
	})
	t.Run("not equal data and keys", func(t *testing.T) {
		assert.NotEqual(t, CalcHash("some_data1", "123"), CalcHash("some_data2", "321"))
	})
}
