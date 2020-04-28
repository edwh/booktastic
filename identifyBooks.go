package main

import (
	"log"
	"regexp"
	"sort"
	"strings"
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

const MAXRESULTS = 1000

type searchResult struct {
	spineindex   int
	searchAuthor string
	searchTitle  string
	foundAuthor  string
	foundTitle   string
}

var searchResults map[searchResult]searchResult
var resultsMux sync.Mutex

func addResult(result searchResult) {
	log.Printf("Add result %+v", result)
	// We want the results to be unique by found values in a single spine.  We might have searched on different
	// variants but found the same thing.
	key := result
	key.searchTitle = ""
	key.searchAuthor = ""
	resultsMux.Lock()

	searchResults[key] = result
	resultsMux.Unlock()
}

func IdentifyBooks(spines []Spine, fragments []OCRFragment) {
	phases := setUpPhases()

	// We need to execute the phases serially as the results of one phase make it more likely that we can find things
	// in later phases.
	for _, p := range phases {
		log.Printf("Execute phase %+v", p)
		start := time.Now()
		searchResults = map[searchResult]searchResult{}
		searchSpines(spines, fragments, p)
		spines, fragments = processSearchResults(spines, fragments)
		duration := time.Since(start)
		log.Printf("Phase %d %+v found %d in %v", p.id, p, len(searchResults), duration)
	}

	log.Printf("All phases complete")

	for _, frag := range fragments {
		if !frag.Used {
			log.Printf("LEFTOVER: spine %d %s", frag.SpineIndex, frag.Description)
		}
	}

	for _, spine := range spines {
		if len(spine.Author) > 0 {
			log.Printf("RESULT: %s - %s", spine.Author, spine.Title)
		}
	}
}

func processSearchResults(spines []Spine, fragments []OCRFragment) ([]Spine, []OCRFragment) {
	// Get the results as an array.
	results := make([]searchResult, 0, len(searchResults))
	for _, v := range searchResults {
		results = append(results, v)
	}

	// We want to process the results in the order of the longest match in a spine first.  That reduces false
	// positives.
	sort.Slice(results, func(i, j int) bool {
		return len(results[i].searchTitle)+len(results[i].searchAuthor) >
			len(results[j].searchTitle)+len(results[j].searchAuthor)
	})

	for _, result := range results {
		log.Printf("Process result %+v", result)
		spines[result.spineindex].Author = result.foundAuthor
		spines[result.spineindex].Title = result.foundTitle
		fragments = flagUsed(fragments, result.spineindex)
		spines, fragments = checkAdjacent(spines, fragments, result)
	}

	return spines, fragments
}

func flagUsed(fragments []OCRFragment, spineindex int) []OCRFragment {
	for i, frag := range fragments {
		if frag.SpineIndex == spineindex {
			log.Printf("Flag used frag %s", frag.Description)
			fragments[i].Used = true
		}
	}

	return fragments
}

func checkAdjacent(spines []Spine, fragments []OCRFragment, result searchResult) ([]Spine, []OCRFragment) {
	// We might have matched on part of a title and have the rest of it in an adjacent spine.  If so it's
	// good to remove it to avoid it causing false matches.
	s := regexp.MustCompile("(?i)" + result.searchTitle)
	residual := strings.TrimSpace(s.ReplaceAllString(result.searchTitle, ""))

	if len(residual) > 0 {
		log.Printf("Residual %s after remove %s", residual, result.searchTitle)

		cmp := []int{}

		if result.spineindex > 0 {
			cmp = append(cmp, result.spineindex-1)
		}

		if result.spineindex < len(spines)-1 {
			cmp = append(cmp, result.spineindex+1)
		}

		s = regexp.MustCompile("(?i)" + residual)

		for _, i := range cmp {
			if len(spines[i].Author) == 0 && s.MatchString(spines[i].Spine) {
				log.Printf("Remove rest of title %s in %s", residual, spines[i].Spine)
				spines[i].Spine = s.ReplaceAllString(spines[i].Spine, "")
			}
		}
	}

	return spines, fragments
}

func setUpPhases() []phase {
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
		if len(spine.Author) > 0 {
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

type SpineOrder struct {
	len   int
	index int
}

func getOrder(spines []Spine) []SpineOrder {
	// Order our search by longest spine first.  This is because the longer the spine is, the more likely
	// it is to have both and author and a subject, and therefore match.  Matching gets it out of the way
	// but also gives us a known author, which can be used to good effect to improve matching on other
	// spines.

	order := []SpineOrder{}

	for i, spine := range spines {
		order = append(order, SpineOrder{
			len:   len(spine.Spine),
			index: i,
		})
	}

	sort.Slice(order, func(i, j int) bool {
		return order[i].len > order[j].len
	})

	return order
}

// We do a whole load of regexp replacement.  Inefficient to create the regexps each time.
type normalizeRegExp struct {
	search  *regexp.Regexp
	replace []byte
}

var normalizeAuthorRegExp []normalizeRegExp
var normalizeTitleRegExp []normalizeRegExp

func normalizeSetupRegexp() {
	if len(normalizeAuthorRegExp) == 0 {
		normalizeAuthorRegExp = []normalizeRegExp{
			{
				// Any numbers in an author are junk.
				search:  regexp.MustCompile(`[0-9]`),
				replace: []byte(""),
			},
			{
				// Remove Dr. as this isn't always present.
				search:  regexp.MustCompile(`Dr\.`),
				replace: []byte(""),
			},
			{
				// Anything in brackets should be removed - not part of the name, could be "(writing as ...)".
				search:  regexp.MustCompile(`(.*)\(.*\)(.*)`),
				replace: []byte("$1$2"),
			},
			{
				// Remove anything which isn't alphabetic.
				search:  regexp.MustCompile(`(?i)[^a-z ]+`),
				replace: []byte(""),
			},
		}

		normalizeTitleRegExp = []normalizeRegExp{
			{
				// Some books have a subtitle, and the catalogues are inconsistent about whether that's included.
				search:  regexp.MustCompile(`(.*?):`),
				replace: []byte("$1"),
			},
			{
				// Anything in brackets should be removed - ditto.
				search:  regexp.MustCompile(`(.*)\(.*\)(.*)`),
				replace: []byte("$1$2"),
			},
			{
				// Remove anything which isn't alphanumeric.
				search:  regexp.MustCompile(`(?i)[^a-z0-9 ]+`),
				replace: []byte(""),
			},
		}
	}
}

func init() {
	log.Printf("Init function")
	normalizeSetupRegexp()
}

func normalizeWithRegExps(str string, list []normalizeRegExp) string {
	for _, re := range list {
		str = string(re.search.ReplaceAll([]byte(str), re.replace))
	}

	str = strings.ToLower(str)

	log.Printf("Normalized to %s", str)

	return str
}

func NormalizeAuthor(author string) string {
	author = normalizeWithRegExps(author, normalizeAuthorRegExp)
	author = strings.TrimSpace(strings.ToLower(author))
	author = removeShortWords(author)

	return author
}

func NormalizeTitle(title string) string {
	title = normalizeWithRegExps(title, normalizeTitleRegExp)
	title = strings.TrimSpace(strings.ToLower(title))
	title = removeShortWords(title)

	return title
}

func searchSpines(spines []Spine, fragments []OCRFragment, phase phase) {
	log.Printf("Search spines %+v", phase)

	order := getOrder(spines)

	if phase.fuzzy {
		// Fuzzy match using DB of known words.  This has the frequency values in it and is therefore likely to
		// yield a better result than what happens within an each ElasticDB search, though we still do that
		// for authors as it works better.
		// TODO
	}

	for _, o := range order {
		// We want to search this spine.  We're hoping it consists of author title, or perhaps title author, but
		// we don't know where the boundary is.  So we want to search breaking at each word.  Use a wait group so that
		// we can do that in parallel.
		spineindex := o.index
		spine := spines[spineindex]

		if len(spine.Author) == 0 {
			// Not yet identified this spine.
			log.Printf("Spine %d %s", spineindex, spine.Spine)
			words := strings.Split(spines[o.index].Spine, " ")
			var wg sync.WaitGroup

			var author, title string

			for wordindex := 0; wordindex+1 < len(words); wordindex++ {
				if phase.authorstart {
					author = strings.Join(words[0:wordindex+1], " ")
					title = strings.Join(words[wordindex+1:len(words)], " ")
					log.Printf("Consider author first split in spine %d at %d %s - %s", spineindex, wordindex, author, title)
				} else {
					title = strings.Join(words[0:wordindex+1], " ")
					author = strings.Join(words[wordindex+1:len(words)], " ")
					log.Printf("Consider author last split in spine %d at %d %s - %s", spineindex, wordindex, author, title)
				}

				wg.Add(1)
				go func(author string, title string, spineindex int, wordindex int) {
					defer wg.Done()
					search(spineindex, author, title, phase.authorplustitle)
				}(author, title, spineindex, wordindex)
			}

			wg.Wait()
		}
	}
}
