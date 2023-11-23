package routers

import (
	"akbarsyarif/coffeeshopgolang/internal/handlers"
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
	route.POST("/", handler.CreateNewProduct)
	route.PATCH("/:productid", handler.UpdateProduct)
	route.DELETE("/:productid", handler.DeleteProduct)
}