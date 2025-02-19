package subscription

import "time"

type Customer struct {
	CustomerId         string    `dynamodbav:"CustomerId"`
	ExternalCustomerId string    `dynamodbav:"ExternalCustomerId"`
	CreatedAt          time.Time `dynamodbav:"CreatedAt"`
	UpdatedAt          time.Time `dynamodbav:"UpdatedAt"`
}

type Subscription struct {
	SubscriptionId         string    `dynamodbav:"SubscriptionId"`
	CustomerId             string    `dynamodbav:"CustomerId"`
	ExternalSubscriptionID string    `dynamodbav:"ExternalSubscriptionId"`
	Plan                   string    `dynamodbav:"Plan"`
	Status                 string    `dynamodbav:"Status"`
	CreatedAt              time.Time `dynamodbav:"CreatedAt"`
	UpdatedAt              time.Time `dynamodbav:"UpdatedAt"`
}
