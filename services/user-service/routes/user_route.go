package routes

import (
	"github.com/abhinandpn/UnVocal/services/user-service/handler"
	"github.com/abhinandpn/UnVocal/services/user-service/repository"
	"github.com/abhinandpn/UnVocal/services/user-service/service"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func SetupRoutes(r *gin.Engine, db *pgxpool.Pool) {

	repo := repository.NewUserRepository(db)
	svc := service.NewUserService(repo)
	h := handler.NewUserHandler(svc)

	// CRUD
	user := r.Group("/users")
	{
		user.POST("/new", h.Register)      // Create a new user
		user.GET("/:uid", h.GetUser)       // Get user by user code
		user.PUT("/:uid", h.UpdateUser)    // Update user by user code
		user.DELETE("/:uid", h.DeleteUser) // Delete user by user code
	}
}
