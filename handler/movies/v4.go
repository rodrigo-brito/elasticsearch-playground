package movies

import (
	"encoding/json"
	"net/http"

	"github.com/golang/glog"
	"github.com/olivere/elastic"
	"github.com/rodrigo-brito/elasticsearch-playground/action/elasticsearch"
	"github.com/spf13/viper"
)

func MultiMatchPrefixShingle(w http.ResponseWriter, r *http.Request) {
	client, err := elasticsearch.GetConnection()
	if err != nil {
		glog.Error(err)
		w.WriteHeader(http.StatusServiceUnavailable)
		return
	}

	term := r.URL.Query().Get("q")

	query := elastic.NewQueryStringQuery(getFuzzyTerm(term)).
		FieldWithBoost("title_shingle", 2).
		FieldWithBoost("theme", 1).
		FieldWithBoost("director_shingle", 1).
		AnalyzeWildcard(true).
		DefaultOperator("AND").
		UseDisMax(true)

	result, err := client.Search().
		Index(viper.GetString("indexName")).
		Type("movies").
		Query(query).
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
