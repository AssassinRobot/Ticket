package main

import (
	"context"
	"log"
	"notification/config"
	"notification/internal/adapters/event/nats"
	"notification/internal/adapters/notification/email"
	"notification/internal/application/core/api"
)

func main() {
	ctx := context.Background()

	err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	natsConn, err := nats.NatsConnection(config.GetNatsURL())
	if err != nil {
		log.Fatalf("Error initializing nats: %v", err)
	}

	emailNotifier := email.NewEmailNotifier(config.GetSMTPConfigs())

	APIadapter := api.NewAPIAdapter(emailNotifier)

	eventConsumer, err := nats.NewNatsEventConsumer(natsConn, APIadapter)
	if err != nil {
		log.Fatalf("Error Jet Stream connection: %v", err)
	}

	errorStream := eventConsumer.ConsumerUserEvents(ctx)
	
	for err := range errorStream {
		log.Printf("Error received: %v\n", err)
	}
}
