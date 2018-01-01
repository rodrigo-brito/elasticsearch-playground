package action

import (
	"context"

	"fmt"
	"github.com/olivere/elastic"
	"github.com/spf13/viper"
)

func CreateIndex(ctx context.Context, client *elastic.Client) error {
	indexName := viper.GetString("indexName")
	mapping := `{
		"settings":{
			"number_of_shards":1,
			"number_of_replicas":0,
			"analysis":{
				"analyzer":{
					"brazilian":{
						"tokenizer":"ngran_tokenizer",
						"filter":[
							"lowercase",
							"asciifolding"
						]
					}
				},
				"tokenizer": {
					"ngran_tokenizer": {
						"type": "ngram",
						"min_gram": 1,
						"max_gram": 3,
						"token_chars": [
							"letter",
							"digit"
						]
					}
				}
			}
		},
		"mappings":{
			"movies":{
				"properties":{
					"id":{
						"type": "integer"
					},
					"title":{
						"type":"text",
						"analyzer": "brazilian"
					},
					"theme":{
						"type":"text",
						"analyzer": "brazilian"
					},
					"director":{
						"type":"text",
						"analyzer": "brazilian"
					},
					"year":{
						"type":"text"
					}
				}
			}
		}
	}`

	if ok, err := client.IndexExists(indexName).Do(ctx); err != nil {
		return err
	} else if ok {
		fmt.Println("Index already exists")
		return nil
	}

	if result, err := client.CreateIndex(indexName).BodyString(mapping).Do(ctx); err != nil {
		return err
	} else if result.Acknowledged {
		fmt.Println("Index created...")
	}

	return InsertFakeData(ctx)
}
