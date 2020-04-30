package main

import (
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
	adjacent        bool
	permuted        bool
	mangled         bool
}

const MAXRESULTS = 1000

type searchResult struct {
	spineindex   int
	phaseid      int
	searchAuthor string
	searchTitle  string
	foundAuthor  string
	foundTitle   string
}

var searchResults map[searchResult]searchResult
var gotSpine map[int]bool
var resultsMux sync.Mutex

func addResult(result searchResult) {
	sugar.Debugf("Add result %+v", result)
	// We want the results to be unique by found values in a single spine.  We might have searched on different
	// variants but found the same thing.
	key := result
	key.searchTitle = ""
	key.searchAuthor = ""

	resultsMux.Lock()

	if !gotSpine[result.spineindex] {
		searchResults[key] = result
		sugar.Infof("Phase %d found result on %d, %s - %s", result.phaseid, result.spineindex, result.foundAuthor, result.foundTitle)
		gotSpine[result.spineindex] = true
	}

	resultsMux.Unlock()
}

func clearResults() {
	searchResults = map[searchResult]searchResult{}
	gotSpine = map[int]bool{}
}

func checkResult(spineindex int) bool {
	resultsMux.Lock()
	ret := gotSpine[spineindex]
	resultsMux.Unlock()
	//sugar.Infof("Check result on %d, %t", spineindex, ret)
	return ret
}

func IdentifyBooks(spines []Spine, fragments []OCRFragment) ([]Spine, []OCRFragment) {
	phases := setUpPhases()

	// We need to execute the phases serially as the results of one phase make it more likely that we can find things
	// in later phases.
	for _, p := range phases {
		sugar.Infof("Execute phase %+v", p)
		sugar.Debugf("Spines at start of phase %+v", spines)
		start := time.Now()

		clearResults()
		searchSpines(spines, fragments, p, 0, len(spines))
		sugar.Debugf("Spines after search phase %+v", spines)
		spines, fragments = processSearchResults(spines, fragments)
		sugar.Debugf("Spines after process %+v", spines)

		if p.mangled || p.permuted || p.adjacent {
			clearResults()
			spines, fragments = searchBrokenSpines(spines, fragments, p)
			spines, fragments = processSearchResults(spines, fragments)
			sugar.Debugf("Spines after broken %+v", spines)
		}

		duration := time.Since(start)
		sugar.Debugf("Phase %d %+v found %d in %v", p.id, p, len(searchResults), duration)
	}

	sugar.Debugf("All phases complete")

	for _, frag := range fragments {
		if !frag.Used {
			sugar.Debugf("LEFTOVER: spine %d %s", frag.SpineIndex, frag.Description)
		}
	}

	for _, spine := range spines {
		if len(spine.Author) > 0 {
			sugar.Debugf("RESULT: %s - %s", spine.Author, spine.Title)
		}
	}

	return spines, fragments
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

	// First record the results.  Need to do this before the next bit as the results contain a spine index which
	// may later shift.
	for _, result := range results {
		sugar.Debugf("Process result %+v", result)
		spines[result.spineindex].Author = result.foundAuthor
		spines[result.spineindex].Title = result.foundTitle
		fragments = flagUsed(fragments, result.spineindex)
		spines, fragments = checkAdjacent(spines, fragments, result)
	}

	// Spines may change at this point.
	spines, fragments = extractKnownAuthors(spines, fragments)

	return spines, fragments
}

func flagUsed(fragments []OCRFragment, spineindex int) []OCRFragment {
	for i, frag := range fragments {
		if frag.SpineIndex == spineindex {
			sugar.Debugf("Flag used frag %s", frag.Description)
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
		sugar.Debugf("Residual %s after remove %s", residual, result.searchTitle)

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
				sugar.Debugf("Remove rest of title %s in %s", residual, spines[i].Spine)
				spines[i].Spine = s.ReplaceAllString(spines[i].Spine, "")
			}
		}
	}

	return spines, fragments
}

func extractKnownAuthors(spines []Spine, fragments []OCRFragment) ([]Spine, []OCRFragment) {
	// People often file books from the same author together.  If we check the authors we have in hand so far
	// then we can ensure that no known author is split across multiple spines.  That can happen sometimes in
	// the Google results.  This means that we will find the author when we are checking broken spines.
	//
	// It also avoids some issues where we can end up using the author from the wrong spine because the correct
	// author is split across more spines than we are currently looking at.
	amap := map[string]bool{}

	for _, spine := range spines {
		if len(spine.Author) > 0 {
			amap[spine.Author] = true
		}
	}

	sugar.Debugf("Currently known authors %+v", amap)

	for author := range amap {
		sugar.Debugf("Check author %s", author)
		authorwords := strings.Split(author, " ")
		wordindex := 0

		for spineindex, spine := range spines {
			if len(spine.Author) == 0 {
				// Not dealt with this spine.
				wi := wordindex
				si := spineindex
				sugar.Debugf("Start at spine %d currently %s", spineindex, spines[spineindex].Spine)

				spinewords := strings.Split(spine.Spine, " ")

				// We only want to start checking at the start of a spine.  If there are other words earlier in the
				// spine they may be a title and by merging with the next spine we might combine two titles.
				swi := 0

				for ok := true; ok; {
					if (si < len(spines)) &&
						(compare(spinewords[swi], authorwords[wi]) >= CONFIDENCE) {
						sugar.Debugf("Found possible author match %s from %s in %s at %d", authorwords[wi], author, spine.Spine, si)
						wi++
						swi++

						if wi >= len(authorwords) {
							ok = false
						} else if swi >= len(spinewords) {
							swi = 0
							si++

							if len(spines[si].Author) > 0 {
								// Wouldn't be safe to merge with a spine that's already matched.
								ok = false
							}

							spinewords = strings.Split(spines[si].Spine, " ")
						}
					} else {
						ok = false
					}
				}

				if wi >= len(authorwords) && si > spineindex {
					// Found author split across spines.
					// TODO This was si >= in PHP.
					sugar.Debugf("Found end of author in spine %d vs %d spine upto word %d", si, spineindex, swi)
					sugar.Debugf("Spines before %+v", spines)
					sugar.Debugf("Merge at %d len %d - %d", spineindex, si, spineindex)
					comspined := spines[spineindex]
					spines, fragments = mergeSpines(spines, fragments, comspined, spineindex, si-spineindex+1)
					sugar.Debugf("Spines now %+v", spines)
				}
			}
		}
	}

	return spines, fragments
}

func removeEmptySpines(spines []Spine, fragments []OCRFragment) ([]Spine, []OCRFragment) {
	for spineindex, spine := range spines {
		if len(strings.TrimSpace(spine.Spine)) == 0 {
			sugar.Debugf("Remove empty spine %d", spineindex)
		}
	}

	return spines, fragments
}

func mergeSpines(spines []Spine, fragments []OCRFragment, comspined Spine, start int, length int) ([]Spine, []OCRFragment) {
	// We have combined multiple adjacent spines into a single one, possibly with some
	// reordering of text.
	sugar.Debugf("Spines before merge at %d len %d %+v", start, length, spines)

	spines[start] = comspined

	// Renumber the spine indexes in the fragments for the spines which we are (re)moving.
	for fragindex, frag := range fragments {
		if frag.SpineIndex > start && frag.SpineIndex <= start+length-1 {
			// These are the ones we're merging.
			fragments[fragindex].SpineIndex = start
		} else if frag.SpineIndex > start+length-1 {
			// These are above
			fragments[fragindex].SpineIndex -= length - 1
		}
	}

	// Remove.
	spines = append(spines[0:start+1], spines[(start+1+length):]...)
	sugar.Debugf("Spines after merge %+v", spines)

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
	ints := [2]int{0, 1}
	id := 0

	// This is empirical - the phases that find stuff.
	gooduns := map[int]bool{
		0:  true,
		1:  true,
		2:  true,
		3:  true,
		4:  true,
		14: true,
	}

	for _, mangled := range ints {
		for _, permuted := range ints {
			for _, adjacent := range ints {
				for _, fuzzy := range bools {
					// Only one of these in any one phase
					if mangled+permuted+adjacent == 1 {
						for _, authorplustitle := range bools {
							for _, authorstart := range bools {
								if gooduns[id] || true {
									phases = append(phases, phase{
										id,
										fuzzy,
										!authorstart, // Search for start in earlier phase as more common.
										authorplustitle,
										adjacent == 1,
										permuted == 1,
										mangled == 1,
									})
								}

								id++
							}
						}
					}
				}
			}
		}
	}

	sugar.Debugf("Phases %+v", phases)

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

type SpineOrder struct {
	len   int
	index int
}

func getOrder(spines []Spine, start int, length int) []SpineOrder {
	// Order our search by longest spine first.  This is because the longer the spine is, the more likely
	// it is to have both and author and a subject, and therefore match.  Matching gets it out of the way
	// but also gives us a known author, which can be used to good effect to improve matching on other
	// spines.

	sugar.Debugf("Get order")
	order := []SpineOrder{}

	for i, spine := range spines {
		if i >= start && i < start+length {
			order = append(order, SpineOrder{
				len:   len(spine.Spine),
				index: i,
			})
		}
	}

	if length > 1 {
		sort.Slice(order, func(i, j int) bool {
			return order[i].len > order[j].len
		})
	}

	sugar.Debugf("Return order %+v", order)

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
	normalizeSetupRegexp()
}

func normalizeWithRegExps(str string, list []normalizeRegExp) string {
	for _, re := range list {
		str = string(re.search.ReplaceAll([]byte(str), re.replace))
	}

	str = strings.ToLower(str)

	sugar.Debugf("Normalized to %s", str)

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

func searchSpines(spines []Spine, fragments []OCRFragment, phase phase, start int, length int) {
	sugar.Debugf("Search spines start %d len %d phase %+v", start, length, phase)

	order := getOrder(spines, start, length)

	if phase.fuzzy {
		// Fuzzy match using DB of known words.  This has the frequency values in it and is therefore likely to
		// yield a better result than what happens within an each ElasticDB search, though we still do that
		// for authors as it works better.
		// TODO
	}

	type searchEntry struct {
		phaseid    int
		spineindex int
		author     string
		title      string
		wordindex  int
	}

	var wg sync.WaitGroup
	searches := []searchEntry{}

	for _, o := range order {
		// We want to search this spine.  We're hoping it consists of author title, or perhaps title author, but
		// we don't know where the boundary is.  So we want to search breaking at each word.  Use a wait group so that
		// we can do that in parallel.
		spineindex := o.index
		spine := spines[spineindex]

		if len(spine.Author) == 0 {
			// Not yet identified this spine.
			sugar.Debugf("Spine %d %s", spineindex, spine.Spine)
			words := strings.Split(spines[o.index].Spine, " ")

			// If it doesn't have two words, it can't have an author and a title.
			if len(words) >= 2 {
				var author, title string

				wordorder := make([]int, len(words)-1)

				for i := range wordorder {
					wordorder[i] = i
				}

				for _, wordindex := range wordorder {
					if phase.authorstart {
						author = strings.Join(words[0:wordindex+1], " ")
						title = strings.Join(words[wordindex+1:len(words)], " ")
						sugar.Debugf("Consider author first split in spine %d at %d %s - %s", spineindex, wordindex, author, title)
					} else {
						title = strings.Join(words[0:wordindex+1], " ")
						author = strings.Join(words[wordindex+1:len(words)], " ")
						sugar.Debugf("Consider author last split in spine %d at %d %s - %s", spineindex, wordindex, author, title)
					}

					wg.Add(1)
					searches = append(searches, searchEntry{
						phaseid:    phase.id,
						spineindex: spineindex,
						author:     author,
						title:      title,
						wordindex:  wordindex,
					})
				}
			}
		}
	}

	// We're playing a balancing game - if we find a result early then we can save on other searches.  So we don't
	// want to do all searches for the same spine simultanously.  Sorting by word order means we are less likely.
	sort.Slice(searches, func(i, j int) bool {
		if searches[i].wordindex != searches[j].wordindex {
			return searches[i].wordindex > searches[j].wordindex
		} else {
			return searches[i].spineindex > searches[j].spineindex
		}
	})

	// Now fire them all off simultaneously.
	for _, s := range searches {
		go func(author string, title string, spineindex int, wordindex int, phaseid int) {
			defer wg.Done()

			// By the time this gets invoked, it's possible that someone else has identified this spine.
			// If so, no point in us searching too.  There's still a timing window where two can identify
			// at the same time, but that's ok - this is just a speedup.
			if !checkResult(spineindex) {
				sugar.Debugf("Now search %s - %s", author, title)
				search(spineindex, author, title, phase.authorplustitle, phaseid)
			} else {
				sugar.Infof("Already identified %d, skip search", spineindex)
			}
		}(s.author, s.title, s.spineindex, s.wordindex, s.phaseid)
	}

	wg.Wait()
}

func searchBrokenSpines(spines []Spine, fragments []OCRFragment, phase phase) ([]Spine, []OCRFragment) {
	// Up to this point we've relied on what Google returns on a single line.  We will have found
	// some books via that route.  But it's common to have the author on one line, and the book on another,
	// or other variations which result in text on a single spine being split.
	//
	// Ideally we'd search all permutations of all the words.  But this is expensive, so we can only go up so far.
	sugar.Debugf("Search broken spines %+v", phase)

	// Mangled spines are slower but they have a separate limit on the number of words.
	max := 4

	for adjacent := 2; adjacent <= max; adjacent++ {
		spineindex := 0

		for ok := true; ok; {
			thisone := spines[spineindex]

			if len(spines[spineindex].Author) == 0 && len(thisone.Spine) > 0 {
				sugar.Debugf("Consider broken spine %s at %d length %d", thisone.Spine, spineindex, adjacent)

				available := true
				healedtext := ""

				for next := spineindex; next < len(spines) && next-spineindex+1 <= adjacent; next++ {
					if len(spines[next].Author) > 0 {
						available = false
					} else {
						healedtext += " " + spines[next].Spine
					}
				}

				if available {
					sugar.Debugf("Available")

					var found bool

					if phase.adjacent {
						spines, fragments, found = searchForAdjacentSpines(spines, fragments, spineindex, adjacent, phase)
					} else if phase.permuted {
						spines, fragments, found = searchForPermutedSpines(spines, fragments, spineindex, adjacent, phase)
					} else if phase.mangled {
						// TODO
						//spines, fragments, found = searchForMangledSpines(spines, fragments, spineindex, adjacent, phase)
					}

					if found {
						// Stop now as our spine index has become invalid.
						// TODO Maybe we could keep going and find more?
						ok = false
					}
				}
			}

			spineindex++

			if spineindex+adjacent >= len(spines) {
				ok = false
			}
		}
	}

	return spines, fragments
}

func permutations(arr []int) [][]int {
	var helper func([]int, int)
	res := [][]int{}

	helper = func(arr []int, n int) {
		if n == 1 {
			tmp := make([]int, len(arr))
			copy(tmp, arr)
			res = append(res, tmp)
		} else {
			for i := 0; i < n; i++ {
				helper(arr, n-1)
				if n%2 == 1 {
					tmp := arr[i]
					arr[i] = arr[n-1]
					arr[n-1] = tmp
				} else {
					tmp := arr[0]
					arr[0] = arr[n-1]
					arr[n-1] = tmp
				}
			}
		}
	}
	helper(arr, len(arr))
	return res
}

func searchForAdjacentSpines(spines []Spine, fragments []OCRFragment, start int, length int, phase phase) ([]Spine, []OCRFragment, bool) {
	// We want to join adjacent spines to see if we can find a match that way.
	//
	// Because we search on normalised values, which include sorting, then by joining a spine we automatically get
	// some permutation searching on that joined value.  For example if we join A B and C D E to get A B C D E, then
	// searching on A B C D E will search on:
	//
	// A - B C D E
	// A B - C D E
	// A B C - D E
	// A B C D - E
	//
	// Just looking at the left hand side, that's the same as searching
	//
	// A
	// A B
	// B A
	// A B C
	// B A C
	// A C B
	// B C A
	// C A B
	// C B A
	// ...etc
	//
	// This doesn't give us permutations such as A E - B C D - for that we need to permute the order in which we
	// join the spines, which is searchPermutedSpines
	//
	// We can't really parallelise this easily since we are looking at multiple spines.  So process
	// the results as they arrive.
	// TODO Bet we could, though.
	found := false

	sugar.Debugf("Before adjacent merge %+v", spines)
	healed := ""

	for spineindex := start; spineindex < start+length; spineindex++ {
		healed += " " + spines[spineindex].Spine
	}

	sugar.Debugf("Merge adjacent text %s", healed)
	newspines := make([]Spine, len(spines))
	copy(newspines, spines)
	newspines[start].Spine = healed

	newspines, newfragments := mergeSpines(newspines, fragments, newspines[start], start, length-1)
	sugar.Debugf("After adjacent merge %+v", newspines)

	// Search using this set of spines to see if we find something.
	clearResults()
	searchSpines(newspines, newfragments, phase, start, 1)

	if len(searchResults) > 0 {
		// We found something.  Use this set.
		sugar.Debugf("Found an adjacent result - use this set")
		sugar.Debugf("Spines before %+v", spines)
		copy(spines, newspines)
		copy(fragments, newfragments)
		sugar.Debugf("Spines after %+v", spines)
		spines, fragments = processSearchResults(spines, fragments)
		sugar.Debugf("Spines process %+v", spines)
		found = true
	}

	return spines, fragments, found
}

func searchForPermutedSpines(spines []Spine, fragments []OCRFragment, start int, length int, phase phase) ([]Spine, []OCRFragment, bool) {
	// We want to join adjacent spines to see if we can find a match that way.
	//
	// Because we search on normalised values, which include sorting, then by joining a spine we automatically get
	// some permutation searching on that joined value.  For example if we join A B and C D E to get A B C D E, then
	// searching on A B C D E will search on:
	//
	// A - B C D E
	// A B - C D E
	// A B C - D E
	// A B C D - E
	//
	// Just looking at the left hand side, that's the same as searching
	//
	// A
	// A B
	// B A
	// A B C
	// B A C
	// A C B
	// B C A
	// C A B
	// C B A
	// ...etc
	//
	// This doesn't give us permutations such as A E - B C D - for that we'd need to permute the spines rather
	// than join adjacent ones.
	//
	// We can't really parallelise this easily since we are looking at multiple spines.  So process
	// the results as they arrive.
	// TODO Bet we could, though.
	found := false
	seq := make([]int, length)

	for i := range seq {
		seq[i] = i + start
	}

	orders := permutations(seq)
	done := false

	for _, order := range orders {
		if !done {
			// Generate a set of spines and fragments which match this permutation.
			sugar.Debugf("Consider permutation %+v", order)
			healedtext := ""

			for _, ent := range order {
				healedtext += " " + spines[ent].Spine
			}

			sugar.Debugf("Healed text %s", healedtext)

			// Now use this text and drop the others.  If we find results then these spines and fragments will become our
			// actual versions.
			//
			// Need to clone as slices are passed by reference (effectively).
			newspines := make([]Spine, len(spines))
			newfragments := make([]OCRFragment, len(fragments))
			copy(newspines, spines)
			copy(newfragments, fragments)
			newspines[start].Spine = healedtext
			newspines, newfragments = mergeSpines(newspines, newfragments, newspines[start], start, length-1)

			// Search using this set of spines to see if we find something.
			clearResults()
			searchSpines(newspines, newfragments, phase, start, 1)

			if len(searchResults) > 0 {
				// We found something.  Use this set.
				// TODO Different permutations might find better results than others?
				sugar.Debugf("Found a permuted result - use this set")
				sugar.Debugf("Spines before %+v", spines)
				copy(spines, newspines)
				copy(fragments, newfragments)
				sugar.Debugf("Spines after %+v", spines)
				spines, fragments = processSearchResults(spines, fragments)
				sugar.Debugf("Spines process %+v", spines)
				done = true
				found = true
			}
		}
	}

	return spines, fragments, found
}

func searchForMangledSpines(spines []Spine, fragments []OCRFragment, start int, length int, phase phase) ([]Spine, []OCRFragment, bool) {
	// More rarely Google breaks in the wrong place.  So now we want to ignore where Google has broken.  Split each of
	// these spines into single-word spines so that we can then permute them.
	//
	// Need to clone as slices are passed by reference (effectively).
	sugar.Debugf("Mangled spine pre-split into words %+v", spines)
	found := false
	newspines := []Spine{}

	if start > 0 {
		newspines = make([]Spine, start)
		copy(newspines, spines[0:start])
	}

	for spineindex := start; spineindex < start+length; spineindex++ {
		spinewords := strings.Split(spines[spineindex].Spine, " ")

		for _, word := range spinewords {
			// Add a new spine for this word.
			sugar.Debugf("Add spine for %s", word)
			newspines = append(newspines, Spine{
				Spine:  word,
				Author: "",
				Title:  "",
			})
		}
	}

	if start+length < len(spines) {
		t := []Spine{}
		t = append(t, newspines...)
		t = append(t, spines[(start+length):]...)
		newspines = t
	}

	added := len(newspines) - len(spines)

	// We can only afford to search permutations upto a certain length, as it's factorial.
	if added+length < 5 {
		// We've added some spines in this process.  We need to renumber the fragments, which is a bit of a faff.
		newlines := []string{}
		for _, spine := range newspines {
			newlines = append(newlines, spine.Spine)
		}

		newfragments := make([]OCRFragment, len(fragments))
		copy(newfragments, fragments)
		newfragments = AddSpineIndex(newlines, newfragments)

		sugar.Debugf("Mangled spine post-split into words added %d spines %+v", added, newspines)

		// Now we search for permutations of these new spines.
		sugar.Infof("Mangled spine search phase %d start %d length %d added %d - start", phase.id, start, length, added)
		newspines, newfragments, found := searchForPermutedSpines(newspines, newfragments, start, length+added, phase)
		sugar.Infof("Mangled spine search phase %d start %d length %d added %d - end %t", phase.id, start, length, added, found)

		if found {
			sugar.Debugf("Found something in mangled spine search")
			spines = newspines
			fragments = newfragments
		} else {
			sugar.Debugf("Nothing found in mangled spine search")
		}
	}

	return spines, fragments, found
}
