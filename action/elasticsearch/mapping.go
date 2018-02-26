package elasticsearch

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
						"tokenizer":"standard",
						"filter":[
							"lowercase",
							"asciifolding"
						]
					},
					"brazilian_ngran":{
						"tokenizer":"ngran_tokenizer",
						"filter":[
							"lowercase",
							"asciifolding"
						]
					},
					"brazilian_shingle":{
						"tokenizer":"standard",
						"filter":[
							"lowercase",
							"asciifolding",
							"filter_shingle"
						]
					}
				},
				"tokenizer": {
					"ngran_tokenizer": {
						"type": "ngram",
						"min_gram": 3,
						"max_gram": 10,
						"token_chars": [
							"letter",
							"digit"
						]
					}
				},
				"filter":{
				   "filter_shingle":{
						"type":"shingle",
						"max_shingle_size": 3,
						"min_shingle_size": 2,
						"output_unigrams": "true",
						"token_separator": ""
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
					"title_ngran":{
						"type":"text",
						"analyzer": "brazilian_ngran",
						"search_analyzer": "brazilian_ngran"
					},
					"title_shingle":{
						"type":"text",
						"analyzer": "brazilian_shingle",
						"search_analyzer": "brazilian_shingle"
					},
					"theme":{
						"type":"text",
						"analyzer": "brazilian"
					},
					"director":{
						"type":"text",
						"analyzer": "brazilian"
					},
					"director_ngran":{
						"type":"text",
						"analyzer": "brazilian_ngran",
						"search_analyzer": "brazilian_ngran"
					},
					"director_shingle":{
						"type":"text",
						"analyzer": "brazilian_shingle",
						"search_analyzer": "brazilian_shingle"
					},
					"title_director":{
						"type":"text",
						"analyzer": "brazilian"
					},
					"title_director_ngran":{
						"type":"text",
						"analyzer": "brazilian_ngran",
						"search_analyzer": "brazilian_ngran"
					},
					"views":{
						"type":"integer"
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