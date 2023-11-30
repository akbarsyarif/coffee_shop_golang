package handlers

import (
	"akbarsyarif/coffeeshopgolang/internal/helpers"
	"akbarsyarif/coffeeshopgolang/internal/models"
	"akbarsyarif/coffeeshopgolang/internal/repositories"
	"akbarsyarif/coffeeshopgolang/pkg"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
)

type HandlerUser struct {
	*repositories.UserRepository
}

func InitializeUserHandler(r *repositories.UserRepository) *HandlerUser {
	return &HandlerUser{r}
}

func (h *HandlerUser) GetAllUser(ctx *gin.Context)  {
	page := ctx.Query("page")
	if page == ""{
		page = "1"
	}

	result, err := h.RepositoryGetAllUser(page)
	if err != nil {
		log.Println(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Internal Server Error",
		})
		return
	}
	if len(result) == 0 {
		ctx.JSON(http.StatusNotFound, gin.H{
			"message": "No User Found",
		})
	}

	metaRes, err := h.RepositoryCountUser()
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

func (h *HandlerUser) GetUserDetail(ctx *gin.Context)  {
	value, ok := ctx.Get("Payload")
	if !ok {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Internal Server Error (Payload)",
		})
		return
	}
	userId := fmt.Sprint(value.(*pkg.Claims).Id)

	result, err := h.RepositoryGetUserDetail(userId)
	if err != nil {
		log.Println(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Internal Server Error",
		})
		return
	}
	if len(result) == 0 {
		ctx.JSON(http.StatusNotFound, gin.H{
			"message": "User Not Found",
			"result": result,
		})
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
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Internal Server Error",
		})
		return
	}
	log.Println(body)
	if body.Full_name == "" || body.Email == "" || body.Pwd == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"Message" : "All Inputs Must Be Filled",
		})
		return
	}

	tx, err := h.Beginx()
	if err != nil {
		log.Println(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Internal Server Error",
		})
		return
	}

	defer tx.Rollback()
	
	result, err := h.RepositoryCreateUser(&body, tx)
	if err != nil {
		pgErr, _ := err.(*pq.Error)
		if pgErr.Code == "23505" {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": "Email is Registered",
			})
			return
		}
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
			userId = Id
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

	_, err = h.RepositoryCreateUserProfile(userId, tx)
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
		"message": "Create User Success",
		"result": gin.H{
			"full_name": body.Full_name,
			"email": body.Email,
		},
	})
}

func (h *HandlerUser) UpdateUser(ctx *gin.Context)  {
	userId := ctx.Param("userid")
	var body models.UserModel
	var imageUrl string

	cld, err := helpers.InitCloudinary()
	if  err != nil {
		log.Println(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"Message": "Internal Server Error",
		})
		return
	}
	formFile, _ := ctx.FormFile("user")
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
		
		publicId := fmt.Sprintf("%s_%s-%s", "golang", "user", userId)
		folder := "coffeeshop/user"
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
		// log.Println("here")
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"Message": "Internal Server Error (binding)",
			"error": err.Error(),
		})
		return
	}
	if _, err := govalidator.ValidateStruct(body); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"Message" : "Please Check your Input",
		})
		return
	}
	if body.Full_name == "" && body.Address == "" && body.Phone_number == "" && imageUrl == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Please Input At Least One Change",
		})
		// panic(err)
		return
	}
	log.Println(body)

	result, err := h.RepositoryUpdateUser(&body, userId, imageUrl)
	if  err != nil {
		log.Println(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"Message": "Internal Server Error",
		})
		return
	}
	check, err := result.RowsAffected()
	if  err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"Message": "Internal Server Error",
		})
		return
	}
	if check == 0 {
		ctx.JSON(http.StatusNotFound, gin.H{
			"Message": "User Not Found",
		})
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Update User Success",
		"result": gin.H{
			"profile_pic": imageUrl,
			"address": body.Address,
			"phone_number": body.Phone_number,
		},
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
	message := fmt.Sprintf("Delete on User Id %s Success", userId)

	ctx.JSON(http.StatusOK, gin.H{
		"message": message,
	})
}