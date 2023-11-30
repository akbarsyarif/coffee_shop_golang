package handlers

import (
	"akbarsyarif/coffeeshopgolang/internal/models"
	"akbarsyarif/coffeeshopgolang/internal/repositories"
	"akbarsyarif/coffeeshopgolang/pkg"
	"fmt"
	"math"
	"strconv"

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

func metaConfig(totalData, page int, host, query string) ([]interface{}){
	var result []interface{}
	// url := fmt.Sprintf("%s%s", host, query)
    totalPage := math.Ceil(float64(totalData) / 5)
    isLastPage := page >= int(totalPage)
	var next any
	var prev any
	result = append(result, page)
	result = append(result, totalPage)
	result = append(result, totalData)
	if isLastPage && page == 1 {
		next = nil
		prev = nil
		result = append(result, next)
		result = append(result, prev)
		return result
	}
	if isLastPage {
		next = nil
		prev = fmt.Sprintf("%s%s", host, query[:len(query)-1] + strconv.Itoa(page-1))
		result = append(result, next)
		result = append(result, prev)
		return result
	}
	if page == 1 {
		next = fmt.Sprintf("%s%s", host, query[:len(query)-1] + strconv.Itoa(page+1))
		prev = nil
		result = append(result, next)
		result = append(result, prev)
		return result
	}

	next = fmt.Sprintf("%s%s", host, query[:len(query)-1] + strconv.Itoa(page+1))
	prev = fmt.Sprintf("%s%s", host, query[:len(query)-1] + strconv.Itoa(page-1))
	result = append(result, next)
	result = append(result, prev)
	return result
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
	value, ok := ctx.Get("Payload")
	if !ok {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Internal Server Error (Payload)",
		})
		return
	}
	userId := fmt.Sprint(value.(*pkg.Claims).Id)
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
		return
	}
	// log.Println(result)
	if len(result) == 0 {
		log.Println(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Internal Server Error",
		})
		return
	}
	metaRes, err := h.RepositoryCountOrder(userId, status)
	if err != nil {
		log.Println(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Internal Server Error",
		})
		return
	}
	// println(metaRes[0])
	a:= ctx.Request.URL.String()
	b:= ctx.Request.Host
	pages, err := strconv.Atoi(page)
	// log.Println(a)
	// log.Println(b)
	// log.Println(pages)
	metaData := metaConfig(metaRes[0], pages, b, a)
	// log.Println(metaData...)

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

func (h *HandlerOrder) GetOrderDetail(ctx *gin.Context)  {
	value, ok := ctx.Get("Payload")
	if !ok {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Internal Server Error (Payload)",
		})
		return
	}
	userId := fmt.Sprint(value.(*pkg.Claims).Id)
	orderId := ctx.Param("orderid")

	result, err := h.RepositoryGetOrderDetail(orderId, userId)
	if err != nil {
		log.Println(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Internal Server Error",
		})
		return
	}
	if len(result) == 0 {
		ctx.JSON(http.StatusNotFound, gin.H{
			"message": "No Order Found",
		})
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
	value, ok := ctx.Get("Payload")
	if !ok {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Internal Server Error (Payload)",
		})
		return
	}
	userId := fmt.Sprint(value.(*pkg.Claims).Id)
	body.User_id = userId
	// log.Println(userId)

	tx, err := h.Beginx()
	if err != nil {
		log.Println(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Internal Server Error",
		})
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
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Internal Server Error",
		})
		return
	}
	result.Rows.Close()

	_, err = h.RepositoryCreateOrderProduct(&body, tx, orderId)
	if err != nil {
		log.Println(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Internal Server Error",
		})
		return
	}

	if err = tx.Commit(); err != nil {
		log.Println(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Internal Server Error",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Create Order Success",
		"result": gin.H{
			"id": orderId,
			"total": body.Total,
			"shipping_name": body.Shipping,
			"product": body.Product,
		},
	})
}

func (h *HandlerOrder) UpdateOrder(ctx *gin.Context)  {
	var body models.OrderModel
	orderId := ctx.Param("orderid")
	if err := ctx.ShouldBind(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, err)
	}
	if body.Status == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Status Cannot Be Empty",
		})
		return
	}

	result, err := h.RepositoryUpdateOrder(&body, orderId)
	if err != nil {
		log.Println(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Internal Server Error",
		})
		return
	}

	if check, _ := result.RowsAffected(); check == 0 {
		ctx.JSON(http.StatusNotFound, gin.H{
			"Message": "Order Not Found",
		})
		return
	}
	if err != nil {
		log.Println(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Internal Server Error",
		})
		return
	}
	message := fmt.Sprintf("Update Status to %s on Order Id %s Success", body.Status, orderId)
	ctx.JSON(http.StatusOK, gin.H{
		"message": message,
	})
}

func (h *HandlerOrder) DeleteOrder(ctx *gin.Context)  {
	orderId := ctx.Param("orderid")
	tx, err := h.Beginx()
	if err != nil {
		log.Println(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Internal Server Error",
		})
		return
	}
	defer tx.Rollback()

	result, err := h.RepositoryDeleteOrder(orderId, tx)
	if err != nil {
		log.Println(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Internal Server Error",
		})
		return
	}
	result, err = h.RepositoryDeleteOrderProduct(orderId, tx)
	if err != nil {
		log.Println(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Internal Server Error",
		})
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
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Internal Server Error",
		})
		return
	}
	message := fmt.Sprintf("Delete on Order Id %s Success", orderId)

	ctx.JSON(http.StatusOK, gin.H{
		"message": message,
	})
}