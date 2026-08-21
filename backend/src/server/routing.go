package server

import (
	"backend/src/controllers/cart"
	sudo "backend/src/controllers/sudo"
	"backend/src/controllers/order"
	"backend/src/controllers/product"
	authHandler "backend/src/controllers/user"
	"backend/src/middlewares"
	auth_service "backend/src/services/auth"
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

    sudoCtrl := &sudo.SudoManager{
        SudoSecret: string(sm.SudoSecret),
    }

    cartCtrl := &cart.CartManager{
        Carts: sm.Carts,
        Products: sm.Products,
    }

    apiPrefix := r.Group("/api") 
    {
        public := apiPrefix.Group("/v2")
        {
            auth := public.Group("/auth")
            {
                auth.GET("/users",      userCtrl.LoginUser(lp))
                auth.POST("/users",     userCtrl.RegisterUser())
                auth.POST("/refresh",   userCtrl.Refresh())
                auth.POST("/logout",    userCtrl.LogoutUser())
            }

            public.GET("/user/verify", sm.Email.VerifyUserVerification())    

            searchApi := public.Group("")
            {  
                searchApi.GET("/products", prdctCtrl.SearchProduct())
            }

            webhook:= public.Group("/webhook")
            {
                webhook.POST("/stripe", paymentCtrl.StripeWebhooks())
            }
        }

        protected := apiPrefix.Group("/v2/protected/")
        protected.Use(middlewares.AuthenticationMiddlewares(string(sm.JwtSecret)))
        {
            sudo := protected.Group("/confirmation")
            {
                sudo.POST("", sudoCtrl.CreateSudo())
            }
            
            user := protected.Group("/user")
            {
                user.POST("/verify/request", sm.Email.UserVerification())
            }

            sudoUser := protected.Group("/sudo/user")
            sudoUser.Use(middlewares.SudoMiddleware(string(sm.SudoSecret)))
            {
                sudoUser.PUT("", 
                    middlewares.AuthorizationMiddleware(auth_service.ActionProfileUpdate), 
                    userCtrl.UpdateUser(),
                )

                sudoUser.DELETE("", 
                    middlewares.AuthorizationMiddleware(auth_service.ActionProfileDelete), 
                    userCtrl.DeleteUser(),
                )
            }

            product := protected.Group("/product")
            {
                product.POST("", 
                    middlewares.AuthorizationMiddleware(auth_service.ActionProductCreate), 
                    prdctCtrl.RegisterProduct(),
                )
            }

            sudoProduct := protected.Group("/sudo/product")
            sudoProduct.Use(middlewares.SudoMiddleware(string(sm.SudoSecret)))
            {
                sudoProduct.PUT("", 
                    middlewares.AuthorizationMiddleware(auth_service.ActionProductUpdate), 
                    prdctCtrl.UpdateProduct(),                    
                )

                sudoProduct.DELETE("", 
                    middlewares.AuthorizationMiddleware(auth_service.ActionProductDelete), 
                    prdctCtrl.RemoveProduct(),
                )
            }
            
            order := protected.Group("/order")
            {
                order.POST("", 
                    middlewares.AuthorizationMiddleware(auth_service.ActionOrderCreate), 
                    orderCtrl.CreateOrder(),
                )
                
                order.POST("/verify/requests", 
                    sm.Email.SendOrderConfirmation(),
                )
            }

            sudoOrder := protected.Group("/sudo/order")
            sudoOrder.Use(middlewares.SudoMiddleware(string(sm.JwtSecret)))
            {
                order.DELETE("", 
                    middlewares.AuthorizationMiddleware(auth_service.ActionOrderCancel), 
                    orderCtrl.CancelOrder(),
                )
            }
                

            cart := protected.Group("/cart")
            {
                cart.POST("", 
                    middlewares.AuthorizationMiddleware(auth_service.ActionCartAdd), 
                    cartCtrl.AddCartItem(),
                )

                cart.PUT("", 
                    middlewares.AuthorizationMiddleware(auth_service.ActionCartUpdate), 
                    cartCtrl.UpdateQuantity(),
                )

                cart.GET("", 
                    middlewares.AuthorizationMiddleware(auth_service.ActionCartRead), 
                    cartCtrl.GetCarts(),
                )

                cart.DELETE("", 
                    middlewares.AuthorizationMiddleware(auth_service.ActionCartRemove), 
                    cartCtrl.RemoveCartItem(),
                )

                cart.DELETE("/clear", 
                    middlewares.AuthorizationMiddleware(auth_service.ActionCartClear), 
                    cartCtrl.ClearCart(),
                )
            }

            payment := protected.Group("/payment")
            {
                payment.POST("/stripe", paymentCtrl.CheckoutOrder())
            }
        }
    }
}