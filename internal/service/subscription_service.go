package service

import (
	"context"
	"github.com/DenisBarabanshchikov/subscription/internal/model"
	"github.com/DenisBarabanshchikov/subscription/internal/port"
	"github.com/DenisBarabanshchikov/subscription/pkg/uuid"
)

type SubscriptionService interface {
	CreateCustomer(ctx context.Context, customerEmail string) (model.Customer, error)
	SubscriberCustomer(ctx context.Context, customerId, plan string) (model.Subscription, error)
	SubscriptionStatus(ctx context.Context, customerId, subscriptionId string) (model.Subscription, error)
}

type subscriptionService struct {
	customer        port.Subscription
	paymentProvider port.PaymentProvider
}

func NewSubscriptionService(customer port.Subscription, paymentProvider port.PaymentProvider) SubscriptionService {
	return &subscriptionService{
		customer:        customer,
		paymentProvider: paymentProvider,
	}
}

func (s subscriptionService) CreateCustomer(ctx context.Context, customerEmail string) (model.Customer, error) {
	externalCustomerId, err := s.paymentProvider.CreateCustomer(ctx, customerEmail)
	if err != nil {
		return model.Customer{}, err
	}
	customer := model.Customer{
		CustomerId:         uuid.GenerateUUID(),
		ExternalCustomerId: externalCustomerId,
	}
	err = s.customer.CreateCustomer(ctx, customer)
	if err != nil {
		return model.Customer{}, err
	}

	return customer, nil
}

func (s subscriptionService) SubscriberCustomer(ctx context.Context, customerId, plan string) (model.Subscription, error) {
	customer, err := s.customer.GetCustomer(context.Background(), customerId)
	if err != nil {
		return model.Subscription{}, err
	}
	if customer == nil {
		return model.Subscription{}, model.NewCustomerNotFoundErr(customerId)
	}
	subscriptionId, err := s.paymentProvider.SubscribeCustomer(ctx, *customer, plan)
	if err != nil {
		return model.Subscription{}, err
	}
	subscription := model.Subscription{
		SubscriptionId:         uuid.GenerateUUID(),
		CustomerId:             customer.CustomerId,
		ExternalSubscriptionID: subscriptionId,
		Plan:                   plan,
		Status:                 "new",
	}
	err = s.customer.CreateSubscription(ctx, subscription)
	if err != nil {
		return model.Subscription{}, err
	}

	return subscription, nil
}

func (s subscriptionService) SubscriptionStatus(ctx context.Context, customerId, subscriptionId string) (model.Subscription, error) {
	customer, err := s.customer.GetCustomer(context.Background(), customerId)
	if err != nil {
		return model.Subscription{}, err
	}
	if customer == nil {
		return model.Subscription{}, model.NewCustomerNotFoundErr(customerId)
	}
	subscription, err := s.customer.GetSubscription(ctx, customerId, subscriptionId)
	if err != nil {
		return model.Subscription{}, err
	}
	if subscription == nil {
		return model.Subscription{}, model.NewSubscriptionNotFoundErr(subscriptionId)
	}

	status, err := s.paymentProvider.GetSubscriptionStatus(ctx, subscription.ExternalSubscriptionID)
	if err != nil {
		return model.Subscription{}, err
	}
	subscription.Status = status

	//TODO update status in DB if this is different

	return *subscription, nil
}
