package main

import (
	"log"
	"sync"
	"time"
)

type phase struct {
	id              int
	fuzzy           bool
	authorstart     bool
	authorplustitle bool
	permuted        bool
	mangled         bool
}

func IdentifyBooks(spines []Spine, fragments []OCRFragment) {
	// We scan the text to identify spines.  We have various techniques for this:
	//
	// - author at start/end of spine
	// - match on author + title or match on author and scan titles
	// - change the order of the spines (permuted spines)
	// - change the order of the words (mangled spines)
	//
	// Some of these are more expensive than others, especially the ordering ones, and work better if we've already
	// identified and flagged as much as possible.  So we run through several phases doing these tests in different
	// orders.
	//
	// The order has been chosen by empirical testing as a combination of success and time - generally
	// the earlier ones are quicker.  If a combination doesn't appear then it has not been effective.
	//
	// Empirically, if we don't find anything in an earlier phase, it isn't worth going on to a later phase.
	phases := setUpPhases()

	// We use a wait group to execute all the phases in parallel.
	var wg sync.WaitGroup

	for _, p := range phases {
		wg.Add(1)
		go func(phase phase) {
			defer wg.Done()
			log.Printf("Execute phase %+v", phase)
			start := time.Now()
			found := 0
			duration := time.Since(start)
			log.Printf("Phase %d %+v found %d in %v", phase.id, p, found, duration)
		}(p)
	}

	wg.Wait()
	log.Printf("All phases complete")
}

func setUpPhases() []phase {
	phases := []phase{}
	bools := [2]bool{true, false}
	id := 0

	for _, fuzzy := range bools {
		for _, mangled := range bools {
			for _, permuted := range bools {
				if !mangled || !permuted {
					for _, authorplustitle := range bools {
						for _, authorstart := range bools {
							phases = append(phases, phase{
								id,
								fuzzy,
								authorstart,
								authorplustitle,
								permuted,
								mangled,
							})

							id++
						}
					}
				}
			}
		}
	}

	return phases
}

func countSuccess(spines []Spine) int {
	count := 0

	for _, spine := range spines {
		if spine.Author != nil {
			count++
		}
	}

	return count
}

func searchForSpines(spines []Spine, fragments []OCRFragment) {
	// We want to search for the spines in ElasticSearch, where we have a list of authors and books.
	//
	// The spine will normally in in the format "Author Title" or "Title Author".  So we can work our
	// way along the words in the spine searching for matches on this.
	//
	// Order our search by longest spine first.  This is because the longer the spine is, the more likely
	// it is to have both and author and a subject, and therefore match.  Matching gets it out of the way
	// but also gives us a known author, which can be used to good effect to improve matching on other
	// spines.
}
