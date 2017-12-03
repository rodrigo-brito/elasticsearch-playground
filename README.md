# Elasticsearch playground - Go
Example of queries for Elasticsearch in Go


### My index mapping
```json
{
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
      "article":{
         "properties":{
            "title":{
               "type":"text",
               "analyzer":"brazilian"
            },
            "content":{
               "type":"text",
               "analyzer":"brazilian"
            }
         }
      }
   }
}
```
