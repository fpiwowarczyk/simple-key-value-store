package store

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEmptyGet(t *testing.T) {
	val, err := Get("no existing key")

	assert.Equal(t, "", val, "Value should be empty")
	assert.Error(t, err, "Error should be returned")
	assert.Equal(t, ErrorNoSuchKey, err, "Should return error no such key")
}

func TestPut(t *testing.T) {
	Put("key1", "value1")

	val, err := Get("key1")

	assert.Equal(t, "value1", val, "Value should be returned")
	assert.NoError(t, err, "Error should be nil")

	Delete("key1")
}

func TestDelete(t *testing.T) {
	Put("key1", "value1")

	val, _ := Get("key1")
	assert.NotEmpty(t, val, "Some value should be returned")

	err := Delete("key1")
	assert.NoError(t, err, "Delete should not throw error ")
	val, err = Get("key1")
	assert.Equal(t, "", val, "Value should be empty")
	assert.Error(t, err, "Error should be returned")
	assert.Equal(t, ErrorNoSuchKey, err, "Should return error no such key")

}
