package conf

import "github.com/spf13/viper"

func init() {
	viper.SetDefault("indexName", "britoflix")
	viper.SetDefault("auth", map[string]string{
		"user":     "elastic",
		"password": "changeme",
	})
}
