//go:build unit

package stripe_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/DenisBarabanshchikov/subscription/internal/adapter/payment_povider/stripe"
	"github.com/DenisBarabanshchikov/subscription/internal/model"
	"github.com/DenisBarabanshchikov/subscription/internal/port"
)

// 1. Define a mock for the Api interface
type mockApi struct {
	mock.Mock
}

// Implement the methods from stripe.Api interface
func (m *mockApi) CreateCustomer(ctx context.Context, email string) (string, error) {
	args := m.Called(ctx, email)
	return args.String(0), args.Error(1)
}

func (m *mockApi) SubscribeCustomer(ctx context.Context, customer model.Customer, price string) (string, error) {
	args := m.Called(ctx, customer, price)
	return args.String(0), args.Error(1)
}

func (m *mockApi) GetSubscriptionStatus(ctx context.Context, subscriptionId string) (string, error) {
	args := m.Called(ctx, subscriptionId)
	return args.String(0), args.Error(1)
}

// TestNewAdapter checks that NewAdapter returns a port.PaymentProvider implementation
func TestNewAdapter(t *testing.T) {
	mockAPI := new(mockApi)
	provider := stripe.NewAdapter(mockAPI)

	assert.Implements(t, (*port.PaymentProvider)(nil), provider)
}

// TestCreateCustomer checks that CreateCustomer calls api.CreateCustomer
func TestCreateCustomer(t *testing.T) {
	ctx := context.Background()
	mockAPI := new(mockApi)
	provider := stripe.NewAdapter(mockAPI)

	mockAPI.
		On("CreateCustomer", ctx, "test@example.com").
		Return("cus_12345", nil).
		Once()

	custID, err := provider.CreateCustomer(ctx, "test@example.com")

	assert.NoError(t, err)
	assert.Equal(t, "cus_12345", custID)
	mockAPI.AssertExpectations(t)
}

// TestSubscribeCustomer checks plan-price mapping and calls api.SubscribeCustomer
func TestSubscribeCustomer(t *testing.T) {
	ctx := context.Background()
	mockAPI := new(mockApi)
	provider := stripe.NewAdapter(mockAPI)

	customer := model.Customer{
		CustomerId:         "customer-id-1",
		ExternalCustomerId: "external-customer-id-1",
	}

	// 1. Test known plan
	mockAPI.
		On("SubscribeCustomer", ctx, customer, "price_1QtWUdIGaC2gk9oobOvUwioa").
		Return("sub_9876", nil).
		Once()

	subID, err := provider.SubscribeCustomer(ctx, customer, "Core")
	assert.NoError(t, err)
	assert.Equal(t, "sub_9876", subID)
	mockAPI.AssertExpectations(t)

	// 2. Test unknown plan -> expect error
	_, err = provider.SubscribeCustomer(ctx, customer, "NonExistentPlan")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unknown plan: NonExistentPlan")
}

// TestSubscribeCustomerAPIFailure checks when api.SubscribeCustomer returns an error
func TestSubscribeCustomerAPIFailure(t *testing.T) {
	ctx := context.Background()
	mockAPI := new(mockApi)
	provider := stripe.NewAdapter(mockAPI)

	customer := model.Customer{
		CustomerId:         "customer-id-1",
		ExternalCustomerId: "external-customer-id-1",
	}

	// For example, "Premium" plan
	mockAPI.
		On("SubscribeCustomer", ctx, customer, "price_1QtWcWIGaC2gk9ooNnWu1RJi").
		Return("", errors.New("api failure")).
		Once()

	_, err := provider.SubscribeCustomer(ctx, customer, "Premium")
	assert.Error(t, err)
	assert.Equal(t, "api failure", err.Error())
	mockAPI.AssertExpectations(t)
}

// TestGetSubscriptionStatus ensures the adapter calls api.GetSubscriptionStatus
func TestGetSubscriptionStatus(t *testing.T) {
	ctx := context.Background()
	mockAPI := new(mockApi)
	provider := stripe.NewAdapter(mockAPI)

	mockAPI.
		On("GetSubscriptionStatus", ctx, "sub_123").
		Return("active", nil).
		Once()

	status, err := provider.GetSubscriptionStatus(ctx, "sub_123")

	assert.NoError(t, err)
	assert.Equal(t, "active", status)
	mockAPI.AssertExpectations(t)
}
