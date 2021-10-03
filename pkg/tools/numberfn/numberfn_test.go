package numberfn_test

import (
	"github.com/bernardosecades/feeder/pkg/tools/numberfn"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestIsAllNumbers(t *testing.T) {
	assert.True(t, numberfn.IsUnsignedInteger("123"))
	assert.True(t, numberfn.IsUnsignedInteger("000123"))

	assert.False(t, numberfn.IsUnsignedInteger("00A123"))
	assert.False(t, numberfn.IsUnsignedInteger("1.5"))
	assert.False(t, numberfn.IsUnsignedInteger("-123"))
}
