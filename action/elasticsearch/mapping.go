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
					"brazilian_ngram":{
						"tokenizer":"ngram_tokenizer",
						"filter":[
							"lowercase",
							"asciifolding",
							"word_joiner"
						]
					},
					"brazilian_shingle":{
						"tokenizer":"standard",
						"filter":[
							"lowercase",
							"asciifolding",
							"filter_shingle"
						]
					},
					"search_analyzer":{
						"tokenizer":"standard",
						"filter":[
							"lowercase",
							"asciifolding",
							"word_joiner"
						]
					}
				},
				"tokenizer": {
					"ngram_tokenizer": {
						"type": "edge_ngram",
						"min_gram": 1,
						"max_gram": 20,
						"token_chars": [
							"letter",
							"digit"
						]
					}
				},
				"filter":{
					"filter_shingle":{
						"type":"shingle",
						"max_shingle_size": 2,
						"min_shingle_size": 2,
						"output_unigrams": "false",
						"token_separator": ""
					},
					"word_joiner": {
						"type": "word_delimiter",
						"catenate_all": true
					}
				}
			}
		},
		"mappings":{
			"artists": {
				"properties":{
					"id":{
						"type": "integer"
					},
					"name":{
						"type": "text",
						"analyzer": "brazilian"
					},
					"name_shingle":{
						"type": "text",
						"analyzer": "brazilian_shingle"
					},
					"genre":{
						"type":"text",
						"analyzer": "brazilian"
					},
					"plays":{
						"type": "integer" 
					}
				}
			},
			"movies":{
				"properties":{
					"id":{
						"type": "integer"
					},
					"title":{
						"type":"text",
						"analyzer": "brazilian"
					},
					"title_ngram":{
						"type":"text",
						"analyzer": "brazilian_ngram",
						"search_analyzer": "brazilian_ngram"
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
					"director_ngram":{
						"type":"text",
						"analyzer": "brazilian_ngram",
						"search_analyzer": "brazilian_ngram"
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
					"title_director_ngram":{
						"type":"text",
						"analyzer": "brazilian_ngram",
						"search_analyzer": "brazilian_ngram"
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
