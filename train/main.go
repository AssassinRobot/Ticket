package main

import (
	"context"
	"log"
	"train/config"
	"train/internal/adapter/database"
	"train/internal/adapter/database/postgres"
	"train/internal/adapter/event/nats"
	"train/internal/adapter/rest"
	"train/internal/application/core/api"
)

func main() {
	ctx := context.Background()

	err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	db, err := database.NewDB(config.GetDatabaseURL())
	if err != nil {
		log.Fatalf("Error initializing database: %v", err)
	}

	natsConn, err := nats.NatsConnection(config.GetNatsURL())
	if err != nil {
		log.Fatalf("Error initializing event bus: %v", err)
	}

	databaseAdapter := postgres.NewPostgresDBAdapter(db)

	trainEventResponder := nats.NewTrainEventResponderAdapter(natsConn, databaseAdapter)

	go func() {
		err := trainEventResponder.ReplyToGetTrainByID(ctx)
		if err != nil {
			log.Printf("Error replying to get train: %v", err)
		}
	}()

	go func() {
		err := trainEventResponder.ReplyToListTrains(ctx)
		if err != nil {
			log.Printf("Error replying to list trains: %v", err)
		}
	}()

	go func() {
		err := trainEventResponder.ReplyToListTrainsFiltered(ctx)
		if err != nil {
			log.Printf("Error replying to list trains filtered: %v", err)
		}
	}()

	api := api.NewAPIAdapter(databaseAdapter)

	trainEventConsumer, err := nats.NewNatsEventConsumer(natsConn, api)
	if err != nil {
		log.Fatalf("Error creating user event publisher: %v", err)
	}

	go func() {
		errStream := trainEventConsumer.ConsumerSeatEvents(ctx)
		for err := range errStream {
			log.Println(err)
		}
	}()

	rest.Start(api, config.GetServerPort())
}
