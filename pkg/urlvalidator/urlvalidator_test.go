package urlvalidator_test

import (
	"fmt"
	"testing"

	"github.com/shalimski/shortener/pkg/urlvalidator"
	"github.com/stretchr/testify/assert"
)

func TestIsURL(t *testing.T) {
	t.Parallel()

	tests := []struct {
		param    string
		expected bool
	}{
		{"", false},
		{"http://foo.bar#com", true},
		{"http://foobar.com", true},
		{"https://foobar.com", true},
		{"foobar.com", true},
		{"http://foobar.coffee/", true},
		{"http://foobar.中文网/", true},
		{"http:www.example.com/main.html", true},
	}
	for _, test := range tests {
		actual := urlvalidator.IsURL(test.param)
		assert.Equal(t, test.expected, actual, fmt.Sprintf("IsURL(%q)", test.param))
	}
}

func TestIsShortURLSuffix(t *testing.T) {
	t.Parallel()

	tests := []struct {
		param    string
		expected bool
	}{
		{"", false},
		{"abc", true},
		{"azip50fke", true},
		{"-----", false},
		{"foobar.com", false},
		{"111111111", true},
		{"123456789011", false},
	}
	for _, test := range tests {
		actual := urlvalidator.IsShortURLSuffix(test.param)
		assert.Equal(t, test.expected, actual, fmt.Sprintf("IsShortURLSuffix(%q)", test.param))
	}
}
