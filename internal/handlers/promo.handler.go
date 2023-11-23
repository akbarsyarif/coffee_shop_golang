package handlers

import (
	"akbarsyarif/coffeeshopgolang/internal/models"
	"akbarsyarif/coffeeshopgolang/internal/repositories"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type HandlerPromo struct {
	*repositories.PromoRepository
}

func InitializePromoHandler(r *repositories.PromoRepository) *HandlerPromo {
	return &HandlerPromo{r}
}

func (h *HandlerPromo) GetPromo(ctx *gin.Context)  {
	result, err := h.RepositoryGetPromo()
	if err != nil {
		log.Println(err)
		ctx.JSON(http.StatusInternalServerError, err)
		// panic(err)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"message": "Success",
		"result": result,
	})
}

func (h *HandlerPromo) CreateNewPromo(ctx *gin.Context)  {
	var body models.PromoModel
	if err := ctx.ShouldBind(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, err)
	}
	log.Println(body)

	result, err := h.RepositoryCreatePromo(&body)
	if err != nil {
		log.Println(err)
		ctx.JSON(http.StatusInternalServerError, err)
		// panic(err)
		return
	}
	res, _ := result.RowsAffected()

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Create Promo Success",
		"result": fmt.Sprintf("%d row added to database", res),
	})
}

func (h *HandlerPromo) UpdatePromo(ctx *gin.Context)  {
	var body models.PromoModel
	promoId := ctx.Param("promoid")
	if err := ctx.ShouldBind(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, err)
	}

	result, err := h.RepositoryUpdatePromo(&body, promoId)
	// if err.Error() == "Please Input at Least One Change" {
	// 	log.Println(err)
	// 	ctx.JSON(http.StatusBadRequest, gin.H{
	// 		"message": err.Error(),
	// 	})
	// 	// panic(err)
	// 	return
	// }
	if check, _ := result.RowsAffected(); check == 0 {
		ctx.JSON(http.StatusNotFound, gin.H{
			"Message": "Promo Not Found",
		})
		return
	}
	if err != nil {
		log.Println(err)
		ctx.JSON(http.StatusInternalServerError, err)
		// panic(err)
		return
	}
	res, _ := result.RowsAffected()


	ctx.JSON(http.StatusOK, gin.H{
		"message": "Update Promo Success",
		"result": fmt.Sprintf("%d row updated in database", res),
	})
}

func (h *HandlerPromo) DeletePromo(ctx *gin.Context)  {
	promoId := ctx.Param("promoid")
	result, err := h.RepositoryDeletePromo(promoId)
	if err != nil {
		log.Println(err)
		ctx.JSON(http.StatusInternalServerError, err)
		// panic(err)
		return
	}

	if check, _ := result.RowsAffected(); check == 0 {
		ctx.JSON(http.StatusNotFound, gin.H{
			"Message": "Promo Not Found",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Delete Promo Success",
	})
}