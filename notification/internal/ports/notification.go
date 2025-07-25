package ports

import "context"

type Notify interface {
	NotifyUserRegistered(ctx context.Context,name,email string)error
	NotifyUserUpdated(ctx context.Context,name,email string)error
}