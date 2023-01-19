package handlers

import (
	"github.com/stretchr/testify/assert"
	"github.com/tiraill/go_collect_metrics/internal/utils"
	"testing"
)

func TestGetRouter(t *testing.T) {
	router := GetRouter(nil, utils.ServerConfig{})
	assert.NotNil(t, router)
}
