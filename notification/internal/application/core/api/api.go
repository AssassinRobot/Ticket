package api

import (
	"context"
	"notification/internal/ports"
)

type APIAdapter struct {
	notifier ports.Notify
}

func NewAPIAdapter(notifier ports.Notify) *APIAdapter {
	return &APIAdapter{
		notifier: notifier,
	}
}

func (a *APIAdapter) UserRegistered(ctx context.Context, name, email string) error {
	return a.notifier.NotifyUserRegistered(ctx, name, email)
}

func (a *APIAdapter) UserUpdated(ctx context.Context, name, email string) error {
	return a.notifier.NotifyUserUpdated(ctx, name, email)
}
