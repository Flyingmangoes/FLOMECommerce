package server

import (
	"backend/src/controllers"
	"backend/src/middlewares"
	"github.com/gin-gonic/gin"
)

func registerRoutes(r *gin.Engine, s *ServerContext) {
	userCtrl := &controllers.UserContext{
        Users:     s.Users,
        Tokens:    s.Tokens,
        JWTSecret: s.JWTSecret,
    }

    storeCtrl := &controllers.StoreContext{
		Users: s.Users,
        Stores: s.Stores,
        JWTSecret: s.JWTSecret,
    }

    // productCtrl := &controllers.ProductContext{
    //     Products:  s.Products,
    //     JWTSecret: s.JWTSecret,
    // }

    // public auth routes
    auth := r.Group("/v1/auth")
    {
        auth.GET("/users",     userCtrl.LoginUser())
        auth.POST("/users",    userCtrl.RegisterUser())
        auth.POST("/refresh", userCtrl.Refresh())
        auth.POST("/logout",  userCtrl.LogoutUser())
    }

    // protected user routes
    protected := r.Group("/v1")
    protected.Use(middlewares.AuthMiddlewares(string(s.JWTSecret)))
    {
        protected.PUT("/users", userCtrl.UpdateUser())
        protected.DELETE("/users", userCtrl.DeleteUser())

        protected.POST("/auth/stores",   storeCtrl.RegisterStore())
		protected.GET("/auth/stores", storeCtrl.LoginStore())
        protected.PUT("/stores",    storeCtrl.UpdateStore())
        protected.DELETE("/stores", storeCtrl.DeleteStore())
    }

    // guest := r.Group("/v1/api")
    // {
    //     guest.GET("/products", productCtrl.SearchProducts())
    //     guest.GET("/stores",   storeCtrl.SearchStores())
    // }

    s_protected := r.Group("/v1/store")
    s_protected.Use(middlewares.AuthMiddlewares(string(s.JWTSecret)))
    s_protected.Use(middlewares.CheckForStore(s.Stores))
    {
        //s_protected.POST("/products",   productCtrl.CreateProducts())
        //s_protected.PUT("/products",    productCtrl.UpdateProducts())
        //s_protected.DELETE("/products", productCtrl.DeleteProducts())
    }
}