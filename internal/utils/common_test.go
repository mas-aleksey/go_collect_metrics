package utils

import (
	"github.com/stretchr/testify/assert"
	"testing"
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
