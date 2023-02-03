package main

import (
	"github.com/go-mysql-org/go-mysql/canal"
	"log"
)

func main() {
	InitConfig()
	InitSearch()
	InitDB()

	if Config.Dump {
		Dump()
	}

	c, err := canal.NewCanal(NewConfig())
	if err != nil {
		log.Fatal(err)
	}

	c.SetEventHandler(&MyEventHandler{})

	position, err := c.GetMasterPos()
	if err != nil {
		panic(err)
	}

	err = c.RunFrom(position)
	if err != nil {
		panic(err)
	}
}
