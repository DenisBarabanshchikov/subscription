package stripe

import (
	"context"
	"github.com/DenisBarabanshchikov/subscription/internal/model"
	"github.com/stripe/stripe-go/v74"
	"github.com/stripe/stripe-go/v74/client"
)

type Api interface {
	CreateCustomer(ctx context.Context, email string) (string, error)
	SubscribeCustomer(ctx context.Context, customer model.Customer, price string) (string, error)
	GetSubscriptionStatus(_ context.Context, subscriptionId string) (string, error)
}

type api struct {
	client *client.API
}

func NewApi(client *client.API) Api {
	return &api{
		client: client,
	}
}

func (a *api) CreateCustomer(_ context.Context, email string) (string, error) {
	customer, err := a.client.Customers.New(&stripe.CustomerParams{
		Email: stripe.String(email),
	})
	if err != nil {
		return "", err
	}

	return customer.ID, nil
}

func (a *api) SubscribeCustomer(_ context.Context, customer model.Customer, price string) (string, error) {
	subParams := &stripe.SubscriptionParams{
		Customer: stripe.String(customer.ExternalCustomerId),
		Items: []*stripe.SubscriptionItemsParams{
			{
				Price: stripe.String(price),
			},
		},
		PaymentBehavior: stripe.String("default_incomplete"),
	}
	subscription, err := a.client.Subscriptions.New(subParams)
	if err != nil {
		return "", err
	}

	return subscription.ID, nil
}

func (a *api) GetSubscriptionStatus(_ context.Context, subscriptionId string) (string, error) {
	subscription, err := a.client.Subscriptions.Get(subscriptionId, nil)
	if err != nil {
		return "", err
	}

	return string(subscription.Status), nil
}
