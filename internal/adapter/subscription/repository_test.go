//go:build integration

package subscription_test

import (
	"context"
	"fmt"
	"github.com/DenisBarabanshchikov/subscription/config"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/DenisBarabanshchikov/subscription/internal/adapter/subscription"
)

func TestDynamoRepository_CreateAndGetCustomer(t *testing.T) {
	repo := subscription.NewDynamoRepository(config.ProvideSubscriptionDynamoConfig())

	// Create a unique test customer.
	customerId := fmt.Sprintf("testcust-%d", time.Now().UnixNano())
	cust := subscription.Customer{
		CustomerId:         customerId,
		ExternalCustomerId: "external-" + customerId,
		CreatedAt:          time.Now().UTC(),
		UpdatedAt:          time.Now().UTC(),
	}

	ctx := context.Background()

	// Insert the customer.
	err := repo.CreateCustomer(ctx, cust)
	assert.NoError(t, err, "failed to create customer")

	// Retrieve the customer.
	retrieved, err := repo.GetCustomer(ctx, customerId)
	assert.NoError(t, err, "failed to get customer")
	assert.NotNil(t, retrieved, "customer not found")
	assert.Equal(t, cust.CustomerId, retrieved.CustomerId)
	assert.Equal(t, cust.ExternalCustomerId, retrieved.ExternalCustomerId)
}

func TestDynamoRepository_CreateAndGetSubscription(t *testing.T) {
	repo := subscription.NewDynamoRepository(config.ProvideSubscriptionDynamoConfig())

	// For a subscription, we first need to ensure a customer exists.
	customerId := fmt.Sprintf("testcust-%d", time.Now().UnixNano())
	cust := subscription.Customer{
		CustomerId:         customerId,
		ExternalCustomerId: "external-" + customerId,
		CreatedAt:          time.Now().UTC(),
		UpdatedAt:          time.Now().UTC(),
	}
	ctx := context.Background()
	err := repo.CreateCustomer(ctx, cust)
	assert.NoError(t, err, "failed to create customer for subscription")

	// Create a unique test subscription.
	subscriptionId := fmt.Sprintf("testsub-%d", time.Now().UnixNano())
	sub := subscription.Subscription{
		SubscriptionId: subscriptionId,
		CustomerId:     customerId,
		Status:         "incomplete",
		CreatedAt:      time.Now().UTC(),
		UpdatedAt:      time.Now().UTC(),
	}

	// Insert the subscription.
	err = repo.CreateSubscription(ctx, sub)
	assert.NoError(t, err, "failed to create subscription")

	// Retrieve the subscription.
	retrievedSub, err := repo.GetSubscription(ctx, customerId, subscriptionId)
	assert.NoError(t, err, "failed to get subscription")
	assert.NotNil(t, retrievedSub, "subscription not found")
	assert.Equal(t, sub.SubscriptionId, retrievedSub.SubscriptionId)
	assert.Equal(t, sub.CustomerId, retrievedSub.CustomerId)
	assert.Equal(t, sub.Status, retrievedSub.Status)
}
