package main

import (
	"context"
	"net/http"
	"os"

	"github.com/go-chi/chi"
	"github.com/golang/glog"
	"github.com/olivere/elastic"

	"github.com/rodrigo-brito/elasticsearch-playground/action/elasticsearch"
	_ "github.com/rodrigo-brito/elasticsearch-playground/conf"
	"github.com/rodrigo-brito/elasticsearch-playground/handler/artists"
	"github.com/rodrigo-brito/elasticsearch-playground/handler/movies"
)

type Project struct {
	ctx           context.Context
	osSignal      chan os.Signal
	elasticClient *elastic.Client
}

func (p *Project) Init(ctx context.Context) {
	client, err := elasticsearch.GetConnection()
	if err != nil {
		glog.Fatal(err)
	}

	if err := elasticsearch.CreateIndex(ctx, client); err != nil {
		glog.Fatal(err)
	}

	r := chi.NewRouter()

	r.Route("/artist", func(r chi.Router) {
		r.Get("/v1", artists.SearchHandler)
	})

	r.Route("/movies", func(r chi.Router) {
		r.Get("/v1", movies.QueryStringWithSlplit)
		r.Get("/v2", movies.MultiMatchNgran)
		r.Get("/v3", movies.MultiMatchPrefix)
		r.Get("/v4", movies.MultiMatchPrefixShingle)
		r.Get("/v5", movies.PrefixPhraseNgran)
	})

	http.ListenAndServe(":3000", r)
}

func main() {
	project := new(Project)
	project.Init(context.Background())
}
