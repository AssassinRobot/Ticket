package nats

import (
	"github.com/nats-io/nats.go"
)

func NatsConnection(url string) (*nats.Conn, error) {	
	nc, err := nats.Connect(url)
	if err != nil {
		return nil, err
	}

	return nc, nil
}
