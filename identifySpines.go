package main

import (
	"encoding/json"
	"log"
	"regexp"
	"strings"
)

// Google OCR returns an array of these.
type OCRFragment struct {
	Locale       string
	Description  string
	BoundingPoly interface{}
}

type Spine struct {
	original string // The original text returned from the OCR.
	spine    string // The current working text
	author   string // Identified author
	title    string // Identified subject
}

func GetLinesAndFragments(str string) ([]string, []OCRFragment) {
	var m []OCRFragment
	json.Unmarshal([]byte(str), &m)

	// First entry is a summary, with newline separators for related text.
	summary := m[0].Description
	lines := strings.Split(summary, "\n")
	log.Printf("Description %s", lines[0])

	// Remaining entries are the fragments.
	fragments := m[1:]

	return lines, fragments
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
	newstr = regexp.MustCompile(`"`).ReplaceAllString(newstr, "")

	// Collapse multiple spaces.
	newstr = regexp.MustCompile(`\s+`).ReplaceAllString(newstr, " ")

	newstr = strings.TrimSpace(newstr)

	if str != newstr {
		log.Printf("Cleaned %s => %s", str, newstr)
	}

	return newstr
}
