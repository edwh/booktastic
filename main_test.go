package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestHelloWorld(t *testing.T) {
	want := "Hi"
	got, err := HelloWorld("Hi");
	assert.Nil(t, err)
	assert.Equal(t, want, got)
}
