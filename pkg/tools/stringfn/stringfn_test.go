package stringfn_test

import (
	"github.com/bernardosecades/feeder/pkg/tools/stringfn"

	"github.com/stretchr/testify/assert"

	"testing"
)

func TestIsAllLetters(t *testing.T) {
	assert.True(t, stringfn.IsOnlyLetters("ABCD"))
	assert.True(t, stringfn.IsOnlyLetters("abc"))

	assert.False(t, stringfn.IsOnlyLetters("ABC1"))
	assert.False(t, stringfn.IsOnlyLetters("ABC-"))
	assert.False(t, stringfn.IsOnlyLetters("1987"))
	assert.False(t, stringfn.IsOnlyLetters("abc@"))
}
