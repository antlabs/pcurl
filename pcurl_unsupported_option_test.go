package pcurl

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_UnsupportedOption(t *testing.T) {
	_, err := ParseAndRequest(`curl --hahaha`)
	assert.Error(t, err)
}
