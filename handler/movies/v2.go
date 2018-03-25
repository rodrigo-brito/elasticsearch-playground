package movies

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/golang/glog"
	"github.com/olivere/elastic"
	"github.com/rodrigo-brito/elasticsearch-playground/action/elasticsearch"
	"github.com/spf13/viper"
)

func MultiMatchNgran(w http.ResponseWriter, r *http.Request) {
	client, err := elasticsearch.GetConnection()
	if err != nil {
		glog.Error(err)
		w.WriteHeader(http.StatusServiceUnavailable)
		return
	}

	term := r.URL.Query().Get("q")
	terms := strings.Split(term, " ")

	var matchQueries []elastic.Query
	for _, term := range terms {
		match := elastic.NewMultiMatchQuery(term).
			FieldWithBoost("title", 2).
			FieldWithBoost("title_ngram", 2).
			FieldWithBoost("theme", 1).
			FieldWithBoost("director", 1).
			FieldWithBoost("director_ngram", 1).
			FieldWithBoost("title_director", 1).
			FieldWithBoost("title_director_ngram", 1).
			Fuzziness("AUTO")
		matchQueries = append(matchQueries, match)
	}
	disMax := elastic.NewDisMaxQuery().Query(matchQueries...).TieBreaker(0)

	result, err := client.Search().
		Index(viper.GetString("indexName")).
		Type("movies").
		Query(disMax).
		Pretty(true).
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
