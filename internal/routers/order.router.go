package routers

import (
	"akbarsyarif/coffeeshopgolang/internal/handlers"
	"akbarsyarif/coffeeshopgolang/internal/middlewares"
	"akbarsyarif/coffeeshopgolang/internal/repositories"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

func RouterOrder(g *gin.Engine, db *sqlx.DB) {
	route := g.Group("/order")
	repository := repositories.InitializeOrderRepository(db)
	handler := handlers.InitializeOrderHandler(repository)

	route.GET("/", middlewares.JWTGate(db, "admin"), handler.GetAllOrder)
	route.GET("/user", middlewares.JWTGate(db, "user"), handler.GetOrderPerUser)
	route.GET("/detail/:orderid", middlewares.JWTGate(db, "user"), handler.GetOrderDetail)
	route.POST("/", middlewares.JWTGate(db, "user"), handler.CreateNewOrder)
	route.PATCH("/:orderid", middlewares.JWTGate(db, "admin"), handler.UpdateOrder)
	route.DELETE("/:orderid", middlewares.JWTGate(db, "admin"), handler.DeleteOrder)
	// route.PATCH("/:productid", handler.UpdateProduct)
	// route.DELETE("/:productid", handler.DeleteProduct)
}