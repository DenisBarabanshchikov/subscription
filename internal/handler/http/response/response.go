package response

type CreateCustomer struct {
	CustomerId         string `json:"customerId"`
	ExternalCustomerId string `json:"externalCustomerId"`
}

type SubscribeCustomer struct {
	SubscriptionId         string `json:"subscriptionId"`
	ExternalSubscriptionId string `json:"externalSubscriptionId"`
}

type SubscriptionStatus struct {
	SubscriptionId         string `json:"subscriptionId"`
	ExternalSubscriptionID string `json:"externalSubscriptionId"`
	Plan                   string `json:"plan"`
	Status                 string `json:"status"`
}
