package routers

import (
	"akbarsyarif/coffeeshopgolang/internal/handlers"
	"akbarsyarif/coffeeshopgolang/internal/middlewares"
	"akbarsyarif/coffeeshopgolang/internal/repositories"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

func RouterUser(g *gin.Engine, db *sqlx.DB) {
	route := g.Group("/user")
	repository := repositories.InitializeUserRepository(db)
	handler := handlers.InitializeUserHandler(repository)

	route.GET("/", middlewares.JWTGate(db, "admin"), handler.GetAllUser)
	route.GET("/detail", middlewares.JWTGate(db, "user", "admin"), handler.GetUserDetail)
	route.POST("/", middlewares.JWTGate(db, "admin"), handler.CreateNewUser)
	route.PATCH("/:userid", middlewares.JWTGate(db, "user", "admin"), handler.UpdateUser)
	route.DELETE("/:userid", middlewares.JWTGate(db, "admin"), handler.DeleteUser)
}