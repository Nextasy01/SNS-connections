package middlewares

import (
	"log"
	"net/http"

	"github.com/Nextasy01/SNS-connections/utils"
	"github.com/gin-gonic/gin"
)

func AuthorizeJWT() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		log.Println("Trying to verify token")
		err := utils.TokenValid(ctx)
		if err != nil {
			log.Println(err)
			ctx.Redirect(http.StatusTemporaryRedirect, "/login")
			ctx.Abort()
			return
		}
		log.Println("JWT Authorization successful!")
		ctx.Next()
	}
}
