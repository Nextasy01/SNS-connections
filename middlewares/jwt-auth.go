package middlewares

import (
	"net/http"

	"github.com/Nextasy01/SNS-connections/utils"
	"github.com/gin-gonic/gin"
)

func AuthorizeJWT() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		err := utils.TokenValid(ctx)
		if err != nil {
			ctx.Redirect(http.StatusPermanentRedirect, "/login")
			ctx.Abort()
			return
		}
		ctx.Next()
	}
}
