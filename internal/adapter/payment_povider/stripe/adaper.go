package stripe

import (
	"context"
	"fmt"
	"github.com/DenisBarabanshchikov/subscription/internal/model"
	"github.com/DenisBarabanshchikov/subscription/internal/port"
)

type adapter struct {
	api Api
}

func NewAdapter(api Api) port.PaymentProvider {
	return &adapter{
		api: api,
	}
}

func (a *adapter) CreateCustomer(ctx context.Context, email string) (string, error) {
	return a.api.CreateCustomer(ctx, email)
}

func (a *adapter) SubscribeCustomer(ctx context.Context, customer model.Customer, plan string) (string, error) {
	price, err := a.getPriceByPlan(ctx, plan)
	if err != nil {
		return "", err
	}
	return a.api.SubscribeCustomer(ctx, customer, price)
}

func (a *adapter) GetSubscriptionStatus(ctx context.Context, subscriptionId string) (string, error) {
	return a.api.GetSubscriptionStatus(ctx, subscriptionId)
}

func (a *adapter) getPriceByPlan(_ context.Context, plan string) (string, error) {
	switch plan {
	case "Core":
		return "price_1QtWUdIGaC2gk9oobOvUwioa", nil
	case "Growth":
		return "price_1QtWcBIGaC2gk9ookwUgcQPj", nil
	case "Premium":
		return "price_1QtWcWIGaC2gk9ooNnWu1RJi", nil
	default:
		return "", fmt.Errorf("unknown plan: %s", plan)
	}
}
