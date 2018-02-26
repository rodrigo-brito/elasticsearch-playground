package conf

import "github.com/spf13/viper"

func init() {
	viper.SetDefault("indexName", "britoflix")
	viper.SetDefault("elastic", map[string]string{
		"user":     "elastic",
		"password": "changeme",
		"address":  "http://127.0.0.1",
		"port":     "9200",
	})
}
