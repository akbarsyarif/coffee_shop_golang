package handlers

import (
	"akbarsyarif/coffeeshopgolang/internal/models"
	"akbarsyarif/coffeeshopgolang/internal/repositories"

	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type HandlerOrder struct {
	*repositories.OrderRepository
}

func InitializeOrderHandler(r *repositories.OrderRepository) *HandlerOrder {
	return &HandlerOrder{r}
}

func (h *HandlerOrder) GetAllOrder(ctx *gin.Context)  {
	result, err := h.RepositoryGetAllOrder()
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

func (h *HandlerOrder) GetOrderPerUser(ctx *gin.Context)  {
	userId := ctx.Param("userid")
	status := ctx.Query("status")
	page := ctx.Query("page")
	if status == ""{
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "You Must Input Status",
		})
		return
	}
	if page == ""{
		page = "1"
	}

	result, err := h.RepositoryGetOrderPerUser(userId, status, page)
	if err != nil {
		log.Println(err)
		ctx.JSON(http.StatusInternalServerError, err)
		// panic(err)
		return
	}
	// log.Println(result)
	if len(result) == 0 {
		ctx.JSON(http.StatusNotFound, gin.H{
			"message": "No Order Found",
		})
		// panic(err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Success",
		"result": result,
	})
}

func (h *HandlerOrder) GetOrderDetail(ctx *gin.Context)  {
	orderId := ctx.Param("orderid")

	result, err := h.RepositoryGetOrderDetail(orderId)
	if err != nil {
		log.Println(err)
		ctx.JSON(http.StatusInternalServerError, err)
		// panic(err)
		return
	}
	// log.Println(result)
	if len(result) == 0 {
		ctx.JSON(http.StatusNotFound, gin.H{
			"message": "No Order Found",
		})
		// panic(err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Success",
		"result": result,
	})
}

func (h *HandlerOrder) CreateNewOrder(ctx *gin.Context)  {
	var body models.OrderModel
	var orderId string
	if err := ctx.ShouldBind(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, err)
		return
	}
	if body.Shipping == "" || body.Total == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "You Must Input shipping_name and total",
		})
		return
	}
	log.Println(body)

	tx, err := h.Beginx()
	if err != nil {
		log.Println(err)
		ctx.JSON(http.StatusInternalServerError, err)
		// panic(err)
		return
	}
	defer tx.Rollback()

	result, err := h.RepositoryCreateOrder(&body, tx)
	for result.Next() {
		var Id string
		err = result.Scan(&Id)
		if Id != "" {
			log.Println(Id)
			orderId = Id
			break
		}
	}
	if err != nil {
		log.Println(err)
		ctx.JSON(http.StatusInternalServerError, err)
		// panic(err)
		return
	}
	result.Rows.Close()

	_, err = h.RepositoryCreateOrderProduct(&body, tx, orderId)
	if err != nil {
		log.Println(err)
		ctx.JSON(http.StatusInternalServerError, err)
		// panic(err)
		return
	}

	if err = tx.Commit(); err != nil {
		log.Println(err)
		ctx.JSON(http.StatusInternalServerError, err)
		// panic(err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Create Order Success",
	})
}

func (h *HandlerOrder) UpdateOrder(ctx *gin.Context)  {
	var body models.OrderModel
	orderId := ctx.Param("orderid")
	if err := ctx.ShouldBind(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, err)
	}

	result, err := h.RepositoryUpdateOrder(&body, orderId)
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
			"Message": "Order Not Found",
		})
		return
	}
	if err != nil {
		log.Println(err)
		ctx.JSON(http.StatusInternalServerError, err)
		// panic(err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Update Order Success",
	})
}

func (h *HandlerOrder) DeleteOrder(ctx *gin.Context)  {
	orderId := ctx.Param("orderid")
	tx, err := h.Beginx()
	if err != nil {
		log.Println(err)
		ctx.JSON(http.StatusInternalServerError, err)
		// panic(err)
		return
	}
	defer tx.Rollback()

	result, err := h.RepositoryDeleteOrder(orderId, tx)
	if err != nil {
		log.Println(err)
		ctx.JSON(http.StatusInternalServerError, err)
		// panic(err)
		return
	}
	result, err = h.RepositoryDeleteOrderProduct(orderId, tx)
	if err != nil {
		log.Println(err)
		ctx.JSON(http.StatusInternalServerError, err)
		// panic(err)
		return
	}

	if check, _ := result.RowsAffected(); check == 0 {
		ctx.JSON(http.StatusNotFound, gin.H{
			"Message": "Order Not Found",
		})
		return
	}

	if err = tx.Commit(); err != nil {
		log.Println(err)
		ctx.JSON(http.StatusInternalServerError, err)
		// panic(err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Delete Order Success",
	})
}