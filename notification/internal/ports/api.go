package ports

import "context"

type APIPort interface {
	UserRegistered(ctx context.Context,name,email string)error
	UserUpdated(ctx context.Context,name,email string)error
}

