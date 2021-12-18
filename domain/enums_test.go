package domain

import (
	"fmt"
	"net/http"
	"testing"
)

func Test_Values_EmptyEnum_ReturnsNotFoundError(t *testing.T) {
	var testEnum EnumList
	testIdx := 0
	val, err := testEnum.Value(testIdx)
	if val != "" {
		t.Errorf("Expected empty string, but got %v", val)
	}
	if err == nil {
		t.Error("Expected error, but got nil")
	}
	if err.StatusCode() != http.StatusNotFound {
		t.Errorf("Expected status not found, but got %v", err.StatusCode())
	}
	if err.Message() != fmt.Sprintf("No item with index %v found", testIdx) {
		t.Errorf("Expected No item with index %v found but got %v", testIdx, err.Message())
	}
}

/*
func Test_A_Values_EmptyEnum_ReturnsNotFoundError(t *testing.T) {
	var testEnum EnumList
	testIdx := 0
	_, err := testEnum.Value(testIdx)

	assert.NotNil(t, err)
}
*/
