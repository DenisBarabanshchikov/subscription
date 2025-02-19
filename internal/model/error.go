package model

import "fmt"

type CustomerNotFoundErr struct {
	msg string
}

func NewCustomerNotFoundErr(customerId string) CustomerNotFoundErr {
	return CustomerNotFoundErr{msg: fmt.Sprintf("customer '%s' not found", customerId)}
}

func (e CustomerNotFoundErr) Error() string {
	return e.msg
}

type SubscriptionNotFoundErr struct {
	msg string
}

func NewSubscriptionNotFoundErr(subscriptionId string) SubscriptionNotFoundErr {
	return SubscriptionNotFoundErr{msg: fmt.Sprintf("subscription '%s' not found", subscriptionId)}
}

func (e SubscriptionNotFoundErr) Error() string {
	return e.msg
}
