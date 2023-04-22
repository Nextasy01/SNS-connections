package middlewares

import (
	"net/http"

	"github.com/Nextasy01/SNS-connections/utils"
	"github.com/gin-gonic/gin"
)

func AlreadyLoggedIn() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token, err := utils.ExtractTokenID(ctx)
		if token != "" && err == nil {
			ctx.SetCookie("error", "you need to log out first!", 10, "/", ctx.Request.URL.Hostname(), false, true)
			ctx.Redirect(http.StatusTemporaryRedirect, "/")
			ctx.Abort()
			return
		}
		ctx.Next()
	}
}
