package routers

import (
	"akbarsyarif/coffeeshopgolang/internal/handlers"
	"akbarsyarif/coffeeshopgolang/internal/repositories"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

func RouterPromo(g *gin.Engine, db *sqlx.DB) {
	route := g.Group("/promo")
	repository := repositories.InitializePromoRepository(db)
	handler := handlers.InitializePromoHandler(repository)

	route.GET("/", handler.GetPromo)
	route.POST("/", handler.CreateNewPromo)
	route.PATCH("/:promoid", handler.UpdatePromo)
	route.DELETE("/:promoid", handler.DeletePromo)
}