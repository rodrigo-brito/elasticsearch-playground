package elasticsearch

import (
	"fmt"
	"github.com/olivere/elastic"
	"github.com/spf13/viper"
)

func GetConnection() (*elastic.Client, error) {
	cfg := viper.GetStringMapString("elastic")
	return elastic.NewClient(
		elastic.SetURL(fmt.Sprintf("%s:%s", cfg["address"], cfg["port"])),
		elastic.SetBasicAuth(cfg["user"], cfg["password"]),
	)
}
