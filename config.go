package main

import (
	"github.com/caarlos0/env/v6"
	"github.com/go-mysql-org/go-mysql/canal"
)

var Config struct {
	DbUrl            string `env:"DB_URL,required"`
	ElasticsearchUrl string `env:"ELASTICSEARCH_URL" envDefault:"http://localhost:9200"`
	Username         string `env:"USERNAME"`
	Password         string `env:"PASSWORD"`
	DBName           string `env:"DB_NAME" envDefault:"fduhole"`
	Dump             bool   `env:"DUMP" envDefault:"false"`
}

func InitConfig() {
	err := env.Parse(&Config)
	if err != nil {
		panic(err)
	}
}

func NewConfig() *canal.Config {
	cfg := canal.NewDefaultConfig()
	cfg.Addr = Config.DbUrl
	cfg.User = Config.Username
	cfg.Password = Config.Password
	cfg.Dump.TableDB = Config.DBName
	cfg.Dump.Tables = []string{"floor"}

	return cfg
}
