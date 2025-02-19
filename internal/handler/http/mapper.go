package http

import (
	"github.com/DenisBarabanshchikov/subscription/internal/handler/http/response"
	"github.com/DenisBarabanshchikov/subscription/internal/model"
)

func mapToCreateCustomerResponse(customer model.Customer) response.CreateCustomer {
	return response.CreateCustomer{
		CustomerId:         customer.CustomerId,
		ExternalCustomerId: customer.ExternalCustomerId,
	}
}

func mapToSubscriberCustomerResponse(subscription model.Subscription) response.SubscribeCustomer {
	return response.SubscribeCustomer{
		SubscriptionId:         subscription.SubscriptionId,
		ExternalSubscriptionId: subscription.ExternalSubscriptionID,
	}
}

func mapToSubscriptionStatusResponse(subscription model.Subscription) response.SubscriptionStatus {
	return response.SubscriptionStatus{
		SubscriptionId:         subscription.SubscriptionId,
		ExternalSubscriptionID: subscription.ExternalSubscriptionID,
		Plan:                   subscription.Plan,
		Status:                 subscription.Status,
	}
}
