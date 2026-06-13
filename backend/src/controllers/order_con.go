package controllers

import (
	"backend/src/middlewares"
	"backend/src/models"
	repo "backend/src/repository"
	orderSrvc "backend/src/services/order"
	"backend/src/utils"
	Logger "backend/src/utils/logger"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

/* ORDER DOCUMENTATION
* 1. CreateOrder
*	Create order is the process that handle the order creation, it take an array of items, and optional location 
*	(optional because we take the default user location from database).
*	
*	The email and id of the buyer retrieved from header that setted by middleware. After that we prep for the Place Order
*	service, the reason i use separate parameter for the place order is to avoid tangled import sequence since that not allowed
*	in golang. The function PlaceOrder return a pointer to the order data or model whatever it called.
*
* 	After the create success it will return a json formatted data, and redirect to stripe checkout payment page...
*	I think i don't suppose to do that in the server side, isn't it supposed to be client side?
*
* 2. CancelOrder
*	Cancel order is a function that it sole purpose is handling the cancel mechanism of order, it required order id and user id to work,
*	user id will be retrieve from jwt token while the oid (order id) retrieved from the client sided in json formatted data.
*
*	it will then be binded to json, And call CancelOrder function from service directory, simply CancelOrder is just service layer
*	that wrap the Order repository with a centralized transaction function (Transaction is a sql function that make sure if and error
	happened in the middle of the process everything reset to it's default state same thing like c.Abort() from gin).
*	
*	If it's CancelOrder service success it return nil.
*/

type OrderManager struct {
	Orders 		repo.OrderStoreInterface
	Users 		repo.UserStoreInterface
	OrderService 	*orderSrvc.OrderService
}

//
//	ORDER REQUEST SCHEMATIC
//

type orderResponse struct {
	ID string
	BuyerID string
	BuyerEmail string
	Total float64
	Location string
	Status string
	CreatedAt time.Time
}

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

func (om *OrderManager) CreateOrder() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req OrderRequest

		if err := c.ShouldBindBodyWithJSON(&req); err != nil {
			Logger.Log.Error("detail", zap.Error(err))
			c.Error(middlewares.ErrBadRequest("Failed to read client request"))
			return
		}

		buyerId := c.GetString("userId")

		user, err := om.Users.FetchUserByID(c.Request.Context(), &buyerId)

		if err != nil {
    		c.Error(middlewares.ErrInternal("Failed to fetch user"))
    		return
		}

		buyerEmail := user.Email

		var location string;

        if req.BuyerLocation == nil {
            location = fmt.Sprintf("%s, %s", user.Country, user.Address)
		}

		items := make([]repo.OrderItemInput, 0)
		for _, item := range req.Items {
			items = append(items, repo.OrderItemInput{
				ProductID: &item.ProductID,
				Quantity: &item.Quantity,
			})
		}
		
		order, err := om.OrderService.PlaceOrder(c.Request.Context(), &orderSrvc.PlaceOrderParams{
			BuyerID: buyerId,
            BuyerEmail: buyerEmail,
            CombinedLocation: location,
            Status: "pending",
            ProductList: items,
		})

		if err != nil {
            Logger.Log.Error("Error in creating order", zap.Error(err))
            c.Error(middlewares.ErrInternal("Failed to place order"))
            return
		}

		c.JSON(http.StatusCreated, gin.H{
			"response": utils.EXIT_SUCCESS,
			"detail": gin.H{
				"info": "order created",
				"order": toOrderResponse(order),
			},
		})
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

		err := om.OrderService.CancelOrder(c.Request.Context(), &orderSrvc.CancelOrderParams{
			BuyerId: buyerId,
			OrderId: req.OID,
		})

		if err != nil {
			Logger.Log.Error("detail", zap.Error(err))
			c.Error(middlewares.ErrInternal("Failed to remove order"))
			return
		}

		c.JSON(http.StatusOK, gin.H{"response": "order removed"})
	}
}