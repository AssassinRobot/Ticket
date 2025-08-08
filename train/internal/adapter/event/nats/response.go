package nats

import (
	"context"
	"fmt"
	"strconv"
	"train/internal/ports"
	"train/utils"

	"github.com/nats-io/nats.go"
)

const (
	SubjectRequestListTrains              = "request.list.trains"
	SubjectRequestListTrainsFiltered      = "request.list.trains.filtered"
	SubjectRequestGetTrainByID            = "request.get.train.byID"
	SubjectRequestCreateTrain             = "request.create.train"
	SubjectRequestUpdateTrain             = "request.update.train"
	SubjectRequestUpdateTrainTravelDetail = "request.update.train.travel.detail"
	SubjectRequestDeleteTrainByID         = "request.delete.train.byID"

	SubjectRequestGetSeatByID        = "request.get.seat.byID"
	SubjectRequestListSeatsByTrainID = "request.list.seats.byTrainID"
	SubjectRequestCreateSeat         = "request.create.seat"
	SubjectRequestUpdateSeatNumber   = "request.update.seat.number"
	SubjectRequestDeleteSeatByID     = "request.delete.seat.byID"
)

type TrainEventResponderAdapter struct {
	natsConn   *nats.Conn
	APIAdapter ports.APIPort
}

func NewTrainEventResponderAdapter(natsConn *nats.Conn, apiAdapter ports.APIPort) *TrainEventResponderAdapter {
	return &TrainEventResponderAdapter{
		natsConn:   natsConn,
		APIAdapter: apiAdapter,
	}
}
func (r *TrainEventResponderAdapter) ReplyToListTrains(ctx context.Context) error {
	subscription, err := r.natsConn.Subscribe(SubjectRequestListTrains, func(msg *nats.Msg) {
		trains, err := r.APIAdapter.ListTrains(ctx)
		if err != nil {
			errString := fmt.Errorf("error list trains: %w", err).Error()
			msg.Respond([]byte(errString))
			return
		}

		serializedTrainsData, err := utils.MarshalListTrainsReplay(trains)
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

func (r *TrainEventResponderAdapter) ReplyToListTrainsFiltered(ctx context.Context) error {
	subscription, err := r.natsConn.Subscribe(SubjectRequestListTrainsFiltered, func(msg *nats.Msg) {
		filters, err := utils.UnmarshalFilters(msg.Data)
		if err != nil {
			errString := fmt.Errorf("error unmarshal train filters: %w", err).Error()
			msg.Respond([]byte(errString))
			return
		}

		trains, err := r.APIAdapter.ListTrainsFiltered(ctx, filters)
		if err != nil {
			errString := fmt.Errorf("error list trains: %w", err).Error()
			msg.Respond([]byte(errString))
			return
		}

		serializedTrainsData, err := utils.MarshalListTrainsReplay(trains)
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

func (r *TrainEventResponderAdapter) ReplyToGetTrainByID(ctx context.Context) error {
	subscription, err := r.natsConn.Subscribe(SubjectRequestGetTrainByID, func(msg *nats.Msg) {
		trainIDstr := string(msg.Data)
		trainID, err := strconv.Atoi(trainIDstr)

		if err != nil {
			errString := fmt.Errorf("error invalid trainID: %d", trainID).Error()

			msg.Respond([]byte(errString))
			return
		}

		train, err := r.APIAdapter.GetTrainByID(ctx, uint(trainID))
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

func (r *TrainEventResponderAdapter) ReplyToCreateTrain(ctx context.Context) error {
	subscription, err := r.natsConn.Subscribe(SubjectRequestCreateTrain, func(msg *nats.Msg) {
		trainName, capacity, err := utils.UnmarshalCreateTrainRequest(msg.Data)
		if err != nil {
			errString := fmt.Errorf("error unmarshal request(CreateTrain): %w", err).Error()
			msg.Respond([]byte(errString))
			return
		}

		if err := r.APIAdapter.CreateTrain(ctx, trainName, capacity); err != nil {
			errString := fmt.Errorf("error create train: %w", err).Error()
			msg.Respond([]byte(errString))
			return
		}

		msg.Respond([]byte("train created successfully"))
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

func (r *TrainEventResponderAdapter) ReplyToUpdateTrainName(ctx context.Context) error {
	subscription, err := r.natsConn.Subscribe(SubjectRequestUpdateTrain, func(msg *nats.Msg) {
		trainID, name, err := utils.UnmarshalUpdateTrainNameRequest(msg.Data)
		if err != nil {
			errString := fmt.Errorf("error unmarshal request(UpdateTrainName): %w", err).Error()
			msg.Respond([]byte(errString))
			return
		}

		if err := r.APIAdapter.UpdateTrain(ctx, uint(trainID), name); err != nil {
			errString := fmt.Errorf("error update train name: %w", err).Error()
			msg.Respond([]byte(errString))
			return
		}

		msg.Respond([]byte("train updated successfully"))
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

func (r *TrainEventResponderAdapter) ReplyToUpdateTrainTravelDetails(ctx context.Context) error {
	subscription, err := r.natsConn.Subscribe(SubjectRequestUpdateTrainTravelDetail, func(msg *nats.Msg) {
		trainID, travelDetails, err := utils.UnmarshalUpdateTrainTravelDetailsRequest(msg.Data)
		if err != nil {
			errString := fmt.Errorf("error unmarshal request(UpdateTrainTravelDetails): %w", err).Error()
			msg.Respond([]byte(errString))
			return
		}

		err = r.APIAdapter.UpdateTrainTravelDetails(ctx, uint(trainID), travelDetails)
		if err != nil {
			errString := fmt.Errorf("error update train travel details: %w", err).Error()
			msg.Respond([]byte(errString))
			return
		}

		msg.Respond([]byte("train travel details updated successfully"))
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

func (r *TrainEventResponderAdapter) ReplyToDeleteTrainByID(ctx context.Context) error {
	subscription, err := r.natsConn.Subscribe(SubjectRequestDeleteTrainByID, func(msg *nats.Msg) {
		trainIDstr := string(msg.Data)
		trainID, err := strconv.Atoi(trainIDstr)

		if err != nil {
			errString := fmt.Errorf("error invalid trainID: %d", trainID).Error()
			msg.Respond([]byte(errString))
			return
		}

		if err := r.APIAdapter.DeleteTrain(ctx, uint(trainID)); err != nil {
			errString := fmt.Errorf("error delete train by ID  : %d,%w", trainID, err).Error()
			msg.Respond([]byte(errString))
			return
		}

		msg.Respond([]byte("train deleted successfully"))
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

func (r *TrainEventResponderAdapter) ReplyToListSeatsByTrainID(ctx context.Context) error {
	subscription, err := r.natsConn.Subscribe(SubjectRequestListSeatsByTrainID, func(msg *nats.Msg) {
		trainIDstr := string(msg.Data)
		trainID, err := strconv.Atoi(trainIDstr)

		if err != nil {
			errString := fmt.Errorf("error invalid trainID: %d", trainID).Error()
			msg.Respond([]byte(errString))
			return
		}

		seats, err := r.APIAdapter.ListSeatsByTrainID(ctx, uint(trainID))
		if err != nil {
			errString := fmt.Errorf("error list seats by train ID  : %d,%w", trainID, err).Error()
			msg.Respond([]byte(errString))
			return
		}

		serializedSeatsData, err := utils.MarshalSeats(seats)
		if err != nil {
			fmt.Println()
			errString := fmt.Errorf("error serialize seats: %v,%w", seats, err).Error()
			msg.Respond([]byte(errString))
			return
		}

		msg.Respond(serializedSeatsData)
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

func (r *TrainEventResponderAdapter) ReplyToGetSeatByID(ctx context.Context) error {
	subscription, err := r.natsConn.Subscribe(SubjectRequestGetSeatByID, func(msg *nats.Msg) {
		seatIDstr := string(msg.Data)
		seatID, err := strconv.Atoi(seatIDstr)

		if err != nil {
			errString := fmt.Errorf("error invalid seatID: %d", seatID).Error()
			msg.Respond([]byte(errString))
			return
		}

		seat, err := r.APIAdapter.GetSeatByID(ctx, uint(seatID))
		if err != nil {
			errString := fmt.Errorf("error get seat by ID  : %d,%w", seatID, err).Error()
			msg.Respond([]byte(errString))
			return
		}

		serializedSeatData, err := utils.MarshalSeat(seat)
		if err != nil {
			fmt.Println()
			errString := fmt.Errorf("error serialize seat: %v,%w", seat, err).Error()
			msg.Respond([]byte(errString))
			return
		}

		msg.Respond(serializedSeatData)
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

func (r *TrainEventResponderAdapter) ReplyToCreateSeat(ctx context.Context) error {
	subscription, err := r.natsConn.Subscribe(SubjectRequestCreateSeat, func(msg *nats.Msg) {
		trainID, seatNumber, err := utils.UnmarshalCreateSeatRequest(msg.Data)
		if err != nil {
			errString := fmt.Errorf("error unmarshal request(CreateSeat): %w", err).Error()
			msg.Respond([]byte(errString))
			return
		}

		if err := r.APIAdapter.CreateSeat(ctx, uint(trainID), uint(seatNumber)); err != nil {
			errString := fmt.Errorf("error create seat: %w", err).Error()
			msg.Respond([]byte(errString))
			return
		}

		msg.Respond([]byte("seat created successfully"))
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

func (r *TrainEventResponderAdapter) ReplyToUpdateSeatNumberBySeatID(ctx context.Context) error {
	subscription, err := r.natsConn.Subscribe(SubjectRequestUpdateSeatNumber, func(msg *nats.Msg) {
		seatID, seatNumber, err := utils.UnmarshalUpdateSeatNumberRequest(msg.Data)
		if err != nil {
			errString := fmt.Errorf("error unmarshal request(UpdateSeatNumberRequest): %w", err).Error()
			msg.Respond([]byte(errString))
			return
		}

		if err := r.APIAdapter.UpdateSeatNumber(ctx, uint(seatID), uint(seatNumber)); err != nil {
			errString := fmt.Errorf("error update seat number: %w", err).Error()
			msg.Respond([]byte(errString))
			return
		}

		msg.Respond([]byte("seat number updated successfully"))
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

func (r *TrainEventResponderAdapter) ReplyToDeleteSeatBySeatID(ctx context.Context) error {
	subscription, err := r.natsConn.Subscribe(SubjectRequestDeleteSeatByID, func(msg *nats.Msg) {
		seatIDstr := string(msg.Data)
		seatID, err := strconv.Atoi(seatIDstr)

		if err != nil {
			errString := fmt.Errorf("error invalid seatID: %d", seatID).Error()
			msg.Respond([]byte(errString))
			return
		}

		if err := r.APIAdapter.DeleteSeat(ctx, uint(seatID)); err != nil {
			errString := fmt.Errorf("error delete seat by ID  : %d,%w", seatID, err).Error()
			msg.Respond([]byte(errString))
			return
		}

		msg.Respond([]byte("seat deleted successfully"))
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
