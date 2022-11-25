package handlers

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetRouter(t *testing.T) {
	router := GetRouter(nil)
	assert.NotNil(t, router)
}
