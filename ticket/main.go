package main

import (
	"log"
	"ticket/config"
	"ticket/internal/adapters/database"
	"ticket/internal/adapters/database/postgres"
	"ticket/internal/adapters/event/nats"
	"ticket/internal/adapters/rest"
	"ticket/internal/application/core/api"
)

func main() {
	err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	dbConnection, err := database.NewDB(config.GetDatabaseURL())
	if err != nil {
		log.Fatal(err)
	}

	databaseAdapter := postgres.NewPostgresDBAdapter(dbConnection)

	natsConn, err := nats.NatsConnection(config.GetNatsURL())
	if err != nil {
		log.Fatal(err)
	}

	eventPublisherAdapter, err := nats.NewEventPublisherAdapter(natsConn)
	if err != nil {
		log.Fatal(err)
	}

	requestSenderAdapter := nats.NewTicketRequestSenderAdapter(natsConn)

	apiAdapter := api.NewAPIAdapter(databaseAdapter,eventPublisherAdapter,requestSenderAdapter)

	rest.Start(apiAdapter,config.GetServerPort())
}
