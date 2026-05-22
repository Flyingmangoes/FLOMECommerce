package server

import (
	"backend/src/controllers"
	"backend/src/middlewares"
	"backend/src/services"

	"github.com/gin-gonic/gin"
)

func registerRoutes(r *gin.Engine, sm *ServerManager, lp *middlewares.LoginPrison) {
	userCtrl := &controllers.UserContext{
        Users:     sm.Users,
        Tokens:    sm.Tokens,
        JWTSecret: sm.JWTSecret,
    }

    storeCtrl := &controllers.StoreManager{
		Users: sm.Users,
        Stores: sm.Stores,
        JWTSecret: sm.JWTSecret,
    }

    prdctCtrl := &controllers.ProductManager{
        Stores: sm.Stores,
        Products: sm.Products,
        JWTSecret: sm.JWTSecret,
    }

    orderCtrl := &controllers.OrderManager{
        Orders: sm.Orders,
        Products: sm.Products,
        Users: sm.Users,
        OrderService: &services.OrderService{Tx: sm.Tx},
        Payment: sm.Payment,
    }

    // public auth routes
    auth := r.Group("/v2/auth")
    {
        auth.GET("/users",     userCtrl.LoginUser(lp))
        auth.POST("/users",    userCtrl.RegisterUser())
        auth.POST("/refresh", userCtrl.Refresh())
        auth.POST("/logout",  userCtrl.LogoutUser())
    }

    // protected user routes
    user := r.Group("/v2/user")
    user.Use(middlewares.AuthMiddlewares(string(sm.JWTSecret)))
    {
        user.POST("/store", storeCtrl.RegisterStore())
        user.Use(middlewares.ConfirmationMiddleware(string(sm.JWTSecret)))
        user.PUT("", userCtrl.UpdateUser())
        user.DELETE("", userCtrl.DeleteUser())
    }

    store := r.Group("/v1/store")
    store.Use(middlewares.AuthMiddlewares(string(sm.JWTSecret)))
    store.Use(middlewares.CheckForStore(sm.Stores))
    store.Use(middlewares.ConfirmationMiddleware(string(sm.JWTSecret)))
    {
        store.PUT("",    storeCtrl.UpdateStore())
        store.DELETE("", storeCtrl.DeleteStore())
    }

    product := r.Group("/v1/store")
    product.Use(middlewares.AuthMiddlewares(string(sm.JWTSecret)))
    product.Use(middlewares.CheckForStore(sm.Stores))
    {
        product.POST("/products", prdctCtrl.RegisterProduct())
        product.Use(middlewares.ConfirmationMiddleware(string(sm.JWTSecret)))
        product.PUT("/products", prdctCtrl.UpdateProduct())
        product.DELETE("/products", prdctCtrl.RemoveProduct())
    }

    order := r.Group("/v1/order")
    order.Use(middlewares.AuthMiddlewares(string(sm.JWTSecret)))
    {
        order.POST("", orderCtrl.CreateOrder())
        order.POST("/checkout", orderCtrl.CheckoutOrder())
        order.Use(middlewares.ConfirmationMiddleware(string(sm.JWTSecret)))
        order.DELETE("", orderCtrl.CancelOrder())
    }

}