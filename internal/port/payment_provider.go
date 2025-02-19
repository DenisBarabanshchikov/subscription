package port

import (
	"context"
	"github.com/DenisBarabanshchikov/subscription/internal/model"
)

type PaymentProvider interface {
	CreateCustomer(ctx context.Context, email string) (string, error)
	SubscribeCustomer(ctx context.Context, customer model.Customer, plan string) (string, error)
	GetSubscriptionStatus(ctx context.Context, subscriptionId string) (string, error)
}
