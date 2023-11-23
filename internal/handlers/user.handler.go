package handlers

import (
	"akbarsyarif/coffeeshopgolang/internal/models"
	"akbarsyarif/coffeeshopgolang/internal/repositories"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type HandlerUser struct {
	*repositories.UserRepository
}

func InitializeUserHandler(r *repositories.UserRepository) *HandlerUser {
	return &HandlerUser{r}
}

func (h *HandlerUser) GetAllUser(ctx *gin.Context)  {
	result, err := h.RepositoryGetAllUser()
	if err != nil {
		log.Println(err)
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"message": "Success",
		"result": result,
	})
}

func (h *HandlerUser) GetUserDetail(ctx *gin.Context)  {
	userId := ctx.Param("userid")
	result, err := h.RepositoryGetUserDetail(userId)
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

func (h *HandlerUser) CreateNewUser(ctx *gin.Context)  {
	var body models.UserModel
	var userId string
	if err := ctx.ShouldBind(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, err)
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
	
	result, err := h.RepositoryCreateUser(&body, tx)
	for result.Next() {
		var Id string
		err = result.Scan(&Id)
		if Id != "" {
			log.Println(Id)
			userId = Id
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

	_, err = h.RepositoryCreateUserProfile(userId, tx)
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
		"message": "Create User Success",
	})
}

func (h *HandlerUser) UpdateUser(ctx *gin.Context)  {
	var body models.UserModel
	userId := ctx.Param("userid")
	if err := ctx.ShouldBind(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, err)
	}

	result, err := h.RepositoryUpdateUser(&body, userId)
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
			"Message": "User Not Found",
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
		"message": "Update User Success",
		"result": fmt.Sprintf("%d row updated in database", res),
	})
}

func (h *HandlerUser) DeleteUser(ctx *gin.Context)  {
	userId := ctx.Param("userid")
	result, err := h.RepositoryDeleteUser(userId)
	if err != nil {
		log.Println(err)
		ctx.JSON(http.StatusInternalServerError, err)
		// panic(err)
		return
	}

	if check, _ := result.RowsAffected(); check == 0 {
		ctx.JSON(http.StatusNotFound, gin.H{
			"Message": "User Not Found",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Delete User Success",
	})
}