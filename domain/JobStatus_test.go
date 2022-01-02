package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_IsValidJobStatus_InvalidStatus_Returns_False(t *testing.T) {
	valid := IsValidJobStatus("bogus")

	assert.NotNil(t, valid)
	assert.EqualValues(t, false, valid)
}

func Test_IsValidJobStatus_ValidStatus_Returns_True(t *testing.T) {
	valid := IsValidJobStatus("running")

	assert.NotNil(t, valid)
	assert.EqualValues(t, true, valid)
}
