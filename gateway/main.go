package main

import (
	"gateway/config"
	"gateway/events/nats"
	"gateway/server"
	"log"
)

func main() {
	err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	natsConn, err := nats.NatsConnection(config.GetNatsURL())
	if err != nil {
		log.Fatalf("Error initializing event bus: %v", err)
	}

	requestHandler := nats.NewNatsRequestSender(natsConn)

	err = server.Start(requestHandler,config.GetPort())
	if err != nil {
		log.Fatal(err)
	}
}
