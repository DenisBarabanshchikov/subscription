package dynamo_client

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

func NewClient(c Config) (*dynamodb.Client, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(c.Region),
	)
	if err != nil {
		return nil, fmt.Errorf("unable to load aws config for dynamodb client: %s", err)
	}

	return dynamodb.NewFromConfig(cfg, func(options *dynamodb.Options) {
		options.BaseEndpoint = c.EndpointUrl
	}), nil
}
