package payment_service

import (
	repo"backend/src/repository"
	"backend/src/middlewares"
	Logger"backend/src/utils/logger"
	"encoding/json"
	"io"
	"net/http"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/stripe/stripe-go/v85"
	"go.uber.org/zap"
)

type PaymentManager struct {
	Products 	repo.ProductStoreInterface
	Orders 		repo.OrderStoreInterface
	Store 		repo.StoreStoreInterface
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

		oid, err := pm.Orders.GetOrderId(c.Request.Context(), buyerId)
		if err != nil {
			Logger.Log.Error("Failed retrieve order id", zap.Error(err))
			c.Error(middlewares.ErrInternal("Retrieve order id failed"))
			return
		}

		orders, order_items, err := pm.Orders.GetOrders(c.Request.Context(), buyerId, oid)

		if err != nil {
			Logger.Log.Error("Failed retrieve order data",zap.Error(err))
			c.Error(middlewares.ErrInternal("Retrieve order data failed"))
			return
		}

		for _, value := range order_items{
			product, err := pm.Products.GetProductByID(c.Request.Context(), value.ProductID)
			if err != nil {
				Logger.Log.Error("Failed to retrieve product data", zap.Error(err))
				c.Error(middlewares.ErrInternal("Failed to retrieve product data"))
				return
			}

			store, err := pm.Store.GetStoreByID(c.Request.Context(), product.StoreID)
			if err != nil {
				Logger.Log.Error("Failed to retrieve store data", zap.Error(err))
				c.Error(middlewares.ErrInternal("Failed to retrieve store data"))
				return
			}

			item_detail := ItemDetail{
				ProductID: value.ProductID,
				ProductName: product.Name,
				ProductDesc: product.Desc,
				StoreName: store.StoreName,
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
			Logger.Log.Error("Failed to create checkout session", zap.Error(err))
			c.Error(middlewares.ErrInternal("Create checkout failed"))
			return
		}

		Logger.Log.Debug(fmt.Sprintf("Client redirected to: %s", sc.URL))
		c.Redirect(http.StatusFound, sc.URL)
	}
}

func(pm *PaymentManager) HandleWebhooks() gin.HandlerFunc {
	return func (c *gin.Context)  {
		payload, err := io.ReadAll(c.Request.Body)
		if err != nil {
			Logger.Log.Error("Failed to read request body", zap.Error(err))
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
			Logger.Log.Error("Failed construct event", zap.Error(err))
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
			pm.Orders.UpdateOrderStatus(c.Request.Context(), orderId, "paid")
			Logger.Log.Info("Order paid", zap.String("order_id", orderId))

		case "checkout.session.expired":
			var session stripe.CheckoutSession
			if err := json.Unmarshal(event.Data.Raw, &session); err != nil {
				c.Status(http.StatusBadRequest)
				return
    		}

			orderId := session.Metadata["order_id"]
			pm.Orders.UpdateOrderStatus(c.Request.Context(), orderId, "cancelled")
			Logger.Log.Info("Order cancelled", zap.String("order_id", orderId))

		default:
			Logger.Log.Info("Unhandled event", zap.String("type", string(event.Type)))
		}

		c.Status(http.StatusOK)
	}
}