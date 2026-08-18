package payment_service

import (
	error_service "backend/src/error"
	repo "backend/src/repository"
	logger_system "backend/src/utils/LoggerSystem"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/stripe/stripe-go/v85"
	"go.uber.org/zap"
)

type PaymentManager struct {
	Products 	repo.ProductStoreInterface
	Orders 		repo.OrderStoreInterface
	Payment 	*PaymentService
}

func (pm *PaymentManager) CheckoutOrder() gin.HandlerFunc {
	return func(c *gin.Context) {
		var (
			buyerEmail string
			orderId string
		)

		buyerId := c.GetString("userId")
		items := make([]ItemDetail, 0)

		orders, order_items, err := pm.Orders.Get(c.Request.Context(), buyerId)
		if err != nil {
			logger_system.Log.Error("Failed retrieve order data",zap.Error(err))
			c.Error(error_service.ErrInternal("Retrieve order data failed"))
			return
		}

		for _, value := range order_items{
			product, err := pm.Products.Get(c.Request.Context(), value.ProductID)
			if err != nil {
				logger_system.Log.Error("Failed to retrieve product data", zap.Error(err))
				c.Error(error_service.ErrInternal("Failed to retrieve product data"))
				return
			}

			item_detail := ItemDetail{
				ProductID: value.ProductID,
				ProductName: product.Name,
				ProductDesc: product.Desc,
				StoreName: "Flomm",
				Price: product.Price,
				Quantity: value.Quantity,
				ImageUrl: product.ProductIMG,
			}

			items = append(items, item_detail)
		}

		for _, value := range orders {
			if value.Status == "pending" {
				buyerEmail = value.BuyerEmail			
				orderId = value.ID
				break
			}
		}

		sc, err := pm.Payment.CreateCheckoutSession(c.Request.Context(), orderId, buyerEmail, items)
		if err != nil {
			logger_system.Log.Error("Failed to create checkout session", zap.Error(err))
			c.Error(error_service.ErrInternal("Create checkout failed"))
			return
		}

		logger_system.Log.Debug(fmt.Sprintf("Client redirected to: %s", sc.URL))
		c.Redirect(http.StatusFound, sc.URL)
	}
}

func(pm *PaymentManager) StripeWebhooks() gin.HandlerFunc {
	return func (c *gin.Context)  {
		payload, err := io.ReadAll(c.Request.Body)
		if err != nil {
			logger_system.Log.Error("Failed to read request body", zap.Error(err))
			c.Status(http.StatusBadRequest)
			return
		}
		
		signHeader := c.GetHeader("Stripe-Signature")
		event, err := pm.Payment.StripeClient.ConstructEvent(
			payload, 
			signHeader, 
			pm.Payment.WebhookKey,
		)
		if err != nil {
			logger_system.Log.Error("Failed construct event", zap.Error(err))
			c.Status(http.StatusBadRequest)
			return
		}

		switch event.Type {
		case "checkout.session.completed":
			var session stripe.CheckoutSession
			if err := json.Unmarshal(event.Data.Raw, &session); err != nil {
				c.Status(http.StatusBadRequest)
				return
			}

			orderId := session.Metadata["order_id"]
			pm.Orders.UpdateStatus(c.Request.Context(), orderId, "paid")
			logger_system.Log.Info("Order paid", zap.String("order_id", orderId))

		case "checkout.session.expired":
			var session stripe.CheckoutSession
			if err := json.Unmarshal(event.Data.Raw, &session); err != nil {
				c.Status(http.StatusBadRequest)
				return
    		}

			orderId := session.Metadata["order_id"]
			pm.Orders.UpdateStatus(c.Request.Context(), orderId, "cancelled")
			logger_system.Log.Info("Order cancelled", zap.String("order_id", orderId))

		default:
			logger_system.Log.Info("Unhandled event", zap.String("type", string(event.Type)))
		}

		c.Status(http.StatusOK)
	}
}