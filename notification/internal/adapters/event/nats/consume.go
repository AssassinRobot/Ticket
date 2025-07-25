package nats

import (
	"context"
	"errors"
	"fmt"
	"notification/internal/ports"
	"notification/utils"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
)

const (
	UserStream                = "USER"
	UserRegisteredSubjectName = "user.registered"
	UserUpdatedSubjectName    = "user.updated"
)

var (
	ErrInvalidSubject = errors.New("invalid  user subject")
)

type NatsEventConsumer struct {
	jetStream jetstream.JetStream
	API       ports.APIPort
}

func NewNatsEventConsumer(nats *nats.Conn, API ports.APIPort) (*NatsEventConsumer, error) {
	jetStream, err := jetstream.New(nats)
	if err != nil {
		return nil, err
	}

	natsEventConsumerInstance := &NatsEventConsumer{
		jetStream: jetStream,
		API:       API,
	}

	return natsEventConsumerInstance, nil
}

func (c *NatsEventConsumer) ConsumerUserEvents(ctx context.Context) <-chan error {
	userConsumer, err := c.jetStream.CreateOrUpdateConsumer(ctx, UserStream, jetstream.ConsumerConfig{
		Durable:       "user_email_consumer",
		MaxAckPending: 5,
	})

	var errorStream = make(chan error)

	if err != nil {
		errorStream <- err

		close(errorStream)

		return errorStream
	}

	go func() {
		defer close(errorStream)

		for {
			msgs, err := userConsumer.Fetch(5)
			if err != nil {
				errorStream <- err
			}

			for msg := range msgs.Messages() {
				switch msg.Subject() {
				case UserRegisteredSubjectName:
					name, email, err := utils.UnMarshalUserData(msg.Data())
					if err != nil {
						errorStream <- err
						ack(msg, errorStream)
						continue
					}

					err = c.API.UserRegistered(ctx, name, email)
					if err != nil {
						errorStream <- err
						ack(msg, errorStream)
						continue
					}

					ack(msg, errorStream)

				case UserUpdatedSubjectName:
					name, email, err := utils.UnMarshalUserData(msg.Data())
					if err != nil {
						errorStream <- err
						ack(msg, errorStream)
						continue
					}

					err = c.API.UserUpdated(ctx, name, email)
					if err != nil {
						errorStream <- err
						ack(msg, errorStream)
						continue
					}

					ack(msg, errorStream)
				default:
					errorStream <- fmt.Errorf("%w:%s", ErrInvalidSubject, msg.Subject())
				}

			}

		}
	}()

	return errorStream
}

func ack(msg jetstream.Msg, errorStream chan<- error) {
	err := msg.Ack()
	if err != nil {
		errorStream <- err
	}
}
