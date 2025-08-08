package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"ticket/config"
	"ticket/internal/adapters/database"
	"ticket/internal/adapters/database/postgres"
	"ticket/internal/adapters/event/nats"
	"ticket/internal/application/core/api"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigChan
		log.Println("Shutting down gracefully...")
		cancel()
	}()

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

	apiAdapter := api.NewAPIAdapter(databaseAdapter, eventPublisherAdapter, requestSenderAdapter)

	eventResponderAdapter := nats.NewTicketEventResponderAdapter(natsConn, apiAdapter)

	go func() {
		err := eventResponderAdapter.ReplyToGetTicketByID(ctx)
		if err != nil {
			log.Fatalf("error ReplyToGetTicketByID:%v", err)
		}
	}()

	go func() {
		err := eventResponderAdapter.ReplayToListTicketsByTrainID(ctx)
		if err != nil {
			log.Fatalf("error ReplyToListTicketsByTrainID:%v", err)
		}
	}()

	go func() {
		err := eventResponderAdapter.ReplayToListTicketsByUserID(ctx)
		if err != nil {
			log.Fatalf("error ReplyToListTicketsByUserID:%v", err)
		}
	}()

	go func() {
		err := eventResponderAdapter.ReplayToBookTicket(ctx)
		if err != nil {
			log.Fatalf("error ReplyToBookTicket:%v", err)
		}
	}()

	go func() {
		err := eventResponderAdapter.ReplayToCancelTicket(ctx)
		if err != nil {
			log.Fatalf("error ReplyToCancelTicket:%v", err)
		}
	}()

	log.Println("Ticket service is starting...")
	
	<- ctx.Done()
}
