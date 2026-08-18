package server

import (
	"backend/src/controllers"
	"backend/src/controllers/cart"
	"backend/src/controllers/order"
	"backend/src/controllers/product"
	authHandler "backend/src/controllers/user"
	"backend/src/middlewares"
	"backend/src/services"
	orderService "backend/src/services/order"
	paymentService "backend/src/services/payment"

	"github.com/gin-gonic/gin"
)

func registerRoutes(r *gin.Engine, sm *ServerManager, lp *middlewares.LoginPrison) {
	userCtrl := &authHandler.UserManager{
        Cart:      sm.Carts,
        Users:     sm.Users,
        Tokens:    sm.Tokens,
        JWTSecret: sm.JwtSecret,
    }

    prdctCtrl := &product.ProductManager{
        Products: sm.Products,
        JWTSecret: sm.JwtSecret,
        Cache: sm.Cacher,
    }

    orderCtrl := &order.OrderManager{
        Orders: sm.Orders,
        Users: sm.Users,
        OrderService: &orderService.OrderService{Tx: sm.Tx},
    }

    paymentCtrl := &paymentService.PaymentManager{
        Orders: sm.Orders,
        Products: sm.Products,
        Payment: sm.Payment,
    }

    sudoCtrl := &controllers.SudoManager{
        SUDOSecret: string(sm.SudoSecret),
    }

    cartCtrl := &cart.CartManager{
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
        }

        public.GET("/user/verify", sm.Email.VerifyUserVerification())    

        api := public.Group("/api")
        {
            api.GET("/users",       userCtrl.SearchUser())
            api.GET("/products",    prdctCtrl.SearchProduct())
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
            user.POST("/verify/request", emailCtrl.SendUserVerificationMail())
            user.POST("/register-store", storeCtrl.RegisterStore())
        }

        sudoUser := protected.Group("/sudo/user")
        sudoUser.Use(middlewares.SudoMiddleware(string(sm.SUDOSecret)))
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

        store := protected.Group("/store")
        store.Use(middlewares.StoreMiddleware(sm.Stores))
        {
            store.PUT("", 
                middlewares.AuthorizationMiddleware(services.ActionStoreUpdate), 
                storeCtrl.UpdateStore(),
            )

            product := store.Group("/product")
            product.Use(middlewares.StoreMiddleware(sm.Stores))
            {
                product.POST("", 
                    middlewares.AuthorizationMiddleware(services.ActionProductCreate), 
                    prdctCtrl.RegisterProduct(),
                )
            }

            sudoProduct := store.Group("/sudo/product")
            sudoProduct.Use(middlewares.SudoMiddleware(string(sm.SUDOSecret)))
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

        sudoStore := protected.Group("/sudo/store")
        sudoStore.Use(middlewares.SudoMiddleware(string(sm.SUDOSecret)))
        {
            sudoStore.DELETE("", 
                middlewares.AuthorizationMiddleware(services.ActionStoreDelete), 
                storeCtrl.DeleteStore(),                
            )            
        }
        
        order := protected.Group("/order")
        order.Use(middlewares.AuthMiddlewares(string(sm.JWTSecret)))
        {
            order.POST("", 
                middlewares.AuthorizationMiddleware(services.ActionOrderCreate), 
                orderCtrl.CreateOrder(),
            )
            
            order.POST("/verify/requests", 
                sm.Email.SendOrderConfirmation(),
            )

            order.POST("/verify", 
                sm.Email.VerifyOrderConfirmation(),
            )
        }

        sudoOrder := protected.Group("/sudo/order")
        sudoOrder.Use(middlewares.SudoMiddleware(string(sm.JWTSecret)))
        {
            order.DELETE("", 
                middlewares.AuthorizationMiddleware(services.ActionOrderCancel), 
                orderCtrl.CancelOrder(),
            )
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
                middlewares.AuthorizationMiddleware(services.ActionCartRead), 
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