package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCleanOCR(t *testing.T) {
	assert.Equal(t, "Hi", CleanOCR("Hi ISBN"))
	assert.Equal(t, "Hi", CleanOCR(" # 0123.45 Hi\" -hi"))
}
