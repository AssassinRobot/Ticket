package main

import (
	"context"
	"log"
	"user/config"
	"user/internal/adapters/database"
	"user/internal/adapters/database/postgres"
	"user/internal/adapters/event/nats"
	"user/internal/adapters/rest"
	"user/internal/application/core/api"
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

	databaseAdapter := postgres.NewDatabasePostgresAdapter(db)

	userEventResponder := nats.NewUserEventResponderAdapter(natsConn, databaseAdapter.GetUserByID)

	var responds = func ()  {
		go func() {
			if err := userEventResponder.ReplyToGetUser(ctx); err != nil {
				log.Printf("Error replying to get user: %v", err)
			}
		}()
	}

	responds()

	userEventPublisher, err := nats.NewEventPublisherAdapter(natsConn)
	if err != nil {
		log.Fatalf("Error creating user event publisher: %v", err)
	}

	api := api.NewAPI(
		databaseAdapter,
		userEventPublisher,
	)
	
	rest.Start(api, config.GetServerPort())

}
