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
	users := r.Group("/users")
	{
		users.POST("/new", h.Register)
		users.POST("/login", h.Login)
		users.POST("/refresh", h.Refresh)
	}

	// Access-token-protected routes
	protected := r.Group("/users")
	protected.Use(middleware.JWTAuth())
	{
		protected.GET("/profile", h.UserProfile)
		protected.DELETE("/:uid", h.DeleteUser)
	}

	// Consider protecting or removing this route later
	r.GET("/users/:uid", h.GetUser)
}
