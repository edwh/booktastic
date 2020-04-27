package main

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

// Queries are executed using channels so that we can perform them in parallel
//func SearchAuthorTitle(query chan<-ElasticQuery) ElasticResult {
//	// Get our client.
//	//es, _ := elasticsearch.NewDefaultClient()
//	//
//	//query
//}
