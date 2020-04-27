package main

import (
	"regexp"
	"strings"
)

func CleanOCR(str string) string {
	// ISBNs often appear on spines.
	str = regexp.MustCompile(`(?i)ISBN`).ReplaceAllString(str, "")

	// Anything with digits separated by dots can't be a real word.
	str = regexp.MustCompile(`\d+\.\d+`).ReplaceAllString(str, "")

	// Anything with leading zeros can't either.
	str = regexp.MustCompile(`0\d+`).ReplaceAllString(str, "")

	// Remove all words that are 1, 2 or 3 digits.  These could legitimately be in some titles but
	// much more often they are ISBN junk.
	str = regexp.MustCompile(`\b\d{1,3}\b`).ReplaceAllString(str, "")

	// Nothing good starts with a dash.
	str = regexp.MustCompile(`\s-\w+(\b|$)`).ReplaceAllString(str, "")

	// # is not a word
	str = regexp.MustCompile(`\s#\s`).ReplaceAllString(str, "")

	// Quotes confuse matters.
	str = regexp.MustCompile(`"`).ReplaceAllString(str, "")

	// Collapse multiple spaces.
	str = regexp.MustCompile(`\s+`).ReplaceAllString(str, " ")

	return strings.TrimSpace(str)
}
