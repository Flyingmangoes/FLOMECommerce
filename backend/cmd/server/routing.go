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
    }

    s_protected := r.Group("/v1/ston")
    s_protected.Use(middlewares.AuthMiddlewares(string(s.JWTSecret)))
    s_protected.Use(middlewares.CheckForStore(s.Stores))
    {
        protected.POST("/stores", storeCtrl.RegisterStore())
        s_protected.PUT("/stores",    storeCtrl.UpdateStore())
        s_protected.DELETE("/stores", storeCtrl.DeleteStore())
    }
}