package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/rodrigo-brito/elasticsearch-playground/action"
	_ "github.com/rodrigo-brito/elasticsearch-playground/conf"
	handle "github.com/rodrigo-brito/elasticsearch-playground/http"

	"github.com/golang/glog"
)

func main() {
	ctx := context.Background()

	client, err := action.GetConnection()
	if err != nil {
		glog.Fatal(err)
	}

	if err := action.CreateIndex(ctx, client); err != nil {
		glog.Fatal(err)
	}
	http.HandleFunc("/elastic", handle.HandleSearch)

	fmt.Println("Server started at localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
