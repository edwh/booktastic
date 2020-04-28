package main

import (
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"log"
	"testing"
)

func TestIdentifyBooks(t *testing.T) {
	log.SetOutput(ioutil.Discard)
	log.SetFlags(0)

	lines, fragments := GetLinesAndFragments(SAMPLE)
	spines := ExtractSpines(lines, fragments)

	IdentifyBooks(spines, fragments)
	assert.NotNil(t, spines[0].Author)
}
