package subscription

import (
	"github.com/DenisBarabanshchikov/subscription/internal/model"
	"time"
)

func mapToCustomerEntity(customer model.Customer) Customer {
	return Customer{
		CustomerId:         customer.CustomerId,
		ExternalCustomerId: customer.ExternalCustomerId,
		CreatedAt:          time.Now(),
		UpdatedAt:          time.Now(),
	}
}

func mapToCustomerModel(customer Customer) model.Customer {
	return model.Customer{
		CustomerId:         customer.CustomerId,
		ExternalCustomerId: customer.ExternalCustomerId,
	}
}

func mapToCustomerModelPtr(customer *Customer) *model.Customer {
	if customer == nil {
		return nil
	}
	res := mapToCustomerModel(*customer)
	return &res
}

func mapSubscriptionToEntity(subscription model.Subscription) Subscription {
	return Subscription{
		SubscriptionId:         subscription.SubscriptionId,
		CustomerId:             subscription.CustomerId,
		ExternalSubscriptionID: subscription.ExternalSubscriptionID,
		Plan:                   subscription.Plan,
		Status:                 subscription.Status,
		CreatedAt:              time.Now(),
		UpdatedAt:              time.Now(),
	}
}

func mapSubscriptionToModel(subscription Subscription) model.Subscription {
	return model.Subscription{
		SubscriptionId:         subscription.SubscriptionId,
		CustomerId:             subscription.CustomerId,
		ExternalSubscriptionID: subscription.ExternalSubscriptionID,
		Plan:                   subscription.Plan,
		Status:                 subscription.Status,
	}
}

func mapSubscriptionToModelPtr(subscription *Subscription) *model.Subscription {
	if subscription == nil {
		return nil
	}
	res := mapSubscriptionToModel(*subscription)
	return &res
}
