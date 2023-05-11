package middlewares

import (
	"log"
	"net/http"

	"github.com/Nextasy01/SNS-connections/utils"
	"github.com/gin-gonic/gin"
)

func AlreadyLoggedIn() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		log.Println("Initializing AlreadyLoggedIn middleware")
		token, err := utils.ExtractTokenID(ctx)
		if token != "" && err == nil {
			log.Println(ctx.Request.URL.Path)
			log.Println("You need to log out first!")
			ctx.SetCookie("error", "you need to log out first!", 10, "/view", ctx.Request.URL.Hostname(), false, true)
			ctx.Redirect(http.StatusTemporaryRedirect, "/view")
			ctx.Abort()
			return
		}
		log.Println("You are not logged in?")
		ctx.Next()
	}
}
