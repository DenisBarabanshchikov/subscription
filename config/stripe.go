package config

import (
	"github.com/DenisBarabanshchikov/subscription/pkg/env"
	"github.com/stripe/stripe-go/v74/client"
)

func ProvideStripeClient() *client.API {
	stripeKey := env.RequiredString("STRIPE_SECRET_KEY")

	sc := &client.API{}
	sc.Init(stripeKey, nil)

	return sc
}
