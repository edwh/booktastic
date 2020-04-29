package main

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"log"
	"path/filepath"
	"strings"
	"testing"
)

var tests = []string{
	"vertical_easy",
	"liz1",
	"liz2",
	"liz3",
	"liz4",
	"liz5",
	"liz7",
	"liz8",
	"liz9",
	"liz10",
	"liz11",
	"liz13",
	"liz14",
	"liz15",
	"liz16",
	"liz17",
	"liz18",
	"liz19",
	"liz20",
	"liz21",
	"liz22",
	"liz23",
	"liz24",
	"liz25",
	"liz26",
	"liz27",
	"liz28",
	"liz29",
	"liz30",
	"liz31",
	"liz33",
	"liz34",
	"liz35",
	"liz36",
	"liz37",
	"liz38",
	"liz39",
	"liz40",
	"liz41",
	"liz43",
	"liz44",
	"liz45",
	"liz46",
	"liz47",
	"liz48",
	"ruth1",
	"ruth2",
	"ruth3",
	"jo1",
	"carol1",
	"carol2",
	"kathryn1",
	"phil1",
	"doug1",
	"doug2",
	"doug3",
	"adam1",
	"adam2",
	"andy1",
	"emma1",
	"suzanne1",
	"suzanne2",
	"suzanne3",
	"suzanne4",
	"suzanne5",
	"tom1",
	"wanda1",
	"caroline1",
	"bryson3",
	"bryson2",
	"bryson",
	"chris1",
	"chris2",
	"crime1",
	"crime2",
	"crime3",
	"basic_horizontal",
	"basic_vertical",
	"gardening",
	"horizontal_overlap",
	"horizontal_overlap2",
}

func runTest(t *testing.T, tests []string) {
	// Run our tests.
	failed := false

	for _, fn := range tests {
		t.Run(fn, func(t *testing.T) {
			ifn := "testdata" + string(filepath.Separator) + fn + ".json"
			data, _ := ioutil.ReadFile(ifn)

			lines, fragments := GetLinesAndFragments(string(data))
			spines, fragments := ExtractSpines(lines, fragments)
			spines, fragments = IdentifyBooks(spines, fragments)

			sugar.Debugf("Spines after test %+v", spines)

			ofn := "testdata" + string(filepath.Separator) + fn + "_books.json"
			sugar.Debugf("Output file %s", ofn)
			odata, _ := ioutil.ReadFile(ofn)

			if len(odata) > 0 {
				sugar.Debugf("Output data %s", odata)

				missed := []Spine{}

				ospines := []Spine{}
				json.Unmarshal([]byte(odata), &ospines)

				olduns := 0
				newuns := 0

				for _, spine := range spines {
					if len(spine.Author) > 0 {
						newuns++
						missing := true
						for _, ospine := range ospines {
							if strings.Compare(strings.ToLower(spine.Author), strings.ToLower(ospine.Author)) == 0 &&
								strings.Compare(strings.ToLower(spine.Title), strings.ToLower(ospine.Title)) == 0 {
								missing = false
							}
						}

						if missing {
							t.Errorf("NOW FOUND: %s - %s\n", spine.Author, spine.Title)
							failed = true
						}
					}
				}

				for ospineindex, ospine := range ospines {
					if len(ospine.Author) > 0 {
						olduns++
						missing := true
						for spineindex, spine := range spines {
							if strings.Compare(strings.ToLower(spine.Author), strings.ToLower(ospine.Author)) == 0 &&
								strings.Compare(strings.ToLower(spine.Title), strings.ToLower(ospine.Title)) == 0 {
								sugar.Debugf("MATCHED: %s - %s at %d vs %d", spine.Author, spine.Title, spineindex, ospineindex)
								missing = false
							}
						}

						if missing {
							missed = append(missed, ospine)
							failed = true
						}
					}
				}

				for _, miss := range missed {
					t.Errorf("MISSING: %s - %s\n", miss.Author, miss.Title)
				}

				if newuns > olduns {
					sugar.Debugf("Better %d vs %d", newuns, olduns)
				} else if newuns < olduns {
					sugar.Debugf("Worse %d vs %d", newuns, olduns)
				} else {
					sugar.Debugf("Same %d", olduns)
				}

			} else {
				sugar.Debugf("No output yet")

				encoded, _ := json.MarshalIndent(spines, "", " ")
				sugar.Debugf(string(encoded))
			}
		})
	}

	assert.False(t, failed)
}

func TestIdentifyBooks(t *testing.T) {
	log.SetOutput(ioutil.Discard)
	log.SetFlags(0)

	runTest(t, tests)
}

func TestEasy(t *testing.T) {
	log.SetOutput(ioutil.Discard)
	log.SetFlags(0)

	runTest(t, []string{
		"vertical_easy",
	})
}
