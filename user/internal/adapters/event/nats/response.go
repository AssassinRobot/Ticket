package nats

import (
	"context"
	"fmt"
	"strconv"
	"user/internal/ports"
	"user/utils"

	"github.com/nats-io/nats.go"
)

const (
	SubjectRequestListUsers      = "request.list.users"
	SubjectRequestGetUserByID    = "request.get.user.byID"
	SubjectRequestCreateUser     = "request.create.user"
	SubjectRequestUpdateUser     = "request.update.user"
	SubjectRequestDeleteUser     = "request.delete.user"
)

type UserEventResponderAdapter struct {
	natsConn *nats.Conn
	API      ports.APIPort
}

func NewUserEventResponderAdapter(natsConn *nats.Conn, api ports.APIPort) *UserEventResponderAdapter {
	return &UserEventResponderAdapter{
		natsConn: natsConn,
		API:      api,
	}
}

func (u *UserEventResponderAdapter) ReplyToListUsers(ctx context.Context) error {
	subscription, err := u.natsConn.Subscribe(SubjectRequestListUsers, func(msg *nats.Msg) {
		users, err := u.API.ListUsers(ctx)
		if err != nil {
			errString := fmt.Errorf("error lists users by ID  : %w", err).Error()
			msg.Respond([]byte(errString))
			return
		}

		serializedData, err := utils.MarshalListUsersReplay(users)
		if err != nil {
			fmt.Println()
			errString := fmt.Errorf("error serialize users: %v,%w", users, err).Error()
			msg.Respond([]byte(errString))
			return
		}

		msg.Respond(serializedData)
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

func (u *UserEventResponderAdapter) ReplyToGetUserByID(ctx context.Context) error {
	subscription, err := u.natsConn.Subscribe(SubjectRequestGetUserByID, func(msg *nats.Msg) {
		userIDstr := string(msg.Data)
		userID, err := strconv.Atoi(userIDstr)
		if err != nil {
			errString := fmt.Errorf("error invalid userID: %d", userID).Error()

			msg.Respond([]byte(errString))
			return
		}

		user, err := u.API.GetUserByID(ctx, uint(userID))
		if err != nil {
			errString := fmt.Errorf("error get user by ID  : %d,%w", userID, err).Error()
			msg.Respond([]byte(errString))
			return
		}

		serializedUserData, err := utils.MarshalUser(user)
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

func (u *UserEventResponderAdapter) ReplayToCreateUser(ctx context.Context) error {
	subscription, err := u.natsConn.Subscribe(SubjectRequestCreateUser, func(msg *nats.Msg) {
		firstName, lastName, email, err := utils.UnmarshalCreateUserRequest(msg.Data)
		if err != nil {
			errString := fmt.Errorf("error unmarshal data  : %s,%w", string(msg.Data), err).Error()
			msg.Respond([]byte(errString))
			return
		}

		_, err = u.API.Register(ctx, firstName, lastName, email)
		if err != nil {
			errString := fmt.Errorf("error create user: %w", err).Error()
			msg.Respond([]byte(errString))
			return
		}

		msg.Respond([]byte("Ok"))
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

func (u *UserEventResponderAdapter) ReplayToUpdateUserByID(ctx context.Context) error {
	subscription, err := u.natsConn.Subscribe(SubjectRequestUpdateUser, func(msg *nats.Msg) {
		userID, firstName, lastName, err := utils.UnmarshalUpdateUserRequest(msg.Data)
		if err != nil {
			errString := fmt.Errorf("error unmarshal data : %s,%w", string(msg.Data), err).Error()
			msg.Respond([]byte(errString))
			return
		}

		err = u.API.UpdateUser(ctx, userID, firstName, lastName)
		if err != nil {
			errString := fmt.Errorf("error update user: %w", err).Error()
			msg.Respond([]byte(errString))
			return
		}

		msg.Respond([]byte("Ok"))
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

func (u *UserEventResponderAdapter) ReplayToDeleteUserByID(ctx context.Context) error {
	subscription, err := u.natsConn.Subscribe(SubjectRequestDeleteUser, func(msg *nats.Msg) {
		userIDstr := string(msg.Data)
		userID, err := strconv.Atoi(userIDstr)
		if err != nil {
			errString := fmt.Errorf("error invalid userID: %d", userID).Error()

			msg.Respond([]byte(errString))
			return
		}

		err = u.API.DeleteUser(ctx, uint(userID))
		if err != nil {
			errString := fmt.Errorf("error delete user %d: %w",userID, err).Error()
			msg.Respond([]byte(errString))
			return
		}

		msg.Respond([]byte("Ok"))
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
