//go:build wireinject
// +build wireinject

package di

import (
	"github.com/DenisBarabanshchikov/subscription/config"
	"github.com/DenisBarabanshchikov/subscription/internal/adapter/payment_povider/stripe"
	"github.com/DenisBarabanshchikov/subscription/internal/adapter/subscription"
	"github.com/DenisBarabanshchikov/subscription/internal/handler/http"
	"github.com/DenisBarabanshchikov/subscription/internal/port"
	"github.com/DenisBarabanshchikov/subscription/internal/service"
	"github.com/google/wire"
	"github.com/stripe/stripe-go/v74/client"
)

var configs = wire.NewSet(
	config.ProvideSubscriptionDynamoConfig,
)

var clients = wire.NewSet(
	config.ProvideStripeClient,
)

var repositories = wire.NewSet(
	subscriptionRepository,
)

var api = wire.NewSet(
	stripeApi,
)

var ports = wire.NewSet(
	subscriptionPort,
	paymentProviderPort,
)

func subscriptionRepository(config subscription.DynamoConfig) subscription.Repository {
	wire.Build(
		subscription.NewDynamoRepository,
	)
	return nil
}

func stripeApi(client *client.API) stripe.Api {
	wire.Build(
		stripe.NewApi,
	)
	return nil
}

func subscriptionPort(repository subscription.Repository) port.Subscription {
	wire.Build(
		subscription.NewAdapter,
	)
	return nil
}

func paymentProviderPort(api stripe.Api) port.PaymentProvider {
	wire.Build(
		stripe.NewAdapter,
	)
	return nil
}

func InitializeHandlers() (*http.Handlers, error) {
	wire.Build(
		configs,
		clients,
		api,
		repositories,
		ports,
		service.NewSubscriptionService,
		http.NewSubscriptionHandler,
		http.NewHandlers,
	)
	return &http.Handlers{}, nil
}
