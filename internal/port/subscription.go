package port

import (
	"context"
	"github.com/DenisBarabanshchikov/subscription/internal/model"
)

type Subscription interface {
	CreateCustomer(ctx context.Context, customer model.Customer) error
	GetCustomer(ctx context.Context, id string) (*model.Customer, error)
	CreateSubscription(ctx context.Context, subscription model.Subscription) error
	GetSubscription(ctx context.Context, customerId, subscriptionId string) (*model.Subscription, error)
}
