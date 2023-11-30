package handlers

import (
	"akbarsyarif/coffeeshopgolang/internal/models"
	"akbarsyarif/coffeeshopgolang/internal/repositories"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
)

type HandlerPromo struct {
	*repositories.PromoRepository
}

func InitializePromoHandler(r *repositories.PromoRepository) *HandlerPromo {
	return &HandlerPromo{r}
}

func (h *HandlerPromo) GetPromo(ctx *gin.Context)  {
	page := ctx.Query("page")
	if page == ""{
		page = "1"
	}

	result, err := h.RepositoryGetPromo(page)
	if err != nil {
		log.Println(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Internal Server Error",
		})
		return
	}
	if len(result) == 0 {
		ctx.JSON(http.StatusNotFound, gin.H{
			"message": "Promo Not Found",
		})
		return
	}

	metaRes, err := h.RepositoryCountPromo()
	if err != nil {
		log.Println(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Internal Server Error",
		})
		return
	}
	a:= ctx.Request.URL.String()
	b:= ctx.Request.Host
	pages, err := strconv.Atoi(page)
	// log.Println(a)
	// log.Println(b)
	// log.Println(pages)
	metaData := metaConfig(metaRes[0], pages, b, a)

	
	ctx.JSON(http.StatusOK, gin.H{
		"message": "Success",
		"result": gin.H{
			"data": result,
			"meta": gin.H{
				"page": metaData[0],
				"total_page": metaData[1],
				"total_data": metaData[2],
				"next": metaData[3],
				"prev": metaData[4],
			},
		},
	})
}

func (h *HandlerPromo) CreateNewPromo(ctx *gin.Context)  {
	var body models.PromoModel
	if err := ctx.ShouldBind(&body); err != nil {
			log.Println(err)
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"message": "Internal Server Error",
			})
			return		
	}
	if _, err := govalidator.ValidateStruct(body); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"Message" : "Please Check your Input",
		})
		return
	}

	_, err := h.RepositoryCreatePromo(&body)
	if err != nil {
		log.Println(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Internal Server Error",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Create Promo Success",
		"result": gin.H{
			"promo_name": body.Promo_name,
			"description" :body.Description,
		},
	})
}

func (h *HandlerPromo) UpdatePromo(ctx *gin.Context)  {
	var body models.PromoModel
	promoId := ctx.Param("promoid")
	if err := ctx.ShouldBind(&body); err != nil {
		log.Println(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Internal Server Error here?",
		})
		return		
	}
	if _, err := govalidator.ValidateStruct(body); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"Message" : "Please Check your Input",
		})
		return
	}
	if body.Promo_name == "" && body.Description == "" && body.Flat_amount == 0 && body.Discount_type == "" &&body.Percent_amount == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Please Input at Least One Change",
		})
		return
	}
	if body.Discount_type != "" {
		if body.Discount_type != "flat" && body.Discount_type != "percent" {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": "Invalid Discount Type",
			})
			return
		}
		if body.Discount_type == "flat" && body.Percent_amount > 0 {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": "Flat Discount Type",
			})
			return
		}
		if body.Discount_type == "percent" && body.Flat_amount > 0 {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": "Percent Discount Type",
			})
			return
		}
	}

	result, err := h.RepositoryUpdatePromo(&body, promoId)
	if err != nil {
		log.Println(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Internal Server Error",
		})
		return
	}
	if check, _ := result.RowsAffected(); check == 0 {
		ctx.JSON(http.StatusNotFound, gin.H{
			"Message": "Promo Not Found",
		})
		return
	}
	// res, _ := result.RowsAffected()


	ctx.JSON(http.StatusOK, gin.H{
		"message": "Update Promo Success",
		"result": body,
	})
}

func (h *HandlerPromo) DeletePromo(ctx *gin.Context)  {
	promoId := ctx.Param("promoid")
	result, err := h.RepositoryDeletePromo(promoId)
	if err != nil {
		log.Println(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Internal Server Error",
		})
		return
	}

	if check, _ := result.RowsAffected(); check == 0 {
		ctx.JSON(http.StatusNotFound, gin.H{
			"Message": "Promo Not Found",
		})
		return
	}
	message := fmt.Sprintf("Delete on Promo Id %s Success", promoId)

	ctx.JSON(http.StatusOK, gin.H{
		"message": message,
	})
}