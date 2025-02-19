package http

import (
	"github.com/DenisBarabanshchikov/subscription/internal/handler/http/request"
	"github.com/DenisBarabanshchikov/subscription/internal/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

type SubscriptionHandler struct {
	subscriptionService service.SubscriptionService
}

func NewSubscriptionHandler(subscriptionService service.SubscriptionService) *SubscriptionHandler {
	return &SubscriptionHandler{
		subscriptionService: subscriptionService,
	}
}

// CreateCustomer handles the customer create request.
// @Description  Creating a new customer
// @Tags         Customer
// @Accept       application/json
// @Produce      json
// @Param        request  body  request.CreateCustomer  true  "Customer data"
// @Success      200  {object}  response.CreateCustomer
// @Failure      400  {object}  response.ErrorResponse
// @Failure      500  {object}  response.ErrorResponse
// @Router       /api/v1/customers [post]
func (h *SubscriptionHandler) CreateCustomer(c *gin.Context) {
	var req request.CreateCustomer
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx := c.Request.Context()

	customer, err := h.subscriptionService.CreateCustomer(ctx, req.Email)
	if err != nil {
		res := handleError(ctx, err)
		c.JSON(res.Code, res)
		return
	}

	c.JSON(http.StatusOK, mapToCreateCustomerResponse(customer))
}

// SubscribeCustomer handles the subscribe customer request.
// @Description  Subscribe a customer (Available plans: Core, Growth, Premium)
// @Tags         Customer
// @Accept       application/json
// @Produce      json
// @Param        customerId    path      string  true  "customerId"
// @Param        request  body  request.SubscribeCustomer  true  "Subscription data"
// @Success      200  {object}  response.SubscribeCustomer
// @Failure      400  {object}  response.ErrorResponse
// @Failure      500  {object}  response.ErrorResponse
// @Router       /api/v1/customers/{customerId}/subscriptions [post]
func (h *SubscriptionHandler) SubscribeCustomer(c *gin.Context) {
	var req request.SubscribeCustomer
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	customerId := c.Param("customerId")

	ctx := c.Request.Context()

	subscription, err := h.subscriptionService.SubscriberCustomer(ctx, customerId, req.Plan)
	if err != nil {
		res := handleError(ctx, err)
		c.JSON(res.Code, res)
		return
	}

	c.JSON(http.StatusOK, mapToSubscriberCustomerResponse(subscription))
}

// GetSubscriptionStatus handles the get subscription status request.
// @Description  Get subscription status
// @Tags         Customer
// @Accept       application/json
// @Produce      json
// @Param        customerId    path      string  true  "customerId"
// @Param        subscriptionId    path      string  true  "subscriptionId"
// @Success      200  {object}  response.SubscriptionStatus
// @Failure      400  {object}  response.ErrorResponse
// @Failure      500  {object}  response.ErrorResponse
// @Router       /api/v1/customers/{customerId}/subscriptions/{subscriptionId} [get]
func (h *SubscriptionHandler) GetSubscriptionStatus(c *gin.Context) {
	customerId := c.Param("customerId")
	subscriptionId := c.Param("subscriptionId")

	ctx := c.Request.Context()

	subscription, err := h.subscriptionService.SubscriptionStatus(ctx, customerId, subscriptionId)
	if err != nil {
		res := handleError(ctx, err)
		c.JSON(res.Code, res)
		return
	}

	c.JSON(http.StatusOK, mapToSubscriptionStatusResponse(subscription))
}

// HandleStripeWebhook handles the stripe webhook.
// @Description  Handles the stripe webhook
// @Tags         Stripe
// @Accept       application/json
// @Produce      json
// @Success      202  "Accepted - no content"
// @Failure      400  {object}  response.ErrorResponse
// @Failure      500  {object}  response.ErrorResponse
// @Router       /api/v1/stripe/webhook [post]
func (h *SubscriptionHandler) HandleStripeWebhook(c *gin.Context) {
	//TODO handle webhooks
	//Suggestion create separate notification service,
	// which will add such events in queue, subscription service will read from this queue
	c.JSON(http.StatusAccepted, nil)
}
