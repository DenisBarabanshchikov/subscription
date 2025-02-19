package subscription

import (
	"context"
	"github.com/DenisBarabanshchikov/subscription/internal/model"
	"github.com/DenisBarabanshchikov/subscription/internal/port"
)

type adapter struct {
	repository Repository
}

func NewAdapter(repository Repository) port.Subscription {
	return &adapter{
		repository: repository,
	}
}

func (a *adapter) CreateCustomer(ctx context.Context, customer model.Customer) error {
	return a.repository.CreateCustomer(ctx, mapToCustomerEntity(customer))
}

func (a *adapter) GetCustomer(ctx context.Context, id string) (*model.Customer, error) {
	customer, err := a.repository.GetCustomer(ctx, id)
	if err != nil {
		return nil, err
	}
	return mapToCustomerModelPtr(customer), nil
}

func (a *adapter) CreateSubscription(ctx context.Context, subscription model.Subscription) error {
	return a.repository.CreateSubscription(ctx, mapSubscriptionToEntity(subscription))
}

func (a *adapter) GetSubscription(ctx context.Context, customerId, subscriptionId string) (*model.Subscription, error) {
	subscription, err := a.repository.GetSubscription(ctx, customerId, subscriptionId)
	if err != nil {
		return nil, err
	}
	return mapSubscriptionToModelPtr(subscription), nil
}
