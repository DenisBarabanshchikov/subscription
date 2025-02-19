//go:build integration

package stripe_test

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"github.com/stripe/stripe-go/v74/client"

	paymentProvider "github.com/DenisBarabanshchikov/subscription/internal/adapter/payment_povider/stripe"
	"github.com/DenisBarabanshchikov/subscription/internal/model"
)

// IntegrationTestSuite is a testify suite for running integration tests against Stripe
type IntegrationTestSuite struct {
	suite.Suite
	api        paymentProvider.Api
	customerID string
	subID      string
}

// SetupSuite runs before the tests in this suite
func (s *IntegrationTestSuite) SetupSuite() {
	// Check that we have a test secret key
	secretKey := os.Getenv("STRIPE_SECRET_KEY")
	if secretKey == "" {
		s.T().Skip("No STRIPE_SECRET_KEY set, skipping Stripe integration tests.")
		return
	}

	// Initialize the stripe client
	sc := &client.API{}
	sc.Init(secretKey, nil)

	// Create the API adapter
	s.api = paymentProvider.NewApi(sc)
}

// TearDownSuite runs after the tests in this suite
func (s *IntegrationTestSuite) TearDownSuite() {
	// Optional: clean up the test data in Stripe by deleting subscription/customer if desired
	// e.g., if s.subID != "" { sc.Subscriptions.Del(s.subID, nil) }
	//       if s.customerID != "" { sc.Customers.Del(s.customerID, nil) }
}

// TestCreateCustomer checks that we can create a real customer in test mode
func (s *IntegrationTestSuite) TestCreateCustomer() {
	ctx := context.Background()

	// Try to create a test customer in Stripe
	custID, err := s.api.CreateCustomer(ctx, "integration-test@mail.com")
	require.NoError(s.T(), err)
	require.NotEmpty(s.T(), custID)

	s.customerID = custID
	s.T().Logf("Created test customer: %s", custID)
}

// TestSubscribeCustomer tries to create a subscription in default_incomplete mode
// Make sure you have a valid test Price ID in your account.
// You can pass it via environment variable or hard-code a test price_XXX for demonstration.
func (s *IntegrationTestSuite) TestSubscribeCustomer() {
	if s.customerID == "" {
		s.T().Skip("No test customer created; skipping TestSubscribeCustomer")
	}

	testPrice := "price_1QtWcWIGaC2gk9ooNnWu1RJi"
	if testPrice == "" {
		s.T().Skip("No STRIPE_TEST_PRICE_ID set, skipping subscription creation test.")
	}

	ctx := context.Background()

	// Build a model.Customer that has the newly created Stripe customer ID
	cust := model.Customer{
		CustomerId:         "internal-id-1",
		ExternalCustomerId: s.customerID, // important
	}

	subID, err := s.api.SubscribeCustomer(ctx, cust, testPrice)
	require.NoError(s.T(), err)
	require.NotEmpty(s.T(), subID)

	s.subID = subID
	s.T().Logf("Created subscription: %s", subID)
}

// TestGetSubscriptionStatus fetches the subscription status and checks if it's "incomplete" by default
func (s *IntegrationTestSuite) TestGetSubscriptionStatus() {
	if s.subID == "" {
		s.T().Skip("No subscription created; skipping TestGetSubscriptionStatus")
	}
	ctx := context.Background()

	status, err := s.api.GetSubscriptionStatus(ctx, s.subID)
	require.NoError(s.T(), err)
	s.T().Logf("Subscription %s status: %s", s.subID, status)

	// Because we used default_incomplete, we expect "incomplete" initially
	require.Equal(s.T(), "incomplete", status)
}

// We use testify's suite runner to run all tests in IntegrationTestSuite in sequence
func TestIntegrationTestSuite(t *testing.T) {
	suite.Run(t, new(IntegrationTestSuite))
}
