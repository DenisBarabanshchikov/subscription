package main

import (
	"github.com/DenisBarabanshchikov/subscription/config"
	"github.com/DenisBarabanshchikov/subscription/di"
	_ "github.com/DenisBarabanshchikov/subscription/docs"
	"github.com/gin-gonic/gin"
	"log"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title       Subscription Service API Documentation
// @version     1.0.0
// @description This is the API documentation for the subscription service.
func main() {
	router := gin.Default()

	// Swagger documentation
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Initialize handlers with dependency injection
	h, err := di.InitializeHandlers()
	if err != nil {
		log.Fatalf("failed to initialize handlers: %v", err)
	}

	// Routes
	api := router.Group("/api/v1")
	{
		api.POST("/customers", h.SubscriptionHandler.CreateCustomer)

		// 2) Create a new subscription for a given customer
		api.POST("/customers/:customerId/subscriptions", h.SubscriptionHandler.SubscribeCustomer)

		// 3) Retrieve a subscription’s status (or details) for a given customer
		api.GET("/customers/:customerId/subscriptions/:subscriptionId", h.SubscriptionHandler.GetSubscriptionStatus)

		// 4) Handle Stripe webhook events
		api.POST("/stripe/webhook", h.SubscriptionHandler.HandleStripeWebhook)
	}

	// Run server
	if err := router.Run(config.ServerAddress); err != nil {
		log.Fatalf("failed to run server: %v", err)
	}
}
