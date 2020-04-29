package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

const SAMPLE = ``

func TestCleanOCR(t *testing.T) {
	assert.Equal(t, "Hi", CleanOCR("Hi ISBN"))
	assert.Equal(t, "Hi", CleanOCR(" # 0123.45 Hi\" -hi"))
}

func TestDimension(t *testing.T) {
	lines, fragments := GetLinesAndFragments(SAMPLE)
	assert.Equal(t, "PMC", lines[0])
	assert.Equal(t, "PMC", fragments[0].Description)
	assert.Equal(t, 333, fragments[0].BoundingPoly.Vertices[0].X)
	assert.Equal(t, 36, MaxDimension(fragments[0].BoundingPoly))
}

func TestPruneSmallText(t *testing.T) {
	lines, fragments := GetLinesAndFragments(SAMPLE)

	// Force some pruning.
	_, _, pruned := PruneSmallText(lines, fragments, 1)
	assert.Equal(t, 83, pruned)
	_, _, pruned = PruneSmallText(lines, fragments, PRUNE_SMALL_TEXT)
	assert.Equal(t, 1, pruned)
}

func TestIdentifySpines(t *testing.T) {
	lines, fragments := GetLinesAndFragments(SAMPLE)
	spines, fragments := ExtractSpines(lines, fragments)
	sugar.Debugf("Spine %+v", spines[0])
	assert.Equal(t, "PMC", spines[0].Spine)
}
