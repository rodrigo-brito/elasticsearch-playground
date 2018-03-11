package http

import (
	"encoding/json"
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
	re, err = regexp.Compile(`\+\-\=&\|\>\<\!\(\)\{\}\[\]\^"~\*\?\:/`)
	if err != nil {
		glog.Error(err)
	}
}
func scapeQueryString(term string) string {
	return re.ReplaceAllString(term, "")
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

	terms := strings.Split(term, " ")
	var queries []elastic.Query
	for _, term := range terms[:len(terms)-1] {
		queries = append(queries, elastic.NewMultiMatchQuery(term).
			FieldWithBoost("title_director_ngram", 2).
			Fuzziness("AUTO").
			MaxExpansions(20).
			MinimumShouldMatch(minimumMatch))
	}

	prefixQuery := elastic.NewQueryStringQuery(
		scapeQueryString(terms[len(terms)-1])+"~AUTO").
		MinimumShouldMatch(minimumMatch).
		AnalyzeWildcard(true).
		FieldWithBoost("title_director_ngram", 1)
	queries = append(queries, prefixQuery)

	boolQuery := elastic.NewBoolQuery().Must(
		queries...,
	)

	result, err := client.Search().
		Index(viper.GetString("indexName")).
		Type("movies").
		Query(boolQuery).
		Pretty(true).
		TrackScores(true).
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
