package action

import (
	"github.com/olivere/elastic"
)

func GetConnection() (*elastic.Client, error) {
	return elastic.NewClient(
		elastic.SetURL("http://172.17.0.2:9200"),
		elastic.SetBasicAuth("elastic", "changeme"),
	)
}
