// Distributed short URL generator,
// each node has an ints interval for generating links in order.
// Intervals of different nodes do not overlap, a coordinator is responsible for this
package generator

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewUrlGenerator(t *testing.T) {
	// start position
	gen, err := NewURLGenerator(&MockCounter{})

	assert.NoError(t, err)

	ctx := context.Background()
	first, err := gen.Next(ctx)
	assert.NoError(t, err)
	assert.Equal(t, first, "b")

	second, err := gen.Next(ctx)
	assert.NoError(t, err)
	assert.Equal(t, second, "c")

	// fixed position
	gen, err = NewURLGenerator(&MockCounter{current: 42})

	assert.NoError(t, err)
	first, err = gen.Next(ctx)
	assert.NoError(t, err)
	assert.Equal(t, first, "rML7")

	second, err = gen.Next(ctx)
	assert.NoError(t, err)
	assert.Equal(t, second, "rML8")
}

type MockCounter struct {
	current int
}

func (m *MockCounter) NextCounter(ctx context.Context) (int, error) {
	m.current = m.current + 1
	return m.current, nil
}

func TestEncode(t *testing.T) {
	tests := []struct {
		name string
		num  int
		want string
	}{
		{"zero", 0, ""},
		{"one", 1, "b"},
		{"onemillion", 1_000_000, "emjc"},
		{"negative", -4, ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Encode(tt.num); got != tt.want {
				t.Errorf("Encode() = %v, want %v", got, tt.want)
			}
		})
	}
}
