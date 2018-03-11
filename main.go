package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/golang/glog"

	"os"
	"os/signal"
	"syscall"

	"github.com/olivere/elastic"
	"github.com/rodrigo-brito/elasticsearch-playground/action/elasticsearch"
	_ "github.com/rodrigo-brito/elasticsearch-playground/conf"
	handle "github.com/rodrigo-brito/elasticsearch-playground/http"
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
	http.HandleFunc("/v1", handle.QueryStringWithSlplit)
	http.HandleFunc("/v2", handle.MultiMatchNgran)
	http.HandleFunc("/v3", handle.MultiMatchPrefix)
	http.HandleFunc("/v4", handle.MultiMatchPrefixShingle)
	http.HandleFunc("/v5", handle.PrefixPhraseNgran)

	p.osSignal = make(chan os.Signal, 2)
	signal.Notify(p.osSignal, os.Interrupt, syscall.SIGTERM)
	go p.close()

	fmt.Println("Server started at localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func (p *Project) close() {
	<-p.osSignal
	fmt.Println("Killing gracefully... :)")
}

func main() {
	project := new(Project)
	project.Init(context.Background())
}
