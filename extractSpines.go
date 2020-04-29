package main

import (
	"encoding/json"
	"log"
	"math"
	"regexp"
	"strings"
)

const PRUNE_SMALL_TEXT = 4

// Google OCR returns an array of these.
type OCRFragment struct {
	Locale       string
	Description  string
	BoundingPoly BoundingPoly
	SpineIndex   int
	Used         bool
}

type BoundingPoly struct {
	Vertices []Vertices
}

type Vertices struct {
	X int
	Y int
}

type Spine struct {
	Spine  string // The current working text
	Author string // Identified author
	Title  string // Identified subject
}

func GetLinesAndFragments(str string) ([]string, []OCRFragment) {
	// TODO Need default values of X and Y as Google omits if 0.
	var m []OCRFragment
	json.Unmarshal([]byte(str), &m)

	// First entry is a summary, with newline separators for related text.
	summary := m[0].Description
	lines := strings.Split(summary, "\n")
	sugar.Debugf("Description %s", lines[0])

	// Remaining entries are the fragments.
	fragments := m[1:]
	for _, f := range fragments {
		f.Used = false
	}

	return lines, fragments
}

func MaxDimension(poly BoundingPoly) int {
	vertices := poly.Vertices

	x := math.Abs(float64(vertices[0].X) - float64(vertices[3].X))
	y := math.Abs(float64(vertices[0].Y) - float64(vertices[3].Y))

	return int(math.Max(x, y))
}

func PruneSmallText(lines []string, fragments []OCRFragment, ratio int) ([]string, []OCRFragment, int) {
	// Very small text on spines is likely to be publishers, ISBN numbers, stuff we've read from the front at an angle,
	// or otherwise junk.  So let's identify the typical letter height, and prune out stuff that's much smaller.
	total := 0
	pruned := 0

	for _, fragment := range fragments {
		total += MaxDimension(fragment.BoundingPoly)
	}

	mean := total / len(fragments)

	sugar.Debugf("Mean max %d", mean)

	newlines := []string{}
	newfragments := []OCRFragment{}
	fragindex := 0

	for lineindex, line := range lines {
		linewords := strings.Split(line, " ")
		newlinewords := []string{}

		for _, word := range linewords {
			if len(word) > 0 {
				sugar.Debugf("Consider %s line %d, fragment %d vs %d", word, lineindex, fragindex, len(fragments))

				fragment := fragments[fragindex]
				if word != fragment.Description {
					sugar.Debugf("ERROR: Mismatch spine/fragment %s vs %s", word, fragment.Description)
					panic("Mismatch spine/fragment")
				} else {
					thismax := MaxDimension(fragment.BoundingPoly)
					sugar.Debugf("Max %d vs %d", thismax, mean)

					if thismax*ratio < mean {
						sugar.Debugf("Prune small text %s size %d vs %d at %d", fragment.Description, thismax, mean, fragindex)
						pruned++
					} else {
						newlinewords = append(newlinewords, fragment.Description)
						fragment.SpineIndex = len(newlines)
						newfragments = append(newfragments, fragment)
					}

					fragindex++
				}
			}
		}

		newlines = append(newlines, strings.Join(newlinewords, " "))
	}

	return newlines, newfragments, pruned
}

func CleanOCR(str string) string {
	// TODO Wasteful to compile the regexp each time.
	newstr := str

	// ISBNs often appear on spines.
	newstr = regexp.MustCompile(`(?i)ISBN`).ReplaceAllString(newstr, "")

	// Anything with digits separated by dots can't be a real word.
	newstr = regexp.MustCompile(`\d+\.\d+`).ReplaceAllString(newstr, "")

	// Anything with leading zeros can't either.
	newstr = regexp.MustCompile(`0\d+`).ReplaceAllString(newstr, "")

	// Remove all words that are 1, 2 or 3 digits.  These could legitimately be in some titles but
	// much more often they are ISBN junk.
	newstr = regexp.MustCompile(`\b\d{1,3}\b`).ReplaceAllString(newstr, "")

	// Nothing good starts with a dash.
	newstr = regexp.MustCompile(`\s-\w+(\b|$)`).ReplaceAllString(newstr, "")

	// # is not a word
	newstr = regexp.MustCompile(`\s#\s`).ReplaceAllString(newstr, "")

	// Quotes confuse matters.
	newstr = regexp.MustCompile(`["']`).ReplaceAllString(newstr, "")

	// Collapse multiple spaces.
	newstr = regexp.MustCompile(`\s+`).ReplaceAllString(newstr, " ")

	newstr = strings.TrimSpace(newstr)

	if str != newstr {
		sugar.Debugf("Cleaned %s => %s", str, newstr)
	}

	return newstr
}

func AddSpineIndex(lines []string, fragments []OCRFragment) []OCRFragment {
	fragindex := 0

	for spineindex, line := range lines {
		words := strings.Split(strings.TrimSpace(line), " ")

		for _, word := range words {
			if len(word) > 0 {
				if word == fragments[fragindex].Description {
					fragments[fragindex].SpineIndex = spineindex
					sugar.Debugf("Frag %d index %d contents %s", fragindex, spineindex, fragments[fragindex].Description)
					fragindex++
				} else {
					log.Fatalf("Mismatch adding spine index %s vs %s", word, fragments[fragindex].Description)
				}
			}
		}
	}

	return fragments
}

func ExtractSpines(lines []string, fragments []OCRFragment) ([]Spine, []OCRFragment) {
	spines := []Spine{}

	fragments = AddSpineIndex(lines, fragments)
	lines, fragments, _ = PruneSmallText(lines, fragments, PRUNE_SMALL_TEXT)

	for _, line := range lines {
		cleaned := CleanOCR(line)

		if len(cleaned) > 0 {
			spines = append(spines, Spine{
				Spine:  cleaned,
				Author: "",
				Title:  "",
			})
		} else {
			// We're removing this spine.  Remove any fragments with this spine index.
			fragments = removeFragmentsForSpine(len(spines), fragments)
		}
	}

	return spines, fragments
}

func removeFragmentsForSpine(spineindex int, fragments []OCRFragment) []OCRFragment {
	newfrags := []OCRFragment{}

	for _, frag := range fragments {
		if frag.SpineIndex != spineindex {
			newfrag := frag

			if newfrag.SpineIndex > spineindex {
				newfrag.SpineIndex--
			}

			newfrags = append(newfrags, newfrag)
		}
	}

	return newfrags
}
