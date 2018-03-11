package http

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/golang/glog"
	"github.com/olivere/elastic"
	"github.com/rodrigo-brito/elasticsearch-playground/action/elasticsearch"
	"github.com/spf13/viper"
)

func MultiMatchPrefix(w http.ResponseWriter, r *http.Request) {
	client, err := elasticsearch.GetConnection()
	if err != nil {
		glog.Error(err)
		w.WriteHeader(http.StatusServiceUnavailable)
		return
	}

	term := r.URL.Query().Get("q")
	terms := strings.Split(term, " ")

	termsCount := len(terms)
	var matchQueries []elastic.Query
	for i := 0; i < termsCount-1; i++ {
		match := elastic.NewMultiMatchQuery(terms[i]).
			FieldWithBoost("title", 2).
			FieldWithBoost("theme", 1).
			FieldWithBoost("director", 2).
			FieldWithBoost("title_director", 1).
			Fuzziness("AUTO")
		matchQueries = append(matchQueries, match)
	}

	prefix := elastic.NewQueryStringQuery(terms[termsCount-1]+"~AUTO*").
		FieldWithBoost("title", 2).
		FieldWithBoost("theme", 1).
		FieldWithBoost("director", 2).
		FieldWithBoost("title_director", 1)

	matchQueries = append(matchQueries, prefix)

	source := elastic.NewFetchSourceContext(true).
		Exclude("*_ngram")

	boolQuery := elastic.NewBoolQuery().Must(matchQueries...)

	result, err := client.Search().
		Index(viper.GetString("indexName")).
		Type("movies").
		Query(boolQuery).
		Pretty(true).
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
