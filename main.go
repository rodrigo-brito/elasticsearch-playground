package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/golang/glog"
	"github.com/olivere/elastic"
)

func main() {
	client, err := elastic.NewClient(elastic.SetURL("http://172.17.0.2:9200"))
	if err != nil {
		panic(err)
	}

	http.HandleFunc("/elastic", func(w http.ResponseWriter, r *http.Request) {
		search := elastic.NewMultiMatchQuery(r.URL.Query().Get("q"), "title", "content").
			Operator("AND").       //Should match all terms
			Type("phrase_prefix"). //Find by prefix
			Slop(5)                //Max difference of terms's order

		result, err := client.Search().
			Index("blog").
			Type("article").
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
	})

	fmt.Println("Server started at 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
