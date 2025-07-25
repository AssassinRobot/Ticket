package nats

import (
	"context"
	"user/internal/application/core/domain"
	"user/utils"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
)

const (
	StreamName                = "USER"
	UserSubject               = "user.*"
	UserRegisteredSubjectName = "user.registered"
	UserUpdatedSubjectName    = "user.updated"
)

type UserEventPublisherAdapter struct {
	jetStream jetstream.JetStream
}

func NewEventPublisherAdapter(nats *nats.Conn) (*UserEventPublisherAdapter, error) {
	jetStream, err := jetstream.New(nats)
	if err != nil {
		return nil, err
	}

	userEventPublisherAdapter := &UserEventPublisherAdapter{
		jetStream: jetStream,
	}

	return userEventPublisherAdapter, nil
}

func (u *UserEventPublisherAdapter) PublishUserRegistered(ctx context.Context, user *domain.User) error {
	_, err := u.jetStream.CreateOrUpdateStream(ctx, jetstream.StreamConfig{
		Name: StreamName,
		Subjects: []string{UserSubject},
		Storage:  jetstream.FileStorage,
	})

	if err != nil {
		return err
	}

	data, err := utils.Marshal(user)
	if err != nil {
		return err
	}

	_, err = u.jetStream.Publish(ctx, UserRegisteredSubjectName, data)

	return err
}

func (u *UserEventPublisherAdapter) PublishUserUpdated(ctx context.Context,user *domain.User) error {
	_, err := u.jetStream.CreateOrUpdateStream(ctx, jetstream.StreamConfig{
		Name: StreamName,
		Subjects: []string{UserSubject},
		Storage:  jetstream.FileStorage,
	})

	if err != nil {
		return err
	}

	data, err := utils.Marshal(user)
	if err != nil {
		return err
	}	
	_, err = u.jetStream.Publish(ctx, UserUpdatedSubjectName, data)

	return err
}
