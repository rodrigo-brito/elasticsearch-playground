package http

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/spf13/viper"

	"github.com/rodrigo-brito/elasticsearch-playground/action"

	"github.com/golang/glog"
	"github.com/olivere/elastic"
)

func HandleSearch(w http.ResponseWriter, r *http.Request) {
	client, err := action.GetConnection()
	if err != nil {
		glog.Error(err)
		w.WriteHeader(http.StatusServiceUnavailable)
		return
	}

	term := r.URL.Query().Get("q")

	termQuery := elastic.NewMultiMatchQuery(term, "title^1.5", "theme", "director").
		Operator("OR").
		Type("most_fields").
		Fuzziness("AUTO").
		CutoffFrequency(0.0001).
		Slop(5)

	exactMatch := elastic.NewMultiMatchQuery(term, "title^1.5", "theme", "director").
		Operator("AND").
		Type("phrase_prefix").
		Slop(5)

	query := elastic.NewDisMaxQuery().
		Query(exactMatch, termQuery).
		TieBreaker(1.1)

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
