package routers

import (
	"akbarsyarif/coffeeshopgolang/internal/handlers"
	"akbarsyarif/coffeeshopgolang/internal/repositories"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

func RouterOrder(g *gin.Engine, db *sqlx.DB) {
	route := g.Group("/order")
	repository := repositories.InitializeOrderRepository(db)
	handler := handlers.InitializeOrderHandler(repository)

	route.GET("/", handler.GetAllOrder)
	route.GET("/:userid", handler.GetOrderPerUser)
	route.GET("/detail/:orderid", handler.GetOrderDetail)
	route.POST("/", handler.CreateNewOrder)
	route.PATCH("/:orderid", handler.UpdateOrder)
	route.DELETE("/:orderid", handler.DeleteOrder)
	// route.PATCH("/:productid", handler.UpdateProduct)
	// route.DELETE("/:productid", handler.DeleteProduct)
}