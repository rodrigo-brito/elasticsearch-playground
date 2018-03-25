package movies

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/golang/glog"
	"github.com/olivere/elastic"
	"github.com/rodrigo-brito/elasticsearch-playground/action/elasticsearch"
	"github.com/spf13/viper"
)

func getFuzzyTerm(term string) string {
	return fmt.Sprintf("%s~AUTO", strings.Join(strings.Split(term, " "), "~AUTO "))
}

func QueryStringWithSlplit(w http.ResponseWriter, r *http.Request) {
	client, err := elasticsearch.GetConnection()
	if err != nil {
		glog.Error(err)
		w.WriteHeader(http.StatusServiceUnavailable)
		return
	}

	term := r.URL.Query().Get("q")

	query := elastic.NewQueryStringQuery(getFuzzyTerm(term)).
		FieldWithBoost("title", 2).
		FieldWithBoost("theme", 1).
		FieldWithBoost("director", 1).
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
