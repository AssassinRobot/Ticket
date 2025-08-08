package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"train/config"
	"train/internal/adapter/database"
	"train/internal/adapter/database/postgres"
	"train/internal/adapter/event/nats"

	"train/internal/application/core/api"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Handle interrupt signal
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigChan
		log.Println("Shutting down gracefully...")
		cancel()
	}()

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

	api := api.NewAPIAdapter(databaseAdapter)

	trainEventResponder := nats.NewTrainEventResponderAdapter(natsConn, api)

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

	go func() {
		err := trainEventResponder.ReplyToGetTrainByID(ctx)
		if err != nil {
			log.Printf("Error replying to get train by ID: %v", err)
		}
	}()

	go func() {
		err := trainEventResponder.ReplyToCreateTrain(ctx)
		if err != nil {
			log.Printf("Error replying to create train: %v", err)
		}
	}()

	go func() {
		err := trainEventResponder.ReplyToUpdateTrainName(ctx)
		if err != nil {
			log.Printf("Error replying to update train: %v", err)
		}
	}()

	go func() {
		err := trainEventResponder.ReplyToUpdateTrainTravelDetails(ctx)
		if err != nil {
			log.Printf("Error replying to update train: %v", err)
		}
	}()

	go func() {
		err := trainEventResponder.ReplyToDeleteTrainByID(ctx)
		if err != nil {
			log.Printf("Error replying to delete train: %v", err)
		}
	}()

	go func() {
		err := trainEventResponder.ReplyToListSeatsByTrainID(ctx)
		if err != nil {
			log.Printf("Error replying to list seats by train ID: %v", err)
		}
	}()

	go func() {
		err := trainEventResponder.ReplyToGetSeatByID(ctx)
		if err != nil {
			log.Printf("Error replying to get seat by ID: %v", err)
		}
	}()

	go func() {
		err := trainEventResponder.ReplyToCreateSeat(ctx)
		if err != nil {
			log.Printf("Error replying to create seat: %v", err)
		}
	}()

	go func() {
		err := trainEventResponder.ReplyToUpdateSeatNumberBySeatID(ctx)
		if err != nil {
			log.Printf("Error replying to update seat: %v", err)
		}
	}()

	go func() {
		err := trainEventResponder.ReplyToDeleteSeatBySeatID(ctx)
		if err != nil {
			log.Printf("Error replying to delete seat: %v", err)
		}
	}()

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

	log.Println("Train service is starting....")

	<-ctx.Done()
}
