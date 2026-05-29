package server

import (
	"backend/src/controllers"
    authHandler"backend/src/controllers/auth"
	"backend/src/middlewares"
	orderService"backend/src/services/order"

	"github.com/gin-gonic/gin"
)

func registerRoutes(r *gin.Engine, sm *ServerManager, lp *middlewares.LoginPrison) {
	userCtrl := &authHandler.UserManager{
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
        Users: sm.Users,
        OrderService: &orderService.OrderService{Tx: sm.Tx},
    }

    paymentCtrl := &controllers.PaymentManager{
        Orders: sm.Orders,
        Products: sm.Products,
        Store: sm.Stores,
        Payment: sm.Payment,
    }

    sudoCtrl := &controllers.SudoManager{
        SUDOSecret: string(sm.SUDOSecret),
    }

    cartCtrl := &controllers.CartManager{
        Carts: sm.Carts,
        Products: sm.Products,
    }

    auth := r.Group("/v2/auth")
    {
        auth.GET("/users",     userCtrl.LoginUser(lp))
        auth.POST("/users",    userCtrl.RegisterUser())
        auth.POST("/refresh", userCtrl.Refresh())
        auth.POST("/logout",  userCtrl.LogoutUser())
    }

    user := r.Group("/v2/user")
    user.Use(middlewares.AuthMiddlewares(string(sm.JWTSecret)))
    {
        user.POST("/store", storeCtrl.RegisterStore())
    }

    user.Use(middlewares.SudoMiddleware(string(sm.JWTSecret)))
    {
        user.PUT("", userCtrl.UpdateUser())
        user.DELETE("", userCtrl.DeleteUser())
    }

    store := r.Group("/v1/store")
    store.Use(middlewares.AuthMiddlewares(string(sm.JWTSecret)))
    store.Use(middlewares.StoreMiddleware(sm.Stores))
    store.Use(middlewares.SudoMiddleware(string(sm.JWTSecret)))
    {
        store.PUT("",    storeCtrl.UpdateStore())
        store.DELETE("", storeCtrl.DeleteStore())
    }

    product := r.Group("/v1/store")
    product.Use(middlewares.AuthMiddlewares(string(sm.JWTSecret)))
    product.Use(middlewares.StoreMiddleware(sm.Stores))
    {
        product.POST("/products", prdctCtrl.RegisterProduct())
    }

    product.Use(middlewares.SudoMiddleware(string(sm.JWTSecret)))
    {
        product.PUT("/products", prdctCtrl.UpdateProduct())
        product.DELETE("/products", prdctCtrl.RemoveProduct())
    }

    order := r.Group("/v1/order")
    order.Use(middlewares.AuthMiddlewares(string(sm.JWTSecret)))
    {
        order.POST("", orderCtrl.CreateOrder())
    }

    order.Use(middlewares.SudoMiddleware(string(sm.JWTSecret)))
    {
        order.DELETE("", orderCtrl.CancelOrder())
    }

    cart := r.Group("/v1/cart")
    cart.Use(middlewares.AuthMiddlewares(string(sm.JWTSecret)))
    {
        cart.POST("", cartCtrl.AddCartItem())
        cart.PUT("", cartCtrl.UpdateQuantity())
        cart.GET("", cartCtrl.GetCarts())
        cart.DELETE("", cartCtrl.RemoveCartItem())
        cart.DELETE("/clear", cartCtrl.ClearCart())
    }

    payment := r.Group("/v1/payment")
    payment.Use(middlewares.AuthMiddlewares(string(sm.JWTSecret)))
    {
        payment.POST("/stripe", paymentCtrl.CheckoutOrder())
    }

    webhook:= r.Group("/v1/webhook")
    {
        webhook.POST("/stripe", paymentCtrl.HandleWebhooks())
    }

    sudo := r.Group("v1/confirmation")
    sudo.Use(middlewares.AuthMiddlewares(string(sm.JWTSecret)))
    {
        sudo.POST("", sudoCtrl.GenerateConfirmation())
    }
}