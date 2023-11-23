package main

import (
	"akbarsyarif/coffeeshopgolang/internal/routers"
	"akbarsyarif/coffeeshopgolang/pkg"
	"log"
)

func main()  {
	database, err := pkg.Postgres()
	if err != nil {
		log.Fatal(err)
	}

	routers := routers.MainRouter(database)
	server := pkg.Server(routers)
	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}

// 2023/11/20 02:04:42 sql: Scan error on column index 7, name "promo_id": converting NULL to string is unsupported

// func main() {
// 	router := gin.Default()

// 	router.GET("/", func(c *gin.Context) {
// 		c.JSON(http.StatusOK, gin.H{
// 			"message": "Welcome",
// 		})
// 	})
// 	router.Run(("0.0.0.0:8000"))
// }