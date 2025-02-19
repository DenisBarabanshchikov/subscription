package config

import (
	"github.com/DenisBarabanshchikov/subscription/internal/adapter/subscription"
	"github.com/DenisBarabanshchikov/subscription/pkg/dynamo_client"
	"github.com/DenisBarabanshchikov/subscription/pkg/env"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"log"
)

func GetDynamoClient() *dynamodb.Client {
	conf := dynamo_client.Config{
		EndpointUrl: env.OptionalStringPtr("DYNAMO_ENDPOINT"),
		Region:      env.RequiredString("AWS_REGION"),
	}

	client, err := dynamo_client.NewClient(conf)
	if err != nil {
		log.Fatal("failed to initialize dynamo client")
	}
	return client
}

func ProvideSubscriptionDynamoConfig() subscription.DynamoConfig {
	return subscription.DynamoConfig{
		Client:       GetDynamoClient(),
		Table:        env.RequiredString("DYNAMO_SUBSCRIPTION_TABLE"),
		QueryTimeout: env.RequiredDuration("DYNAMO_SUBSCRIPTION_TIMEOUT"),
	}
}
