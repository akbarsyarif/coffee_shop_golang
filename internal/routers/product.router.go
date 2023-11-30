package routers

import (
	"akbarsyarif/coffeeshopgolang/internal/handlers"
	"akbarsyarif/coffeeshopgolang/internal/middlewares"
	"akbarsyarif/coffeeshopgolang/internal/repositories"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

func RouterProduct(g *gin.Engine, db *sqlx.DB) {
	route := g.Group("/product")
	repository := repositories.InitializeRepository(db)
	handler := handlers.InitializeHandler(repository)

	route.GET("/", handler.GetAllProduct)
	route.GET("/:productid", handler.GetProductDetail)
	route.POST("/", middlewares.JWTGate(db, "admin"), handler.CreateNewProduct)
	route.PATCH("/:productid", middlewares.JWTGate(db, "admin"), handler.UpdateProduct)
	route.DELETE("/:productid", middlewares.JWTGate(db, "admin"), handler.DeleteProduct)
}