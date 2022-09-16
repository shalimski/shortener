package randomstring_test

import (
	"testing"

	"github.com/shalimski/shortener/pkg/randomstring"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	first := randomstring.New(7)
	assert.NotEmpty(t, first)

	second := randomstring.New(7)
	assert.NotEmpty(t, second)

	assert.Equal(t, len(first), len(second))

	third := randomstring.New(5)
	assert.NotEmpty(t, third)

	assert.Greater(t, len(second), len(third))

	assert.NotEqual(t, first, second)
	assert.NotEqual(t, first, third)
}

