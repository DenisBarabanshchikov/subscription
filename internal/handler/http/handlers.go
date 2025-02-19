package http

type Handlers struct {
	SubscriptionHandler *SubscriptionHandler
}

func NewHandlers(
	subscriptionHandler *SubscriptionHandler,
) *Handlers {
	return &Handlers{
		SubscriptionHandler: subscriptionHandler,
	}
}
