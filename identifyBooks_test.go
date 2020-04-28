package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestIdentifyBooks(t *testing.T) {
	//log.SetOutput(ioutil.Discard)
	//log.SetFlags(0)

	lines, fragments := GetLinesAndFragments(SAMPLE)
	spines, fragments := ExtractSpines(lines, fragments)

	IdentifyBooks(spines, fragments)
	assert.NotNil(t, spines[0].Author)
}
