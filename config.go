package main

import (
	"github.com/caarlos0/env/v10"
)

var Config struct {
	DbUrl            string `env:"DB_URL,required"`
	ElasticsearchUrl string `env:"ELASTICSEARCH_URL" envDefault:"http://localhost:9200"`
}

func InitConfig() {
	err := env.Parse(&Config)
	if err != nil {
		panic(err)
	}
}
