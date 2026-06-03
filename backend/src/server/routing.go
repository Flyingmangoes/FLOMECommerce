package server

import (
	"backend/src/controllers"
	authHandler "backend/src/controllers/auth"
	"backend/src/middlewares"
	"backend/src/services"
	orderService "backend/src/services/order"
	paymentService "backend/src/services/payment"

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

    paymentCtrl := &paymentService.PaymentManager{
        Orders: sm.Orders,
        Products: sm.Products,
        Store: sm.Stores,
        Payment: sm.Payment,
    }

    emailCtrl := sm.Email

    sudoCtrl := &controllers.SudoManager{
        SUDOSecret: string(sm.SUDOSecret),
    }

    cartCtrl := &controllers.CartManager{
        Carts: sm.Carts,
        Products: sm.Products,
    }

    auth := r.Group("/v2/auth")
    {
        auth.GET("/users",      userCtrl.LoginUser(lp))
        auth.POST("/users",     userCtrl.RegisterUser())
        auth.POST("/refresh",   userCtrl.Refresh())
        auth.POST("/logout",    userCtrl.LogoutUser())

        verificationAuth := r.Group("/verify")
        verificationAuth.Use(middlewares.VerificationMiddleware(sm.VERIFICATION_SECRET))
        auth.POST("/request", emailCtrl.SendVerificationMail())
        auth.POST("",         emailCtrl.VerifyEmail())
    }

    user := r.Group("/v2/user")
    user.Use(middlewares.AuthMiddlewares(string(sm.JWTSecret)))
    {
        user.POST("/store", storeCtrl.RegisterStore())
    }

        sudoUser := r.Group("sudo")
        sudoUser.Use(middlewares.SudoMiddleware(string(sm.JWTSecret)))
    {
        sudoUser.PUT("", 
            middlewares.AuthorizationMiddleware(services.ActionProfileUpdate), 
            userCtrl.UpdateUser(),
        )

        sudoUser.DELETE("", 
            middlewares.AuthorizationMiddleware(services.ActionProfileDelete), 
            userCtrl.DeleteUser(),
        )
    }

    store := r.Group("/v1/store")
    store.Use(middlewares.AuthMiddlewares(string(sm.JWTSecret)))
    store.Use(middlewares.StoreMiddleware(sm.Stores))
    store.Use(middlewares.SudoMiddleware(string(sm.JWTSecret)))
    {
        store.PUT("", 
            middlewares.AuthorizationMiddleware(services.ActionStoreUpdate), 
            storeCtrl.UpdateStore(),
        )

        store.DELETE("", 
            middlewares.AuthorizationMiddleware(services.ActionStoreDelete), 
            storeCtrl.DeleteStore(),
        )
    }

    product := r.Group("/v1/store/product")
    product.Use(middlewares.AuthMiddlewares(string(sm.JWTSecret)))
    product.Use(middlewares.StoreMiddleware(sm.Stores))
    {
        product.POST("", 
            middlewares.AuthorizationMiddleware(services.ActionProductCreate), 
            prdctCtrl.RegisterProduct(),
        )

        sudoProduct := r.Group("sudo")
        sudoProduct.Use(middlewares.SudoMiddleware(string(sm.JWTSecret)))
        {
            sudoProduct.PUT("", 
                middlewares.AuthorizationMiddleware(services.ActionProductUpdate), 
                prdctCtrl.UpdateProduct(),
            )
            sudoProduct.DELETE("", 
                middlewares.AuthorizationMiddleware(services.ActionProductDelete), 
                prdctCtrl.RemoveProduct(),
            )
        }
    }

    order := r.Group("/v1/order")
    order.Use(middlewares.AuthMiddlewares(string(sm.JWTSecret)))
    {
        order.POST("", 
            middlewares.AuthorizationMiddleware(services.ActionOrderCreate), 
            orderCtrl.CreateOrder(),
        )

        sudoOrder := r.Group("sudo")
        sudoOrder.Use(middlewares.SudoMiddleware(string(sm.JWTSecret)))
        {
            order.DELETE("", 
                middlewares.AuthorizationMiddleware(services.ActionOrderCancel), 
                orderCtrl.CancelOrder(),
            )
        }
    }

    cart := r.Group("/v1/cart")
    cart.Use(middlewares.AuthMiddlewares(string(sm.JWTSecret)))
    {
        cart.POST("", 
            middlewares.AuthorizationMiddleware(services.ActionCartAdd), 
            cartCtrl.AddCartItem(),
        )

        cart.PUT("", 
            middlewares.AuthorizationMiddleware(services.ActionCartUpdate), 
            cartCtrl.UpdateQuantity(),
        )

        cart.GET("", 
            middlewares.AuthorizationMiddleware(services.ActionCartSelfRead), 
            cartCtrl.GetCarts(),
        )

        cart.DELETE("", 
            middlewares.AuthorizationMiddleware(services.ActionCartRemove), 
            cartCtrl.RemoveCartItem(),
        )

        cart.DELETE("/clear", 
            middlewares.AuthorizationMiddleware(services.ActionCartClear), 
            cartCtrl.ClearCart(),
        )
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