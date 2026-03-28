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

	// v1 auth
	auth := r.Group("/v1/auth" ) 
	{
		auth.GET("/user", userCtrl.Login())
		auth.POST("/user", userCtrl.Register())
		auth.POST("/refresh", userCtrl.Refresh())
        auth.POST("/logout", userCtrl.Logout())
	}

	protected := r.Group("/v1")
    protected.Use(middlewares.AuthMiddlewares(string(s.JWTSecret)))				
    {
        protected.PUT("/user", userCtrl.Update())
        protected.DELETE("/user", userCtrl.Delete())
    }
}