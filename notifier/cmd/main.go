package main

import (
	"github.com/nozgurozturk/marvin/notifier/internal/app"
	"github.com/nozgurozturk/marvin/notifier/internal/config"
	"github.com/nozgurozturk/marvin/notifier/internal/service"
	"github.com/nozgurozturk/marvin/notifier/internal/storage"
	"sync"
)

func main() {
	config.Set()
	mongo, err := storage.MongoConnect()
	if err != nil {
		return
	}


	s := storage.New(mongo)

	repoService := service.NewRepoService(s.Repos())
	subscriberService := service.NewSubscriberService(s.Subscribers())
	var wg sync.WaitGroup
	wg.Add(1)
	go app.SendNotification(subscriberService, repoService)
	wg.Wait()
}
