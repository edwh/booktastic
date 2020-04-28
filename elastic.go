package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/agnivade/levenshtein"
	"github.com/elastic/go-elasticsearch/v7"
	"github.com/patrickmn/go-cache"
	"log"
	"strings"
	"time"
)

const INDEX = "booktastic"
const CONFIDENCE = 75

type ElasticQuery struct {
	Author string
	Title  string
}

type ElasticResult struct {
	Author       string
	Title        string
	NormalAuthor string
	NormalTitle  string
}

// We use a cache to reduce searches, as our parallelisation can often result in the same combinations.
var elasticCache *cache.Cache = nil

func init() {
	log.Printf("Create cache")
}

func removeShortWords(str string) string {
	words := strings.Split(strings.TrimSpace(str), " ")
	ret := []string{}

	for _, word := range words {
		word = strings.TrimSpace(word)

		if len(word) > 3 {
			ret = append(ret, word)
		}
	}

	return strings.Join(ret, " ")
}

func getElastic() *elasticsearch.Client {
	// Get our client.  We need a separate one because we're very parallelised here.
	es, _ := elasticsearch.NewDefaultClient()

	return es
}

func getCache() {
	if elasticCache == nil {
		log.Printf("Create cache")
		elasticCache = cache.New(cache.NoExpiration, 10*time.Minute)
	}
}

// Queries are executed using channels so that we can perform them in parallel
func SearchAuthorTitle(author string, title string) {
	if len(author) > 0 && len(title) > 0 {
		getCache()

		// See if we have an entry cached which will save the query.
		var r map[string]interface{}

		if x, found := elasticCache.Get(author + "-" + title); found {
			log.Printf("Found cache entry %s-%s", author, title)
			r = x.(map[string]interface{})
		} else {
			// No cache entry - query.
			es := getElastic()

			// Empirical testing shows that using a fuzziness of 2 for author all the time gives good results.
			var buf bytes.Buffer
			query := map[string]interface{}{
				"query": map[string]interface{}{
					"bool": map[string]interface{}{
						"must": []interface{}{
							map[string]interface{}{
								"fuzzy": map[string]interface{}{
									"normalauthor": map[string]interface{}{
										"value":     author,
										"fuzziness": 2,
									},
								},
							},
							map[string]interface{}{
								"fuzzy": map[string]interface{}{
									"normaltitle": map[string]interface{}{
										"value":     title,
										"fuzziness": 2,
									},
								},
							},
						},
					},
				},
			}

			if err := json.NewEncoder(&buf).Encode(query); err != nil {
				log.Fatalf("Error encoding query: %s", err)
			}

			// Perform the search request.
			res, err := es.Search(
				es.Search.WithContext(context.Background()),
				es.Search.WithIndex(INDEX),
				es.Search.WithBody(&buf),
				//es.Search.WithPretty(),
				es.Search.WithSize(5),
			)

			if err != nil {
				log.Fatalf("Error getting response: %s", err)
			}

			defer res.Body.Close()

			if res.IsError() {
				var e map[string]interface{}
				if err := json.NewDecoder(res.Body).Decode(&e); err != nil {
					log.Fatalf("Error parsing the response body: %s", err)
				} else {
					// Print the response status and error information.
					log.Fatalf("[%s] %s: %s",
						res.Status(),
						e["error"].(map[string]interface{})["type"],
						e["error"].(map[string]interface{})["reason"],
					)
				}
			}

			if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
				log.Fatalf("Error parsing the response body: %s", err)
			}

			// Save in cache for next time.
			elasticCache.Set(author+"-"+title, r, cache.NoExpiration)
		}

		for _, hit := range r["hits"].(map[string]interface{})["hits"].([]interface{}) {
			log.Printf(" * ID=%s, %s", hit.(map[string]interface{})["_id"], hit.(map[string]interface{})["_source"])
			data := hit.(map[string]interface{})["_source"]
			hitauthor := fmt.Sprintf("%v", data.(map[string]interface{})["normalauthor"])
			hittitle := fmt.Sprintf("%v", data.(map[string]interface{})["normaltitle"])

			if len(hitauthor) > 0 && len(hittitle) > 0 {
				authperc := compare(author, hitauthor)
				titperc := compare(title, hittitle)

				log.Printf("Author + title match %d, %d, %s - %s vs %s - %s", authperc, titperc, author, title, hitauthor, hittitle)
				if authperc >= CONFIDENCE && titperc >= CONFIDENCE && sanityCheck(hitauthor, hittitle) {
					log.Printf("FOUND: Author + Title match %d, %d %+v", authperc, titperc, data)
				}
			}
		}

		log.Println(strings.Repeat("=", 37))
	}
}

func sanityCheck(author, title string) bool {
	// We see some matches where the author and title are basically the same.  Might be true for autobiographies but
	// more likely junk.
	if strings.Contains(author, title) || strings.Contains(title, author) {
		return false
	}

	return true
}

func compare(str1, str2 string) int {
	len1 := len(str1)
	len2 := len(str2)

	var lenratio float32
	lenratio = float32(len1) / float32(len2)

	var pc int

	if strings.Contains(str1, str2) || strings.Contains(str2, str1) &&
		lenratio >= 0.5 && lenratio <= 2 {
		// One inside the other is pretty good as long as they're not too different in length.
		pc = CONFIDENCE
	} else {
		dist := levenshtein.ComputeDistance(str1, str2)

		var max int

		if len1 > len2 {
			max = len1
		} else {
			max = len2
		}

		pc = 100 - 100*dist/max
	}

	return pc
}

func search(author string, title string, authorplustitle bool) {
	authwords := strings.Split(author, " ")

	// Require an author to have one part of their name which isn't very short.  Probably discriminates against
	// Chinese people who use initials, so not ideal.
	longenough := false

	for _, word := range authwords {
		if len(word) > 3 {
			longenough = true
		}
	}

	if !longenough {
		log.Printf("Reject too short author %s", author)
		return
	}

	// There are some titles which are very short, but they are more likely to just be false junk.
	if len(title) < 4 {
		log.Printf("Reject too short title %s", title)
		return
	}

	log.Printf("Search for %s - %s", author, title)

	if authorplustitle {
		log.Printf("author - title")
		SearchAuthorTitle(author, title)
	} else {
		log.Printf("author only")
		// TODO

		log.Printf("title only")
		// TODO
	}
}
