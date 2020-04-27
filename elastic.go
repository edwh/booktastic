package main

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/elastic/go-elasticsearch/v7"
	"log"
	"regexp"
	"strings"
)

const INDEX = "booktastic"

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

func NormalizeAuthor(author string) string {
	// TODO Wasteful to build regex each time.

	// Any numbers in an author are junk.
	author = regexp.MustCompile(`[0-9]`).ReplaceAllString(author, "")

	// Remove Dr. as this isn't always present.
	author = regexp.MustCompile(`Dr.`).ReplaceAllString(author, "")

	// Anything in brackets should be removed - not part of the name, could be "(writing as ...)".
	author = regexp.MustCompile(`(.*)\(.*\)(.*)`).ReplaceAllString(author, "$1$2")

	// Remove anything which isn't alphabetic.
	author = regexp.MustCompile(`(?i)[^a-z ]+`).ReplaceAllString(author, "")

	author = strings.TrimSpace(strings.ToLower(author))

	author = removeShortWords(author)

	return author
}

func NormalizeTitle(title string) string {
	// Some books have a subtitle, and the catalogues are inconsistent about whether that's included.
	title = regexp.MustCompile(`(.*?):`).ReplaceAllString(title, "$1")

	// Anything in brackets should be removed - ditto.
	title = regexp.MustCompile(`(.*)\(.*\)(.*)`).ReplaceAllString(title, "$1$2")

	// Remove anything which isn't alphanumeric.
	title = regexp.MustCompile(`(?i)[^a-z0-9 ]+`).ReplaceAllString(title, "")

	title = strings.TrimSpace(strings.ToLower(title))

	title = removeShortWords(title)

	return title
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

// Queries are executed using channels so that we can perform them in parallel
func SearchAuthorTitle(author string, title string) {
	if len(author) > 0 && len(title) > 0 {
		// Get our client.  We need a separate one because we're very parallelised here.
		es, _ := elasticsearch.NewDefaultClient()

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
								"normalauthor": map[string]interface{}{
									"value":     author,
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
			es.Search.WithPretty(),
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

		var r map[string]interface{}

		if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
			log.Fatalf("Error parsing the response body: %s", err)
		}
		// Print the response status, number of results, and request duration.
		log.Printf(
			"Search for %s - %s [%s] %d hits; took: %dms",
			author,
			title,
			res.Status(),
			int(r["hits"].(map[string]interface{})["total"].(map[string]interface{})["value"].(float64)),
			int(r["took"].(float64)),
		)
		// Print the ID and document source for each hit.
		for _, hit := range r["hits"].(map[string]interface{})["hits"].([]interface{}) {
			log.Printf(" * ID=%s, %s", hit.(map[string]interface{})["_id"], hit.(map[string]interface{})["_source"])
		}

		log.Println(strings.Repeat("=", 37))
	}
}

func search(author string, title string, authorplustitle bool) {
	author2 := NormalizeAuthor(author)

	authwords := strings.Split(author2, " ")

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

	author = author2

	// There are some titles which are very short, but they are more likely to just be false junk.
	if len(title) < 4 {
		log.Printf("Reject too short title %s", title)
		return
	}

	title = NormalizeTitle(title)

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
