package routers

import (
	"akbarsyarif/coffeeshopgolang/internal/handlers"
	"akbarsyarif/coffeeshopgolang/internal/middlewares"
	"akbarsyarif/coffeeshopgolang/internal/repositories"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

func RouterPromo(g *gin.Engine, db *sqlx.DB) {
	route := g.Group("/promo")
	repository := repositories.InitializePromoRepository(db)
	handler := handlers.InitializePromoHandler(repository)

	route.GET("/", handler.GetPromo)
	route.POST("/", middlewares.JWTGate(db, "admin"), handler.CreateNewPromo)
	route.PATCH("/:promoid", middlewares.JWTGate(db, "admin"), handler.UpdatePromo)
	route.DELETE("/:promoid", middlewares.JWTGate(db, "admin"), handler.DeletePromo)
}