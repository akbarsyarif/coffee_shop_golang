package routers

import (
	"akbarsyarif/coffeeshopgolang/internal/handlers"
	"akbarsyarif/coffeeshopgolang/internal/repositories"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

func RouterUser(g *gin.Engine, db *sqlx.DB) {
	route := g.Group("/user")
	repository := repositories.InitializeUserRepository(db)
	handler := handlers.InitializeUserHandler(repository)

	route.GET("/", handler.GetAllUser)
	route.GET("/:userid", handler.GetUserDetail)
	route.POST("/", handler.CreateNewUser)
	route.PATCH("/:userid", handler.UpdateUser)
	route.DELETE("/:userid", handler.DeleteUser)
}