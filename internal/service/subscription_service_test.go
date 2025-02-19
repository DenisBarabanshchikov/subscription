//go:build unit

package service_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/DenisBarabanshchikov/subscription/internal/model"
	"github.com/DenisBarabanshchikov/subscription/internal/service"
)

// --- Mocks ---

// mockSubscription implements port.Subscription.
type mockSubscription struct {
	mock.Mock
}

func (m *mockSubscription) CreateCustomer(ctx context.Context, customer model.Customer) error {
	args := m.Called(ctx, customer)
	return args.Error(0)
}

func (m *mockSubscription) GetCustomer(ctx context.Context, id string) (*model.Customer, error) {
	args := m.Called(ctx, id)
	if cust, ok := args.Get(0).(*model.Customer); ok {
		return cust, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *mockSubscription) CreateSubscription(ctx context.Context, subscription model.Subscription) error {
	args := m.Called(ctx, subscription)
	return args.Error(0)
}

func (m *mockSubscription) GetSubscription(ctx context.Context, customerId, subscriptionId string) (*model.Subscription, error) {
	args := m.Called(ctx, customerId, subscriptionId)
	if sub, ok := args.Get(0).(*model.Subscription); ok {
		return sub, args.Error(1)
	}
	return nil, args.Error(1)
}

// mockPaymentProvider implements port.PaymentProvider.
type mockPaymentProvider struct {
	mock.Mock
}

func (m *mockPaymentProvider) CreateCustomer(ctx context.Context, email string) (string, error) {
	args := m.Called(ctx, email)
	return args.String(0), args.Error(1)
}

func (m *mockPaymentProvider) SubscribeCustomer(ctx context.Context, customer model.Customer, plan string) (string, error) {
	args := m.Called(ctx, customer, plan)
	return args.String(0), args.Error(1)
}

func (m *mockPaymentProvider) GetSubscriptionStatus(ctx context.Context, subscriptionId string) (string, error) {
	args := m.Called(ctx, subscriptionId)
	return args.String(0), args.Error(1)
}

// --- Unit Tests ---

func TestCreateCustomer_Success(t *testing.T) {
	ctx := context.Background()

	mockSub := new(mockSubscription)
	mockPay := new(mockPaymentProvider)

	email := "test@mail.com"
	externalCustomerID := "ext_cus_123"

	// Expect the payment provider to create the customer and return an external ID.
	mockPay.
		On("CreateCustomer", ctx, email).
		Return(externalCustomerID, nil).Once()

	// The service calls subscription.CreateCustomer with a model.Customer that has ExternalCustomerId set.
	// Since CustomerId is generated (via uuid.GenerateUUID), we only verify that it is non-empty
	// and that ExternalCustomerId is as expected.
	mockSub.
		On("CreateCustomer", ctx, mock.MatchedBy(func(c model.Customer) bool {
			return c.ExternalCustomerId == externalCustomerID && c.CustomerId != ""
		})).
		Return(nil).Once()

	svc := service.NewSubscriptionService(mockSub, mockPay)
	cust, err := svc.CreateCustomer(ctx, email)
	assert.NoError(t, err)
	assert.Equal(t, externalCustomerID, cust.ExternalCustomerId)
	assert.NotEmpty(t, cust.CustomerId)

	mockPay.AssertExpectations(t)
	mockSub.AssertExpectations(t)
}

func TestCreateCustomer_PaymentProviderError(t *testing.T) {
	ctx := context.Background()

	mockSub := new(mockSubscription)
	mockPay := new(mockPaymentProvider)

	email := "test@mail.com"
	expectedErr := errors.New("payment provider error")

	mockPay.
		On("CreateCustomer", ctx, email).
		Return("", expectedErr).Once()

	svc := service.NewSubscriptionService(mockSub, mockPay)
	cust, err := svc.CreateCustomer(ctx, email)
	assert.Error(t, err)
	assert.Equal(t, expectedErr, err)
	assert.Empty(t, cust.CustomerId)

	// No call expected to subscription.CreateCustomer
	mockPay.AssertExpectations(t)
}

func TestSubscriberCustomer_Success(t *testing.T) {
	ctx := context.Background()

	mockSub := new(mockSubscription)
	mockPay := new(mockPaymentProvider)

	customerId := "cust_123"
	plan := "Core"
	externalSubID := "ext_sub_456"

	// Prepare an existing customer.
	existingCustomer := &model.Customer{
		CustomerId:         customerId,
		ExternalCustomerId: "ext_cus_123",
	}
	// When GetCustomer is called, return our existing customer.
	mockSub.
		On("GetCustomer", mock.Anything, customerId).
		Return(existingCustomer, nil).Once()

	// Expect the payment provider to subscribe the customer.
	mockPay.
		On("SubscribeCustomer", ctx, *existingCustomer, plan).
		Return(externalSubID, nil).Once()

	// Expect CreateSubscription to be called with a subscription that has the proper fields.
	mockSub.
		On("CreateSubscription", ctx, mock.MatchedBy(func(s model.Subscription) bool {
			return s.CustomerId == customerId &&
				s.ExternalSubscriptionID == externalSubID &&
				s.Plan == plan &&
				s.Status == "new" &&
				s.SubscriptionId != ""
		})).
		Return(nil).Once()

	svc := service.NewSubscriptionService(mockSub, mockPay)
	sub, err := svc.SubscriberCustomer(ctx, customerId, plan)
	assert.NoError(t, err)
	assert.Equal(t, customerId, sub.CustomerId)
	assert.Equal(t, externalSubID, sub.ExternalSubscriptionID)
	assert.Equal(t, plan, sub.Plan)
	assert.Equal(t, "new", sub.Status)
	assert.NotEmpty(t, sub.SubscriptionId)

	mockSub.AssertExpectations(t)
	mockPay.AssertExpectations(t)
}

func TestSubscriberCustomer_CustomerNotFound(t *testing.T) {
	ctx := context.Background()

	mockSub := new(mockSubscription)
	mockPay := new(mockPaymentProvider)

	nonExistentCustomerID := "nonexistent"
	plan := "Core"

	// Return nil to simulate that the customer was not found.
	mockSub.
		On("GetCustomer", mock.Anything, nonExistentCustomerID).
		Return(nil, nil).Once()

	svc := service.NewSubscriptionService(mockSub, mockPay)
	sub, err := svc.SubscriberCustomer(ctx, nonExistentCustomerID, plan)
	assert.Error(t, err)
	assert.Equal(t, model.NewCustomerNotFoundErr(nonExistentCustomerID).Error(), err.Error())
	assert.Empty(t, sub.SubscriptionId)

	mockSub.AssertExpectations(t)
}

func TestSubscriberCustomer_PaymentProviderError(t *testing.T) {
	ctx := context.Background()

	mockSub := new(mockSubscription)
	mockPay := new(mockPaymentProvider)

	customerId := "cust_123"
	plan := "Core"
	existingCustomer := &model.Customer{
		CustomerId:         customerId,
		ExternalCustomerId: "ext_cus_123",
	}
	expectedErr := errors.New("subscribe error")

	mockSub.
		On("GetCustomer", mock.Anything, customerId).
		Return(existingCustomer, nil).Once()

	mockPay.
		On("SubscribeCustomer", ctx, *existingCustomer, plan).
		Return("", expectedErr).Once()

	svc := service.NewSubscriptionService(mockSub, mockPay)
	sub, err := svc.SubscriberCustomer(ctx, customerId, plan)
	assert.Error(t, err)
	assert.Equal(t, expectedErr, err)
	assert.Empty(t, sub.SubscriptionId)

	mockSub.AssertExpectations(t)
	mockPay.AssertExpectations(t)
}

func TestSubscriberCustomer_CreateSubscriptionError(t *testing.T) {
	ctx := context.Background()

	mockSub := new(mockSubscription)
	mockPay := new(mockPaymentProvider)

	customerId := "cust_123"
	plan := "Core"
	externalSubID := "ext_sub_456"
	expectedErr := errors.New("create subscription error")

	existingCustomer := &model.Customer{
		CustomerId:         customerId,
		ExternalCustomerId: "ext_cus_123",
	}
	mockSub.
		On("GetCustomer", mock.Anything, customerId).
		Return(existingCustomer, nil).Once()

	mockPay.
		On("SubscribeCustomer", ctx, *existingCustomer, plan).
		Return(externalSubID, nil).Once()

	mockSub.
		On("CreateSubscription", ctx, mock.Anything).
		Return(expectedErr).Once()

	svc := service.NewSubscriptionService(mockSub, mockPay)
	sub, err := svc.SubscriberCustomer(ctx, customerId, plan)
	assert.Error(t, err)
	assert.Equal(t, expectedErr, err)
	assert.Empty(t, sub.SubscriptionId)

	mockSub.AssertExpectations(t)
	mockPay.AssertExpectations(t)
}

func TestSubscriptionStatus_Success(t *testing.T) {
	ctx := context.Background()

	mockSub := new(mockSubscription)
	mockPay := new(mockPaymentProvider)

	customerId := "cust_123"
	subscriptionId := "sub_abc"
	externalSubID := "ext_sub_789"
	status := "active"

	existingCustomer := &model.Customer{
		CustomerId:         customerId,
		ExternalCustomerId: "ext_cus_123",
	}
	existingSubscription := &model.Subscription{
		SubscriptionId:         subscriptionId,
		CustomerId:             customerId,
		ExternalSubscriptionID: externalSubID,
		Status:                 "new",
	}

	mockSub.
		On("GetCustomer", mock.Anything, customerId).
		Return(existingCustomer, nil).Once()

	mockSub.
		On("GetSubscription", ctx, customerId, subscriptionId).
		Return(existingSubscription, nil).Once()

	mockPay.
		On("GetSubscriptionStatus", ctx, externalSubID).
		Return(status, nil).Once()

	svc := service.NewSubscriptionService(mockSub, mockPay)
	sub, err := svc.SubscriptionStatus(ctx, customerId, subscriptionId)
	assert.NoError(t, err)
	assert.Equal(t, status, sub.Status)

	mockSub.AssertExpectations(t)
	mockPay.AssertExpectations(t)
}

func TestSubscriptionStatus_CustomerNotFound(t *testing.T) {
	ctx := context.Background()

	mockSub := new(mockSubscription)
	mockPay := new(mockPaymentProvider)

	customerId := "nonexistent"
	subscriptionId := "sub_abc"

	mockSub.
		On("GetCustomer", mock.Anything, customerId).
		Return(nil, nil).Once()

	svc := service.NewSubscriptionService(mockSub, mockPay)
	sub, err := svc.SubscriptionStatus(ctx, customerId, subscriptionId)
	assert.Error(t, err)
	assert.Equal(t, model.NewCustomerNotFoundErr(customerId).Error(), err.Error())
	assert.Empty(t, sub.SubscriptionId)

	mockSub.AssertExpectations(t)
}

func TestSubscriptionStatus_SubscriptionNotFound(t *testing.T) {
	ctx := context.Background()

	mockSub := new(mockSubscription)
	mockPay := new(mockPaymentProvider)

	customerId := "cust_123"
	subscriptionId := "nonexistent"

	existingCustomer := &model.Customer{
		CustomerId:         customerId,
		ExternalCustomerId: "ext_cus_123",
	}
	mockSub.
		On("GetCustomer", mock.Anything, customerId).
		Return(existingCustomer, nil).Once()

	// Return nil for GetSubscription to simulate subscription not found.
	mockSub.
		On("GetSubscription", ctx, customerId, subscriptionId).
		Return(nil, nil).Once()

	svc := service.NewSubscriptionService(mockSub, mockPay)
	sub, err := svc.SubscriptionStatus(ctx, customerId, subscriptionId)
	assert.Error(t, err)
	assert.Equal(t, model.NewSubscriptionNotFoundErr(subscriptionId).Error(), err.Error())
	assert.Empty(t, sub.SubscriptionId)

	mockSub.AssertExpectations(t)
}

func TestSubscriptionStatus_PaymentProviderError(t *testing.T) {
	ctx := context.Background()

	mockSub := new(mockSubscription)
	mockPay := new(mockPaymentProvider)

	customerId := "cust_123"
	subscriptionId := "sub_abc"
	externalSubID := "ext_sub_789"
	expectedErr := errors.New("payment provider error")

	existingCustomer := &model.Customer{
		CustomerId:         customerId,
		ExternalCustomerId: "ext_cus_123",
	}
	existingSubscription := &model.Subscription{
		SubscriptionId:         subscriptionId,
		CustomerId:             customerId,
		ExternalSubscriptionID: externalSubID,
		Status:                 "new",
	}

	mockSub.
		On("GetCustomer", mock.Anything, customerId).
		Return(existingCustomer, nil).Once()

	mockSub.
		On("GetSubscription", ctx, customerId, subscriptionId).
		Return(existingSubscription, nil).Once()

	mockPay.
		On("GetSubscriptionStatus", ctx, externalSubID).
		Return("", expectedErr).Once()

	svc := service.NewSubscriptionService(mockSub, mockPay)
	sub, err := svc.SubscriptionStatus(ctx, customerId, subscriptionId)
	assert.Error(t, err)
	assert.Equal(t, expectedErr, err)
	assert.Empty(t, sub.SubscriptionId)

	mockSub.AssertExpectations(t)
	mockPay.AssertExpectations(t)
}
