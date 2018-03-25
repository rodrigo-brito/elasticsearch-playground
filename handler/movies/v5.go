package movies

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"github.com/golang/glog"
	"github.com/olivere/elastic"
	"github.com/rodrigo-brito/elasticsearch-playground/action/elasticsearch"
	"github.com/spf13/viper"
)

var re *regexp.Regexp

func init() {
	var err error
	re, err = regexp.Compile(`[\+\-\=&\|\>\<\!\(\)\{\}\[\]\^"~\*\?\:/]`)
	if err != nil {
		glog.Error(err)
	}
}

func scapeQueryString(term string) string {
	return re.ReplaceAllString(term, "")
}

func injectFuzzySufix(term string) string {
	terms := strings.Split(strings.TrimSpace(scapeQueryString(term)), " ")
	return fmt.Sprintf("%s~AUTO*", strings.Join(terms, "~AUTO "))
}

func PrefixPhraseNgran(w http.ResponseWriter, r *http.Request) {
	client, err := elasticsearch.GetConnection()
	if err != nil {
		glog.Error(err)
		w.WriteHeader(http.StatusServiceUnavailable)
		return
	}

	const minimumMatch = "80%"

	term := r.URL.Query().Get("q")

	prefixExact := elastic.NewQueryStringQuery(fmt.Sprintf("%s*", scapeQueryString(term))).
		MinimumShouldMatch(minimumMatch).
		AnalyzeWildcard(true).
		FieldWithBoost("title_director_ngram", 1).
		Analyzer("standard")

	prefixWithFuzzy := elastic.NewQueryStringQuery(injectFuzzySufix(term)).
		MinimumShouldMatch(minimumMatch).
		AnalyzeWildcard(true).
		FieldWithBoost("title_director_ngram", 1).
		Analyzer("standard")

	query := elastic.NewDisMaxQuery().
		Query(
			prefixExact.Boost(3),
			prefixWithFuzzy.Boost(1),
		).TieBreaker(0.1)

	source := elastic.NewFetchSourceContext(true).
		Exclude("*_ngram")

	result, err := client.Search().
		Index(viper.GetString("indexName")).
		Type("movies").
		Query(query).
		Pretty(true).
		TrackScores(true).
		FetchSourceContext(source).
		Do(r.Context())

	if err != nil {
		w.WriteHeader(http.StatusServiceUnavailable)
		glog.Error(err)
		return
	}

	json, err := json.Marshal(result.Hits)
	if err != nil {
		w.WriteHeader(http.StatusServiceUnavailable)
		glog.Error(err)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(json)
}
