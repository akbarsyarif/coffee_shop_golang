package routers

import (
	"akbarsyarif/coffeeshopgolang/internal/middlewares"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

func MainRouter(db *sqlx.DB) *gin.Engine {
	router := gin.Default()
	
	router.Use(middlewares.CORSMiddleware)

	RouterUser(router, db)
	RouterProduct(router, db)
	RouterOrder(router, db)
	RouterPromo(router, db)
	RouterAuth(router, db)
	return router
}