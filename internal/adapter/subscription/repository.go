package subscription

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/pkg/errors"
	"time"
)

type Repository interface {
	CreateCustomer(ctx context.Context, entity Customer) error
	GetCustomer(ctx context.Context, customerId string) (*Customer, error)
	CreateSubscription(ctx context.Context, entity Subscription) error
	GetSubscription(ctx context.Context, customerId, subscriptionId string) (*Subscription, error)
}

type DynamoConfig struct {
	Client       *dynamodb.Client
	Table        string
	QueryTimeout time.Duration
}

type dynamoRepository struct {
	client       *dynamodb.Client
	table        string
	queryTimeout time.Duration
}

func NewDynamoRepository(config DynamoConfig) Repository {
	return &dynamoRepository{
		client:       config.Client,
		table:        config.Table,
		queryTimeout: config.QueryTimeout,
	}
}

func (d *dynamoRepository) CreateCustomer(ctx context.Context, entity Customer) error {
	ctx, cancel := context.WithTimeout(context.WithoutCancel(ctx), d.queryTimeout)
	defer cancel()

	atr, err := attributevalue.MarshalMap(&entity)
	if err != nil {
		return errors.Wrapf(err, "failed to marshal dynamo customer entity")
	}

	pk := fmt.Sprintf("CUSTOMER#%s", entity.CustomerId)
	sk := fmt.Sprintf("CUSTOMER#%s", entity.CustomerId)
	atr["PK"] = &types.AttributeValueMemberS{Value: pk}
	atr["SK"] = &types.AttributeValueMemberS{Value: sk}

	input := &dynamodb.PutItemInput{
		Item:                atr,
		TableName:           aws.String(d.table),
		ConditionExpression: aws.String("attribute_not_exists(PK)"),
	}

	_, err = d.client.PutItem(ctx, input)
	if err != nil {
		return errors.Wrapf(err, "failed to put dynamo customer entity")
	}

	return nil
}

func (d *dynamoRepository) GetCustomer(ctx context.Context, customerId string) (*Customer, error) {
	ctx, cancel := context.WithTimeout(context.WithoutCancel(ctx), d.queryTimeout)
	defer cancel()

	pk := fmt.Sprintf("CUSTOMER#%s", customerId)
	sk := fmt.Sprintf("CUSTOMER#%s", customerId)

	input := &dynamodb.GetItemInput{
		Key: map[string]types.AttributeValue{
			"PK": &types.AttributeValueMemberS{Value: pk},
			"SK": &types.AttributeValueMemberS{Value: sk},
		},
		TableName: aws.String(d.table),
	}

	result, err := d.client.GetItem(ctx, input)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get dynamo customer entity")
	}

	return unmarshalCustomerEntity(result)
}

func (d *dynamoRepository) CreateSubscription(ctx context.Context, entity Subscription) error {
	ctx, cancel := context.WithTimeout(context.WithoutCancel(ctx), d.queryTimeout)
	defer cancel()

	atr, err := attributevalue.MarshalMap(&entity)
	if err != nil {
		return errors.Wrapf(err, "failed to marshal dynamo subscription entity")
	}

	pk := fmt.Sprintf("CUSTOMER#%s", entity.CustomerId)
	sk := fmt.Sprintf("SUBSCRIPTION#%s", entity.SubscriptionId)
	atr["PK"] = &types.AttributeValueMemberS{Value: pk}
	atr["SK"] = &types.AttributeValueMemberS{Value: sk}

	input := &dynamodb.PutItemInput{
		Item:                atr,
		TableName:           aws.String(d.table),
		ConditionExpression: aws.String("attribute_not_exists(PK)"),
	}

	_, err = d.client.PutItem(ctx, input)
	if err != nil {
		return errors.Wrapf(err, "failed to put dynamo subscription entity")
	}

	return nil
}

func (d *dynamoRepository) GetSubscription(ctx context.Context, customerId, subscriptionId string) (*Subscription, error) {
	ctx, cancel := context.WithTimeout(context.WithoutCancel(ctx), d.queryTimeout)
	defer cancel()

	pk := fmt.Sprintf("CUSTOMER#%s", customerId)
	sk := fmt.Sprintf("SUBSCRIPTION#%s", subscriptionId)

	input := &dynamodb.GetItemInput{
		Key: map[string]types.AttributeValue{
			"PK": &types.AttributeValueMemberS{Value: pk},
			"SK": &types.AttributeValueMemberS{Value: sk},
		},
		TableName: aws.String(d.table),
	}

	result, err := d.client.GetItem(ctx, input)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get dynamo subscription entity")
	}

	return unmarshalSubscriptionEntity(result)
}

func unmarshalCustomerEntity(result *dynamodb.GetItemOutput) (*Customer, error) {
	if result.Item == nil {
		return nil, nil
	}
	var entity Customer
	if err := attributevalue.UnmarshalMap(result.Item, &entity); err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal dynamo customer entity")
	}
	return &entity, nil
}

func unmarshalSubscriptionEntity(result *dynamodb.GetItemOutput) (*Subscription, error) {
	if result.Item == nil {
		return nil, nil
	}
	var entity Subscription
	if err := attributevalue.UnmarshalMap(result.Item, &entity); err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal dynamo subscription entity")
	}
	return &entity, nil
}
