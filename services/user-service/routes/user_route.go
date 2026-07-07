package routes

import (
	"github.com/abhinandpn/UnVocal/services/user-service/handler"
	"github.com/abhinandpn/UnVocal/services/user-service/middleware"
	"github.com/abhinandpn/UnVocal/services/user-service/repository"
	"github.com/abhinandpn/UnVocal/services/user-service/service"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func SetupRoutes(r *gin.Engine, db *pgxpool.Pool) {

	repo := repository.NewUserRepository(db)
	svc := service.NewUserService(repo)
	h := handler.NewUserHandler(svc)

	// Public routes
	user := r.Group("/users")
	{
		user.POST("/new", h.Register)
		user.POST("/login", h.Login)
	}

	// Protected routes
	protected := r.Group("/users")
	protected.Use(middleware.JWTAuth())
	{
		protected.GET("/profile", h.UserProfile)
		protected.PUT("/:uid", h.UpdateUser)
		protected.DELETE("/:uid", h.DeleteUser)
		// protected.POST("/logout", h.Logout)
	}

	// Public route
	r.GET("/users/:uid", h.GetUser)
}
