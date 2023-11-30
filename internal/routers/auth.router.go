package routers

import (
	"akbarsyarif/coffeeshopgolang/internal/handlers"
	"akbarsyarif/coffeeshopgolang/internal/middlewares"
	"akbarsyarif/coffeeshopgolang/internal/repositories"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

func RouterAuth(g *gin.Engine, db *sqlx.DB) {
	route := g.Group("/auth")
	repository := repositories.InitializeAuthRepository(db)
	handler := handlers.InitializeAuthHandler(repository)

	route.POST("/", handler.Register)
	route.POST("/login", handler.Login)
	route.DELETE("/", middlewares.JWTGate(db, "user", "admin"), handler.Logout)
	// route.PATCH("/:productid", handler.UpdateProduct)
	// route.DELETE("/:productid", handler.DeleteProduct)
}