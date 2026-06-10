package server

import (
	"backend/src/controllers"
	authHandler "backend/src/controllers/user"
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
        Cache: sm.Cacher,
    }

    prdctCtrl := &controllers.ProductManager{
        Stores: sm.Stores,
        Products: sm.Products,
        JWTSecret: sm.JWTSecret,
        Cache: sm.Cacher,
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

    public := r.Group("/v2")
    {
        auth := public.Group("/auth")
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

        api := r.Group("/api")
        {
            api.GET("/users", userCtrl.SearchUser())
            api.GET("/products", prdctCtrl.SearchProduct())
            api.GET("/stores", storeCtrl.SearchStore())
        }

        webhook:= public.Group("/webhook")
        {
            webhook.POST("/stripe", paymentCtrl.StripeWebhooks())
        }

        sudo := public.Group("/confirmation")
        sudo.Use(middlewares.AuthMiddlewares(string(sm.JWTSecret)))
        {
            sudo.POST("", sudoCtrl.GenerateConfirmation())
        }
    }

    protected := r.Group("/v2/protected/")
    protected.Use(middlewares.AuthMiddlewares(string(sm.JWTSecret)))
    {
        user := protected.Group("/user")
        {
            user.POST("/register-store", storeCtrl.RegisterStore())

            sudoUser := user.Group("/sudo")
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
        }

        store := r.Group("/store")
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

            product := store.Group("/product")
            product.Use(middlewares.StoreMiddleware(sm.Stores))
            {
                product.POST("", 
                    middlewares.AuthorizationMiddleware(services.ActionProductCreate), 
                    prdctCtrl.RegisterProduct(),
                )

                sudoProduct := product.Group("/sudo")
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
        }
        
        order := protected.Group("/order")
        order.Use(middlewares.AuthMiddlewares(string(sm.JWTSecret)))
        {
            order.POST("", 
                middlewares.AuthorizationMiddleware(services.ActionOrderCreate), 
                orderCtrl.CreateOrder(),
            )

            sudoOrder := order.Group("/sudo")
            sudoOrder.Use(middlewares.SudoMiddleware(string(sm.JWTSecret)))
            {
                order.DELETE("", 
                    middlewares.AuthorizationMiddleware(services.ActionOrderCancel), 
                    orderCtrl.CancelOrder(),
                )
            }
        }

        cart := protected.Group("/cart")
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

        payment := protected.Group("/payment")
        payment.Use(middlewares.AuthMiddlewares(string(sm.JWTSecret)))
        {
            payment.POST("/stripe", paymentCtrl.CheckoutOrder())
        }
    }
}