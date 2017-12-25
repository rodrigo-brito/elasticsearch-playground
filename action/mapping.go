package action

import (
	"context"

	"github.com/olivere/elastic"
)

const indexName = "britoflix"

func CreateIndex(ctx context.Context, client *elastic.Client) error {
	mapping := `{
		"settings":{
			"number_of_shards":1,
			"number_of_replicas":0,
			"analysis":{
				"analyzer":{
					"brazilian":{
						"tokenizer":"standard",
						"filter":[
							"lowercase",
							"asciifolding"
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
					"category":{
						"type":"text",
						"analyzer": "brazilian" 
					}
				}
			}
		}
	}`

	ok, err := client.IndexExists(indexName).Do(ctx)
	if err != nil {
		return err
	}
	if ok {
		return nil

	}
	_, err = client.CreateIndex(indexName).BodyString(mapping).Do(ctx)
	return err
}
