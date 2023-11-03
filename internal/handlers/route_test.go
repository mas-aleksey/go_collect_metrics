package handlers

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/tiraill/go_collect_metrics/internal/utils"
)

func TestGetRouter(t *testing.T) {
	router := GetRouter(nil, utils.ServerConfig{})
	assert.NotNil(t, router)
}
