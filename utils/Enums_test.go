package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Values_EmptyEnum_ReturnsNotFoundError(t *testing.T) {
	var testEnum Enum
	testIdx := 0
	val, err := testEnum.Value(testIdx)

	assert.Nil(t, val)
	assert.NotNil(t, err)

}
