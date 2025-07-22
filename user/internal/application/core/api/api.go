package api

import (
	"context"
	"user/internal/application/core/domain"
	"user/internal/ports"
)

type API struct {
	ports.DatabasePort
	ports.UserEventPublisher
	ports.UserEventResponder
}

func NewAPI(
	databasePort ports.DatabasePort,
	userEventPublisher ports.UserEventPublisher,
) *API {
	return &API{
		DatabasePort:        databasePort,
		UserEventPublisher:  userEventPublisher,
	}
}	

func (api *API) 	Register(ctx context.Context, firstName, lastName, email string) (*domain.User, error) {
	user := &domain.User{
		FirstName: firstName,
		LastName:  lastName,
		Email:     email,
	}

	if err := api.DatabasePort.SaveUser(ctx, user); err != nil {
		return nil, err
	}

	err := api.UserEventPublisher.PublishUserRegistered(ctx, user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (api *API) GetUserByID(ctx context.Context, id uint) (*domain.User, error) {
	return api.DatabasePort.GetUserByID(ctx, id)
}

func (api *API) ListUsers(ctx context.Context) ([]domain.User, error) {
	users, err := api.DatabasePort.ListUsers(ctx)
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (api *API) UpdateUser(ctx context.Context, id uint, firstName, lastName string) error {
	user,err := api.DatabasePort.UpdateUser(ctx, id, firstName, lastName)
	if err != nil {
		return err
	}

	err = api.UserEventPublisher.PublishUserUpdated(ctx,user)
	if err != nil {
		return err
	}

	return nil
}

func (api *API) DeleteUser(ctx context.Context, id uint) error {
	return api.DatabasePort.DeleteUser(ctx, id)
}
