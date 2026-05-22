package controllers

import (
	"backend/src/middlewares"
	"backend/src/models"
	"backend/src/repository"
	"backend/src/services"
	Logger "backend/src/utils/logger"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type OrderManager struct {
	Orders 		repository.OrderStoreInterface
	Store 		repository.StoreStoreInterface
	Users 		repository.UserStoreInterface
	Products 	repository.ProductStoreInterface	 
	OrderService 	*services.OrderService
	Payment			*services.PaymentService
}

//
// Order Structure
//

type OrderItemRequest struct {
    ProductID 	string `json:"productId" binding:"required"`
    Quantity  	int    `json:"quantity"  binding:"required,min=1"`
}

type OrderRequest struct {
	Items []OrderItemRequest `json:"items" binding:"required"`
    BuyerLocation *string       `json:"location" binding:"omitempty"`
}	

type CancelOrderRequest struct {
	OID string `json:"orderId" binding:"required"`
	Confirmation bool `json:"confirmation" binding:"required"`
}

type orderResponse struct {
	ID string
	BuyerID string
	BuyerEmail string
	Total float64
	Location string
	Status string
	CreatedAt time.Time
}

//
// Handler
//

func toOrderResponse(o *models.Order) orderResponse {
	return orderResponse{
		ID: o.ID,
		BuyerID: o.BuyerID,
		BuyerEmail: o.BuyerEmail,
		Total: o.TotalPrice,
		Location: o.Location,
		Status: o.Status,
		CreatedAt: o.CreatedAt,
	}
}

/* ORDER DOCUMENTATION
* 1. CreateOrder
*	create order is a function that used by gin.IRoutes and used OrderRequest struct
*	as the parameter that then be  binded using gin.Context.ShouldBindBodyWithJSON
*
*	for buyer id aka user id we use GetString(), because the user id already set 
*	in the auth middleware which used for retrieve the user data which will be 
*	used to fill an array of OrderItemInput that used in-- 
*	OrderStoreInterface.PlaceOrder that require PlaceOrderParams to work.
*	
* 2. CheckoutOrder
* 3. CancelOrder
*/

func (om *OrderManager) CreateOrder() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req OrderRequest

		if err := c.ShouldBindBodyWithJSON(&req); err != nil {
			Logger.Log.Error("detail", zap.Error(err))
			c.Error(middlewares.ErrBadRequest("Failed to read client request"))
			return
		}

		buyerId := c.GetString("userId")

		user, err := om.Users.GetUserByID(c.Request.Context(), buyerId)

		if err != nil {
    		c.Error(middlewares.ErrInternal("Failed to fetch user"))
    		return
		}

		buyerEmail := user.Email

		var location string;

        if req.BuyerLocation == nil {
            location = fmt.Sprintf("%s, %s", user.Country, user.Address)
		}

		items := make([]repository.OrderItemInput, 0)
		for _, item := range req.Items {
			items = append(items, repository.OrderItemInput{
				ProductID: &item.ProductID,
				Quantity: &item.Quantity,
			})
		}
		
		order, err := om.OrderService.PlaceOrder(c.Request.Context(), &services.PlaceOrderParams{
			BuyerID: buyerId,
            BuyerEmail: buyerEmail,
            CombinedLocation: location,
            Status: "pending",
            ProductList: items,
		})

		if err != nil {
            Logger.Log.Error("detail", zap.Error(err))
            c.Error(middlewares.ErrInternal("Failed to place order"))
            return
		}

		c.JSON(http.StatusCreated, gin.H{
			"response": "order created",
			"detail": gin.H{
				"order": toOrderResponse(order),
			},
		})

		c.Redirect(http.StatusFound, "/checkout")
	}
}

func (om *OrderManager) CheckoutOrder() gin.HandlerFunc {
	return func(c *gin.Context) {
		buyerId := c.GetString("userId")
		items := make([]services.ItemDetail, 0)

		orders, order_items, err := om.Orders.GetOrders(c.Request.Context(), &repository.OrderStoreParams{
			BuyerID: &buyerId,
		})

		if err != nil {
			Logger.Log.Error("Failed retrieve order data",zap.Error(err))
			c.Error(middlewares.ErrInternal("Retrieve order data failed"))
			return
		}

		for i := 0;i < len(order_items); i++{
			product, err := om.Products.GetProductByID(c.Request.Context(), order_items[i].ProductID)
			if err != nil {
				Logger.Log.Error("Failed to retrieve product data", zap.Error(err))
				c.Error(middlewares.ErrInternal("Failed to retrieve product data"))
				return
			}

			store, err := om.Store.GetStoreByID(c.Request.Context(), product.StoreID)
			if err != nil {
				Logger.Log.Error("Failed to retrieve store data", zap.Error(err))
				c.Error(middlewares.ErrInternal("Failed to retrieve store data"))
				return
			}

			item_detail := services.ItemDetail{
				ProductID: order_items[i].ProductID,
				ProductName: product.Name,
				ProductDesc: product.Desc,
				StoreName: store.StoreName,
				Quantity: order_items[i].Quantity,
			}

			items = append(items, item_detail)
		}

		order_details := make([]services.OrderDetail, 0)
		for _, value := range orders {
			od := services.OrderDetail{
				OrderID: value.ID,
				BuyerID: value.BuyerID,
				BuyerEmail: value.BuyerEmail,
				Location: value.Location,
			}

			order_details = append(order_details, od)
		}

		sc, err := om.Payment.CreateCheckoutSession(c.Request.Context(), &items, &order_details)

		if err != nil {
			Logger.Log.Error("Failed to create checkout session", zap.Error(err))
			c.Error(middlewares.ErrInternal("Create checkout failed"))
			return
		}

		Logger.Log.Debug(fmt.Sprintf("Client redirected to: %s", sc.URL))
		c.Redirect(http.StatusFound, sc.URL)
	}
}

func (om *OrderManager) CancelOrder() gin.HandlerFunc {
	return func (c *gin.Context) {
		var req CancelOrderRequest

		if err := c.ShouldBindBodyWithJSON(req); err != nil {
			Logger.Log.Error("detail", zap.Error((err)))
			c.Error(middlewares.ErrBadRequest("Failed to read client request"))
			return
		}

		buyerId := c.GetString("userId")

		err := om.OrderService.CancelOrder(c.Request.Context(), &services.CancelOrderParams{
			BuyerId: buyerId,
			OrderId: req.OID,
			Confirmation: req.Confirmation,
		})

		if err != nil {
			Logger.Log.Error("detail", zap.Error(err))
			c.Error(middlewares.ErrInternal("Failed to remove order"))
			return
		}

		c.JSON(http.StatusOK, gin.H{"response": "order removed"})
	}
}
