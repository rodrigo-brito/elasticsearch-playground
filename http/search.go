package http

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/spf13/viper"

	"github.com/rodrigo-brito/elasticsearch-playground/action"

	"fmt"
	"github.com/golang/glog"
	"github.com/olivere/elastic"
	"strings"
)

func getFuzzyTerm(term string) string {
	return fmt.Sprintf("%s~AUTO", strings.Join(strings.Split(term, " "), "~AUTO "))
}

func HandleSearch(w http.ResponseWriter, r *http.Request) {
	client, err := action.GetConnection()
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
		Do(context.Background())

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
