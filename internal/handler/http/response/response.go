package ressponse

type CreateCustomer struct {
	CustomerId       string `json:"customerId"`
	StripeCustomerId string `json:"stripeCustomerId"`
}
