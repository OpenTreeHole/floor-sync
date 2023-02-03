package main

import (
	"github.com/caarlos0/env/v6"
	"net/url"
)

var Config struct {
	DbUrl            url.URL `env:"DB_URL,required"`
	ElasticsearchUrl url.URL `env:"ELASTICSEARCH_URL" envDefault:"http://localhost:9200"`
	Dump             bool    `env:"DUMP" envDefault:"false"`
}

func InitConfig() {
	err := env.Parse(&Config)
	if err != nil {
		panic(err)
	}
}
