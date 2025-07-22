package nats

import (
	"context"
	"fmt"
	"strconv"
	"user/internal/application/core/domain"
	"user/utils"

	"github.com/nats-io/nats.go"
)

const (
	SubjectGetUser = "request.user.get"
)

type UserEventResponderAdapter struct {
	natsConn       *nats.Conn
	getUserHandler func(ctx context.Context, userID uint) (*domain.User, error)
}

func NewUserEventResponderAdapter(natsConn *nats.Conn, getUserHandler func(ctx context.Context, userID uint) (*domain.User, error)) *UserEventResponderAdapter {
	return &UserEventResponderAdapter{
		natsConn:       natsConn,
		getUserHandler: getUserHandler,
	}
}

func (u *UserEventResponderAdapter) ReplyToGetUser(ctx context.Context) error {
	subscription, err := u.natsConn.Subscribe(SubjectGetUser, func(msg *nats.Msg) {
		userIDstr := string(msg.Data)
		userID, err := strconv.Atoi(userIDstr)
		if err != nil {
			errString := fmt.Errorf("error invalid userID: %d", userID).Error()

			msg.Respond([]byte(errString))
			return
		}

		user, err := u.getUserHandler(ctx, uint(userID))
		if err != nil {
			errString := fmt.Errorf("error get user by ID  : %d,%w", userID, err).Error()
			msg.Respond([]byte(errString))
			return
		}

		serializedUserData, err := utils.Marshal(user)
		if err != nil {
			fmt.Println()
			errString := fmt.Errorf("error serialize user: %v,%w", user, err).Error()
			msg.Respond([]byte(errString))
			return
		}

		msg.Respond(serializedUserData)
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
