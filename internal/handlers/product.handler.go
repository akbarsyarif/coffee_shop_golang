package handlers

import (
	"akbarsyarif/coffeeshopgolang/internal/models"
	"akbarsyarif/coffeeshopgolang/internal/repositories"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type HandlerProduct struct {
	*repositories.ProductRepository
}

func InitializeHandler(r *repositories.ProductRepository) *HandlerProduct {
	return &HandlerProduct{r}
}

func (h *HandlerProduct) GetAllProduct(ctx *gin.Context)  {
	// search := ctx.Query("search")
	// category := ctx.Query("category")
	log.Println(ctx.Request.URL.Query())
	min_price := ctx.Query("min_price")
	max_price := ctx.Query("max_price")
	page := ctx.Query("page")
	if page == "" {
		page = "1"
	}

	var res1 int
	var res2 int
	var err error
	log.Println(res1, res2)
	
	if min_price != "" {
		res1, err = strconv.Atoi(min_price)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": "Input Number Only For Price",
			})
			return
		}
	}
	if max_price != "" {
		res2, err = strconv.Atoi(max_price)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": "Input Number Only For Price",
			})
			return
		}
	}
	if min_price != "" && max_price != "" {
		if res1 > res2 {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": "Min Price Must Be Lower Than Max Price",
			})
			return
		}
	}

	result, err := h.RepositoryGetAllProduct()
	if err != nil {
		log.Println(err)
		ctx.JSON(http.StatusInternalServerError, err)
		// panic(err)
		return
	}
	if len(result) == 0 {
		ctx.JSON(http.StatusNotFound, gin.H{
			"message": "No Product Found",
			"result": result,
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Success",
		"result": result,
	})
}

func (h *HandlerProduct) GetProductDetail(ctx *gin.Context)  {
	productId := ctx.Param("productid")
	result, err := h.RepositoryGetProductDetail(productId)
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


func (h *HandlerProduct) CreateNewProduct(ctx *gin.Context)  {
	var body models.ProductModel
	if err := ctx.ShouldBind(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, err)
	}
	log.Println(body)

	result, err := h.RepositoryCreateProduct(&body)
	if err != nil {
		log.Println(err)
		ctx.JSON(http.StatusInternalServerError, err)
		// panic(err)
		return
	}
	res, _ := result.RowsAffected()

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Create Product Success",
		"result": fmt.Sprintf("%d row added to database", res),
	})
}

func (h *HandlerProduct) UpdateProduct(ctx *gin.Context)  {
	var body models.ProductModel
	productId := ctx.Param("productid")
	if err := ctx.ShouldBind(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, err)
	}

	result, err := h.RepositoryUpdateProduct(&body, productId)
	if err.Error() == "Please Input at Least One Change" {
		log.Println(err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		// panic(err)
		return
	}
	if check, _ := result.RowsAffected(); check == 0 {
		ctx.JSON(http.StatusNotFound, gin.H{
			"Message": "Product Not Found",
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
		"message": "Update Product Success",
		"result": fmt.Sprintf("%d row updated in database", res),
	})
}

func (h *HandlerProduct) DeleteProduct(ctx *gin.Context)  {
	productId := ctx.Param("productid")
	result, err := h.RepositoryDeleteProduct(productId)
	if err != nil {
		log.Println(err)
		ctx.JSON(http.StatusInternalServerError, err)
		// panic(err)
		return
	}

	if check, _ := result.RowsAffected(); check == 0 {
		ctx.JSON(http.StatusNotFound, gin.H{
			"Message": "Product Not Found",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Delete Product Success",
	})
}