package main

import (
	_ "github.com/nozgurozturk/marvin/server/docs"
	"github.com/nozgurozturk/marvin/server/internal/config"
	"github.com/nozgurozturk/marvin/server/internal/router"
	"github.com/nozgurozturk/marvin/server/internal/storage"
	"log"
)

// @title Marvin
// @version 0.0.1
// @description API for Marvin outdated package dependency notification service

// @contact.name Ozgur Ozturk
// @contact.email ozgur@nozgurozturk.com

// @host localhost:8080

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization

func main() {

	cnf := config.Set()
	mongo, err := storage.MongoConnect()
	if err != nil {
		return
	}
	redis, err := storage.RedisConnect()
	if err != nil {
		return
	}

	s := storage.New(mongo, redis)
	r := router.New(s)
	// go app.Check(r.Service.Subscriber())
	err = r.Router.Listen(cnf.HTTP.Port)
	if err != nil {
		log.Fatal(err)
	}
}
