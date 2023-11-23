package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

func MainRouter(db *sqlx.DB) *gin.Engine {
	router := gin.Default()
	
	RouterUser(router, db)
	RouterProduct(router, db)
	RouterOrder(router, db)
	RouterPromo(router, db)
	return router
}