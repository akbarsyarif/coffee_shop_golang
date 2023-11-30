package handlers

import (
	"akbarsyarif/coffeeshopgolang/internal/models"
	"akbarsyarif/coffeeshopgolang/internal/repositories"
	"akbarsyarif/coffeeshopgolang/pkg"
	"log"
	"net/http"
	"strings"

	"github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
)

type HandlerAuth struct {
	*repositories.AuthRepository
}

func InitializeAuthHandler(r *repositories.AuthRepository) *HandlerAuth {
	return &HandlerAuth{r}
}

func (h *HandlerAuth) Register(ctx *gin.Context)  {
	var body models.UserModel
	var userId string
	if err := ctx.ShouldBind(&body); err != nil {
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}

	if _ , err := govalidator.ValidateStruct(body); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"Message" : "All Inputs Must Be Filled",
		})
		return
	}

	c := pkg.HashConfig{
		Time:    3,
		Memory:  64 * 1024,
		Threads: 2,
		KeyLen:  32,
		SaltLen: 16,
	}
	hashedPassword, err := c.GenHashedPassword(body.Pwd)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Internal Server Error",
		})
		return
	}

	tx, err := h.Beginx()
	if err != nil {
		log.Println(err)
		ctx.JSON(http.StatusInternalServerError, err)
		// panic(err)
		return
	}
	defer tx.Rollback()

	result, err := h.RepositoryRegister(&body, hashedPassword, tx)
	// log.Println(result)
	if err != nil {
		pgErr, _ := err.(*pq.Error)
		if pgErr.Code == "23505" {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": "Email is Registered",
			})
			return
		}
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}
	for result.Next() {
		var id string
		err = result.Scan(&id)
		if id != "" {
			log.Println(id)
			userId = id
			break
		}
	}
	if err != nil {
		log.Println(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Internal Server Error",
		})
		// panic(err)
		return
	}
	result.Rows.Close()

	err = h.RepositoryCreateUserProfile(userId, tx)
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
		"message": "Register Success",
		"result": gin.H{
			"full_name": body.Full_name,
			"email": body.Email,
		},
	})
}

func (h *HandlerAuth) Login(ctx *gin.Context) {
	body := &models.GetUserInfoModel{}
	if err := ctx.ShouldBind(body); err != nil {
		log.Println(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"Message": "Internal Server Error",
		})
		return
	}

	if _, err := govalidator.ValidateStruct(body); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"Message" : "Email And Password Cannot Be Empty",
		})
		return
	}

	result, err := h.RepositoryGetPassword(body)
	if err != nil {
		log.Println(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"Message": "Internal Server Error",
		})
		return
	}
	if len(result) == 0 {
		ctx.JSON(http.StatusNotFound, gin.H{
			"message": "Email or Password is Wrong",
		})
		return
	}

	hc := pkg.HashConfig{}
	isValid, err := hc.ComparePassword(body.Pwd, result[0].Pwd)
	if err != nil {
		log.Println(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"Message": "Internal Server Error",
		})
		return
	}
	if !isValid {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"message": "Email or Password is Wrong",
		})
		return
	}
	// log.Println(result[0].Isverified)
	if result[0].Isverified == false {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"message": "Please Verify Your Email First",
		})
		return
	}

	payload := pkg.NewPayload(result[0].Id, result[0].Role_name)
	token, err := payload.GenerateToken()
	if err != nil {
		log.Println(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"Message": "Internal Server Error",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Login Success",
		"result": gin.H{
			"token": token,
			"user_info": gin.H{
				"id": result[0].Id,
				"email": body.Email,
				"full_name": result[0].Full_name,
			},
		},
	})
}

func (h *HandlerAuth) Logout(ctx *gin.Context)  {
	bearerToken := ctx.GetHeader("Authorization")
	token := strings.Split(bearerToken, " ")[1]
	// log.Println(token)
	if err := h.RepositoryLogout(token); err != nil {
		log.Println(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"Message": "Internal Server Error",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Logout Success",
	})
}