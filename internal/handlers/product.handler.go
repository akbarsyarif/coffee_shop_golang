package handlers

import (
	"akbarsyarif/coffeeshopgolang/internal/helpers"
	"akbarsyarif/coffeeshopgolang/internal/models"
	"akbarsyarif/coffeeshopgolang/internal/repositories"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
)

type HandlerProduct struct {
	*repositories.ProductRepository
}

func InitializeHandler(r *repositories.ProductRepository) *HandlerProduct {
	return &HandlerProduct{r}
}

func (h *HandlerProduct) GetAllProduct(ctx *gin.Context)  {
	var params models.ProductParams
	if err := ctx.ShouldBindQuery(&params); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Internal Server Error",
		})
	return
	}
	if _, err := govalidator.ValidateStruct(params); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"Message" : "Invalid Query Params",
			"err": err.Error(),
		})
		return
	}
	if params.Page == "" {
		params.Page = "1"
	}

	var res1 int
	var res2 int
	var err error
	log.Println(res1, res2)
	
	if params.Min_price != "" {
		res1, err = strconv.Atoi(params.Min_price)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": "Input Number Only For Price",
			})
			return
		}
	}
	if params.Max_price != "" {
		res2, err = strconv.Atoi(params.Max_price)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": "Input Number Only For Price",
			})
			return
		}
	}
	if params.Min_price != "" && params.Max_price != "" {
		if res1 > res2 {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": "Min Price Must Be Lower Than Max Price",
			})
			return
		}
	}

	result, err := h.RepositoryGetAllProduct(&params)
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

	metaRes, err := h.RepositoryCountProduct(&params)
	if err != nil {
		log.Println(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Internal Server Error",
		})
		return
	}
	a:= ctx.Request.URL.String()
	b:= ctx.Request.Host
	pages, err := strconv.Atoi(params.Page)
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

func (h *HandlerProduct) GetProductDetail(ctx *gin.Context)  {
	productId := ctx.Param("productid")
	result, err := h.RepositoryGetProductDetail(productId)
	if err != nil {
		log.Println(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Internal Server Error",
		})
		// panic(err)
		return
	}
	if len(result) == 0 {
		ctx.JSON(http.StatusNotFound, gin.H{
			"message": "Product Not Found",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Success",
		"result": result,
	})
}


func (h *HandlerProduct) CreateNewProduct(ctx *gin.Context)  {
	var body models.ProductModel
	var productId string
	if err := ctx.ShouldBind(&body); err != nil {
		log.Println(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Internal Server Error?",
		})
		return
	}
	log.Println(body)

	result, err := h.RepositoryCreateProduct(&body)
	if err != nil {
		log.Println(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Internal Server Error",
		})
		return
	}
	for result.Next() {
		var Id string
		err = result.Scan(&Id)
		if Id != "" {
			log.Println(Id)
			productId = Id
			break
		}
	}
	if err != nil {
		log.Println(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Internal Server Error",
		})
		return
	}
	result.Rows.Close()
	
	var imageUrl string

	cld, err := helpers.InitCloudinary()
	if  err != nil {
		log.Println(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"Message": "Internal Server Error",
		})
		return
	}
	formFile, _ := ctx.FormFile("product")
	if formFile != nil {
		file, err := formFile.Open()
		if  err != nil {
			log.Println(err)
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"Message": "Internal Server Error",
			})
			return
		}
		defer file.Close()
		
		publicId := fmt.Sprintf("%s_%s-%s", "golang", "product", productId)
		folder := "coffeeshop/product"
		res, err := cld.Uploader(ctx, file, publicId, folder)
		if  err != nil {
			log.Println(err)
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"Message": "Internal Server Error",
			})
			return
		}
		imageUrl = res.SecureURL
	}

	_, err = h.RepositoryInsertimage(imageUrl, productId)
	if  err != nil {
		log.Println(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"Message": "Internal Server Error",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Create Product Success",
		"result": gin.H{
			"product_image": imageUrl,
			"product_name": body.Product_name,
			"description": body.Description,
			"price": body.Price,
		},
	})
}

func (h *HandlerProduct) UpdateProduct(ctx *gin.Context)  {
	var body models.ProductModel
	productId := ctx.Param("productid")
	var imageUrl string

	cld, err := helpers.InitCloudinary()
	if  err != nil {
		log.Println(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"Message": "Internal Server Error",
		})
		return
	}
	formFile, _ := ctx.FormFile("product")
	if formFile != nil {
		file, err := formFile.Open()
		if  err != nil {
			log.Println(err)
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"Message": "Internal Server Error",
			})
			return
		}
		defer file.Close()
		
		publicId := fmt.Sprintf("%s_%s-%s", "golang", "product", productId)
		folder := "coffeeshop/product"
		res, err := cld.Uploader(ctx, file, publicId, folder)
		if  err != nil {
			log.Println(err)
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"Message": "Internal Server Error",
			})
			return
		}
		imageUrl = res.SecureURL
	}

	if err := ctx.ShouldBind(&body); err != nil {
		log.Println(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Internal Server Error?",
		})
		return
	}
	if body.Product_name == "" && body.Description == "" && body.Price == 0 && body.Category == "" && body.Promo == "" && imageUrl == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Please Input At Least One Change",
		})
		return
	}

	result, err := h.RepositoryUpdateProduct(&body, productId, imageUrl)
	if err != nil {
		log.Println(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Internal Server Error?",
		})
		return
	}
	if check, _ := result.RowsAffected(); check == 0 {
		ctx.JSON(http.StatusNotFound, gin.H{
			"Message": "Product Not Found",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Update Product Success",
		"result": gin.H{
			"product_name": body.Product_name,
			"price": body.Price,
			"description": body.Description,
		},
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
	message := fmt.Sprintf("Delete on Product Id %s Success", productId)

	ctx.JSON(http.StatusOK, gin.H{
		"message": message,
	})
}