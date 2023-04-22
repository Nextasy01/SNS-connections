package middlewares

import (
	"log"

	"github.com/Nextasy01/SNS-connections/routes"
	"github.com/Nextasy01/SNS-connections/utils"
	"github.com/gin-gonic/gin"
)

func CheckGoogleAuth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		user_id, err := utils.ExtractTokenID(ctx)
		if err != nil {
			log.Println(err)
			ctx.Next()
			return
		}
		google_acc, err := routes.DB.GetAccByUserId(user_id)
		if err != nil {
			log.Println(err)
			ctx.Next()
			return
		}
		id, err := ctx.Cookie("google_id")
		if err != nil || id != google_acc.ID.String() {
			ctx.SetCookie("google_id", google_acc.ID.String(), 24*3600, "/", ctx.Request.URL.Hostname(), false, true)
			log.Println("Cookie google_id was set!")
			ctx.Next()
			return
		}
		ctx.Next()

	}
}
