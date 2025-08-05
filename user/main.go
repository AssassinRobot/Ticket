package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"user/config"
	"user/internal/adapters/database"
	"user/internal/adapters/database/postgres"
	"user/internal/adapters/event/nats"
	"user/internal/application/core/api"
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

	databaseAdapter := postgres.NewDatabasePostgresAdapter(db)

	userEventPublisher, err := nats.NewEventPublisherAdapter(natsConn)
	if err != nil {
		log.Fatalf("Error creating user event publisher: %v", err)
	}

	api := api.NewAPI(
		databaseAdapter,
		userEventPublisher,
	)

	userEventResponder := nats.NewUserEventResponderAdapter(natsConn, api)

	go func() {
		err := userEventResponder.ReplyToListUsers(ctx)

		if err != nil {
			log.Printf("Error replying to get user: %v", err)
		}
	}()

	go func() {
		err := userEventResponder.ReplyToGetUserByID(ctx)

		if err != nil {
			log.Printf("Error replying to get user: %v", err)
		}
	}()

	go func() {
		err := userEventResponder.ReplayToCreateUser(ctx)

		if err != nil {
			log.Printf("Error replying to get user: %v", err)
		}
	}()

	go func() {
		err := userEventResponder.ReplayToUpdateUserByID(ctx)

		if err != nil {
			log.Printf("Error replying to get user: %v", err)
		}
	}()

	go func() {
		err := userEventResponder.ReplayToDeleteUserByID(ctx)
		if err != nil {
			log.Printf("Error replying to get user: %v", err)
		}
	}()

	<-ctx.Done()
}
