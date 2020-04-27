package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestIdentifyBooks(t *testing.T) {
	lines, fragments := GetLinesAndFragments(SAMPLE)
	spines := ExtractSpines(lines, fragments)

	IdentifyBooks(spines, fragments)
	assert.NotNil(t, spines[0].Author)
}
