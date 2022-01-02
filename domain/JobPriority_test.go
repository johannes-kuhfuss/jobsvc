package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_IsValidPriority_InvalidPriority_Returns_False(t *testing.T) {
	valid := IsValidPriority("bogus")

	assert.NotNil(t, valid)
	assert.EqualValues(t, false, valid)
}

func Test_IsValidPriority_ValidPriority_Returns_True(t *testing.T) {
	valid := IsValidPriority("high")

	assert.NotNil(t, valid)
	assert.EqualValues(t, true, valid)
}
