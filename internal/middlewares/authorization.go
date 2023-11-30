package middlewares

import (
	"akbarsyarif/coffeeshopgolang/internal/repositories"
	"akbarsyarif/coffeeshopgolang/pkg"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

func JWTGate(db *sqlx.DB, allowedRole ...string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
	bearerToken := ctx.GetHeader("Authorization")
	if bearerToken == "" {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"message": "Please Login First",
		})
		return
	}
	if !strings.Contains(bearerToken, "Bearer ") {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"message": "Please Login Again",
		})
		return
	}
	token := strings.Replace(bearerToken, "Bearer ", "", -1)
	auth := repositories.InitializeAuthRepository(db)
	res, err := auth.RepositoryCheckToken(token)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "Internal Server Error",
		})
		return
	}
	if len(res) > 0 {
		ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H {
			"message": "Please Login Again",
		})
		return
	}

	// ok, err := checkToken(token)
	// if err != nil {
	// 	ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
	// 		"message": "Internal Server Error",
	// 	})
	// 	return
	// }
	// if len(ok) > 0 {
	// 	ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
	// 		"message": "Please Login First",
	// 	})
	// }

	payload, err := pkg.VerifyToken(token)
	if err != nil {
		if strings.Contains(err.Error(), "expired") {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"Message": "Please Login Again",
			})
			return
		}
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, err)
		return
	}

	allowed := false
	for _, role := range allowedRole {
		if payload.Role == role {
			allowed = true
			break
		}
	}
	if !allowed {
		ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{
			"message": "Access Denied",
		})
	}

	ctx.Set("Payload", payload)
	ctx.Next()
}
}

