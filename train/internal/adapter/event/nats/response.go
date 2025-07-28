package nats

import (
	"context"
	"fmt"
	"strconv"
	"train/internal/application/core/domain"
	"train/utils"

	"github.com/nats-io/nats.go"
)

const (
	SubjectListTrain         = "request.trains.list"
	SubjectGetTrain          = "request.trains.get"
	SubjectListTrainFiltered = "request.trains.list.filter"
)

type TrainEventResponderAdapter struct {
	natsConn *nats.Conn
	TrainEventHandlers
}

type TrainEventHandlers interface {
	GetTrainByID(ctx context.Context, trainID uint) (*domain.Train, error)
	ListTrains(ctx context.Context) ([]domain.Train, error)
	ListTrainsFiltered(ctx context.Context, filter *domain.TrainFilters) ([]domain.Train, error)
}

func NewTrainEventResponderAdapter(natsConn *nats.Conn, trainHandlers TrainEventHandlers) *TrainEventResponderAdapter {
	return &TrainEventResponderAdapter{
		natsConn:           natsConn,
		TrainEventHandlers: trainHandlers,
	}
}

func (u *TrainEventResponderAdapter) ReplyToGetTrainByID(ctx context.Context) error {
	subscription, err := u.natsConn.Subscribe(SubjectGetTrain, func(msg *nats.Msg) {
		trainIDstr := string(msg.Data)
		trainID, err := strconv.Atoi(trainIDstr)

		if err != nil {
			errString := fmt.Errorf("error invalid trainID: %d", trainID).Error()

			msg.Respond([]byte(errString))
			return
		}

		train, err := u.TrainEventHandlers.GetTrainByID(ctx, uint(trainID))
		if err != nil {
			errString := fmt.Errorf("error get train by ID  : %d,%w", trainID, err).Error()
			msg.Respond([]byte(errString))
			return
		}

		serializedTrainData, err := utils.MarshalTrain(train)
		if err != nil {
			fmt.Println()
			errString := fmt.Errorf("error serialize train: %v,%w", train, err).Error()
			msg.Respond([]byte(errString))
			return
		}

		msg.Respond(serializedTrainData)
	})

	if err != nil {
		return err
	}

	go func() {
		<-ctx.Done()
		subscription.Unsubscribe()
	}()

	return nil
}

func (u *TrainEventResponderAdapter) ReplyToListTrains(ctx context.Context) error {
	subscription, err := u.natsConn.Subscribe(SubjectListTrain, func(msg *nats.Msg) {
		trains, err := u.TrainEventHandlers.ListTrains(ctx)
		if err != nil {
			errString := fmt.Errorf("error list trains: %w", err).Error()
			msg.Respond([]byte(errString))
			return
		}

		serializedTrainsData, err := utils.MarshalTrains(trains)
		if err != nil {
			fmt.Println()
			errString := fmt.Errorf("error serialize train: %v,%w", trains, err).Error()
			msg.Respond([]byte(errString))
			return
		}

		msg.Respond(serializedTrainsData)
	})

	if err != nil {
		return err
	}

	go func() {
		<-ctx.Done()
		subscription.Unsubscribe()
	}()

	return nil
}

func (u *TrainEventResponderAdapter) ReplyToListTrainsFiltered(ctx context.Context) error {
	subscription, err := u.natsConn.Subscribe(SubjectListTrainFiltered, func(msg *nats.Msg) {
		filters, err := utils.UnMarshalFilters(msg.Data)
		if err != nil {
			errString := fmt.Errorf("error unmarshal train filters: %w", err).Error()
			msg.Respond([]byte(errString))
			return
		}

		trains, err := u.TrainEventHandlers.ListTrainsFiltered(ctx,filters)
		if err != nil {
			errString := fmt.Errorf("error list trains: %w", err).Error()
			msg.Respond([]byte(errString))
			return
		}

		serializedTrainsData, err := utils.MarshalTrains(trains)
		if err != nil {
			fmt.Println()
			errString := fmt.Errorf("error serialize train: %v,%w", trains, err).Error()
			msg.Respond([]byte(errString))
			return
		}

		msg.Respond(serializedTrainsData)
	})

	if err != nil {
		return err
	}

	go func() {
		<-ctx.Done()
		subscription.Unsubscribe()
	}()

	return nil
}
