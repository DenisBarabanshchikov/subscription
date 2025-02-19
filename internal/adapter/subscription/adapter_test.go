//go:build unit

package subscription_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/DenisBarabanshchikov/subscription/internal/adapter/subscription"
	"github.com/DenisBarabanshchikov/subscription/internal/model"
)

// mockRepository is a mock of the subscription.Repository interface
type mockRepository struct {
	mock.Mock
}

func (m *mockRepository) CreateCustomer(ctx context.Context, customer subscription.Customer) error {
	args := m.Called(ctx, customer)
	return args.Error(0)
}

func (m *mockRepository) GetCustomer(ctx context.Context, id string) (*subscription.Customer, error) {
	args := m.Called(ctx, id)
	// We use a type assertion. In case of nil, handle appropriately.
	if ce, ok := args.Get(0).(*subscription.Customer); ok {
		return ce, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *mockRepository) CreateSubscription(ctx context.Context, sub subscription.Subscription) error {
	args := m.Called(ctx, sub)
	return args.Error(0)
}

func (m *mockRepository) GetSubscription(ctx context.Context, customerId, subscriptionId string) (*subscription.Subscription, error) {
	args := m.Called(ctx, customerId, subscriptionId)
	if se, ok := args.Get(0).(*subscription.Subscription); ok {
		return se, args.Error(1)
	}
	return nil, args.Error(1)
}

// TestCreateCustomer checks that the adapter calls repo.CreateCustomer with correct data
func TestCreateCustomer(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(mockRepository)

	// Initialize the adapter with our mock repository
	adapter := subscription.NewAdapter(mockRepo)

	// Define a sample domain model.Customer.
	// Note: We leave the CreatedAt/UpdatedAt fields unset in the input.
	cust := model.Customer{
		CustomerId:         "cust_123",
		ExternalCustomerId: "cus_67890",
	}

	// Instead of comparing the entire mapped struct (which includes timestamps),
	// use a custom matcher that checks only the fields of interest.
	mockRepo.
		On("CreateCustomer", ctx, mock.MatchedBy(func(c subscription.Customer) bool {
			return c.CustomerId == "cust_123" &&
				c.ExternalCustomerId == "cus_67890"
		})).
		Return(nil).
		Once()

	// Call the adapter method
	err := adapter.CreateCustomer(ctx, cust)

	// Assertions
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

// TestCreateCustomer_Error checks how the adapter handles repository errors
func TestCreateCustomer_Error(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(mockRepository)
	adapter := subscription.NewAdapter(mockRepo)

	cust := model.Customer{
		CustomerId: "cust_999",
	}

	// Suppose the repository returns an error
	expectedErr := errors.New("repository failure")
	mockRepo.
		On("CreateCustomer", ctx, mock.Anything).
		Return(expectedErr).
		Once()

	err := adapter.CreateCustomer(ctx, cust)
	assert.Equal(t, expectedErr, err)
	mockRepo.AssertExpectations(t)
}

func TestGetCustomer(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(mockRepository)
	adapter := subscription.NewAdapter(mockRepo)

	// Define a test ID
	customerID := "cust_123"

	// The repository returns an entity.
	repoEntity := subscription.Customer{
		CustomerId: "cust_123",
		// For this test, CreatedAt and UpdatedAt can be set to any value.
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	mockRepo.
		On("GetCustomer", ctx, customerID).
		Return(&repoEntity, nil).
		Once()

	cust, err := adapter.GetCustomer(ctx, customerID)
	assert.NoError(t, err)
	// Check that the domain model is mapped from the returned entity
	assert.Equal(t, "cust_123", cust.CustomerId)

	mockRepo.AssertExpectations(t)
}

// TestGetCustomer_NotFound tests error scenario for GetCustomer
func TestGetCustomer_NotFound(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(mockRepository)
	adapter := subscription.NewAdapter(mockRepo)

	missingID := "not_exist"
	repoErr := errors.New("not found")

	// The repository fails to find the customer
	mockRepo.
		On("GetCustomer", ctx, missingID).
		Return((*subscription.Customer)(nil), repoErr).
		Once()

	cust, err := adapter.GetCustomer(ctx, missingID)
	assert.Nil(t, cust)
	assert.Equal(t, repoErr, err)

	mockRepo.AssertExpectations(t)
}

func TestCreateSubscription(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(mockRepository)
	adapter := subscription.NewAdapter(mockRepo)

	// Define a sample domain subscription model.
	sub := model.Subscription{
		SubscriptionId: "sub_abc",
		CustomerId:     "cust_123",
		Status:         "incomplete",
	}

	// Expected entity after mapping. Again, we ignore timestamp fields by using a matcher.
	mockRepo.
		On("CreateSubscription", ctx, mock.MatchedBy(func(s subscription.Subscription) bool {
			return s.SubscriptionId == "sub_abc" &&
				s.CustomerId == "cust_123" &&
				s.Status == "incomplete"
		})).
		Return(nil).
		Once()

	err := adapter.CreateSubscription(ctx, sub)
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestCreateSubscription_Error(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(mockRepository)
	adapter := subscription.NewAdapter(mockRepo)

	sub := model.Subscription{SubscriptionId: "sub_error", CustomerId: "cust_999"}
	repoErr := errors.New("failed to create subscription")

	mockRepo.
		On("CreateSubscription", ctx, mock.Anything).
		Return(repoErr).
		Once()

	err := adapter.CreateSubscription(ctx, sub)
	assert.Equal(t, repoErr, err)
	mockRepo.AssertExpectations(t)
}

func TestGetSubscription(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(mockRepository)
	adapter := subscription.NewAdapter(mockRepo)

	customerID := "cust_123"
	subscriptionID := "sub_456"

	// The repository returns an entity.
	repoEntity := subscription.Subscription{
		SubscriptionId: "sub_456",
		CustomerId:     "cust_123",
		Status:         "incomplete",
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}
	mockRepo.
		On("GetSubscription", ctx, customerID, subscriptionID).
		Return(&repoEntity, nil).
		Once()

	sub, err := adapter.GetSubscription(ctx, customerID, subscriptionID)
	assert.NoError(t, err)
	assert.Equal(t, "sub_456", sub.SubscriptionId)
	assert.Equal(t, "cust_123", sub.CustomerId)
	assert.Equal(t, "incomplete", sub.Status)

	mockRepo.AssertExpectations(t)
}

func TestGetSubscription_Error(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(mockRepository)
	adapter := subscription.NewAdapter(mockRepo)

	custID := "unknown_cust"
	subID := "unknown_sub"
	repoErr := errors.New("not found")

	mockRepo.
		On("GetSubscription", ctx, custID, subID).
		Return((*subscription.Subscription)(nil), repoErr).
		Once()

	sub, err := adapter.GetSubscription(ctx, custID, subID)
	assert.Nil(t, sub)
	assert.Equal(t, repoErr, err)

	mockRepo.AssertExpectations(t)
}
