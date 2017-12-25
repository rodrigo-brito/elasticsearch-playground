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

	search := elastic.NewMultiMatchQuery(r.URL.Query().Get("q"), "title", "content").
		Operator("AND").       //Should match all terms
		Type("phrase_prefix"). //Find by prefix
		Slop(5)                //Max difference of terms's order

	result, err := client.Search().
		Index(viper.GetString("indexName")).
		Type("movie").
		Query(search).
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
