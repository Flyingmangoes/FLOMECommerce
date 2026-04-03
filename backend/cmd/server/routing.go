package server

import (
	"backend/src/controllers"
	"backend/src/middlewares"
	"github.com/gin-gonic/gin"
)

func registerRoutes(r *gin.Engine, s *ServerContext) {
	userCtrl := &controllers.UserContext{
        Users: s.Users,
        Products:s.Products,
        Orders: s.Orders,
		Tokens: s.Tokens,
		JWTSecret: s.JWTSecret,
    }

	storeCtrl := &controllers.StoreContext{
        Products:s.Products,
        Orders: s.Orders,
	}

	// v1 auth
	auth := r.Group("/v1/auth" ) 
	{
		auth.GET("/user", userCtrl.Login())
		auth.POST("/user", userCtrl.Register())
		auth.POST("/refresh", userCtrl.Refresh())
        auth.POST("/logout", userCtrl.Logout())
	}

	// v1 normal searching
	guest := r.Group("/v1/api") 
	{
		guest.GET("/product", userCtrl.SearchProduct())
	}

	protected := r.Group("/v1")
    protected.Use(middlewares.AuthMiddlewares(string(s.JWTSecret)))				
    {
        protected.PUT("/user", userCtrl.Update())
        protected.DELETE("/user", userCtrl.Delete())
    }

	s_protected := r.Group("/v1/api/store")
	s_protected.Use(middlewares.AuthMiddlewares(string(s.JWTSecret)))
	s_protected.Use(middlewares.CheckForStore())
	{
		s_protected.POST("/product", storeCtrl.CreateProduct())
		s_protected.PUT("/product", storeCtrl.UpdateProduct())
		s_protected.DELETE("/product", storeCtrl.RemoveProduct())
	}
}