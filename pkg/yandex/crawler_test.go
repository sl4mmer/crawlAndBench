package yandex

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSearch(t *testing.T) {
	queries := []string{"тест"}
	for _, q := range queries {
		res, err := Search(q)
		assert.NoError(t, err)
		assert.True(t, len(res) > 0)
	}
}
