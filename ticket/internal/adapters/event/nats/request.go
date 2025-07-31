package nats

import (
	"context"
	"strconv"
	"ticket/internal/application/core/domain"
	"ticket/utils"

	"github.com/nats-io/nats.go"
)

const (
	SubjectGetTrain          = "request.trains.get"
)

type TicketRequestSenderAdapter struct {
	natsConn *nats.Conn
}

func NewTicketRequestSenderAdapter(natsConn *nats.Conn) *TicketRequestSenderAdapter {
	return &TicketRequestSenderAdapter{
		natsConn: natsConn,
	}
}

func (t *TicketRequestSenderAdapter) RequestGetTrainByID(ctx context.Context, trainID uint) (*domain.Train, error) {
	strTrainID := strconv.Itoa(int(trainID))
	msg, err := t.natsConn.RequestWithContext(ctx, SubjectGetTrain, []byte(strTrainID))
	if err != nil {
		return nil, err
	}

	train,err := utils.UnMarshalTrain(msg.Data)
	if err != nil {
		return nil, err
	}

	return  train,nil
}
