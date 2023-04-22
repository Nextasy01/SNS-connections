package middlewares

import (
	"log"

	"github.com/Nextasy01/SNS-connections/routes"
	"github.com/Nextasy01/SNS-connections/utils"
	"github.com/gin-gonic/gin"
)

func CheckInstagramAuth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		user_id, err := utils.ExtractTokenID(ctx)
		if err != nil {
			log.Println(err)
			ctx.Abort()
			return
		}
		instagram_acc, err := routes.DB.GetInstaAccByUserId(user_id)
		if err != nil {
			ctx.Next()
			log.Println(err)
			return
		}
		id, err := ctx.Cookie("instagram_id")
		if err != nil || id != instagram_acc.ID.String() {
			ctx.SetCookie("instagram_id", instagram_acc.ID.String(), 24*3600, "/", ctx.Request.URL.Hostname(), false, true)
			log.Println("Cookie instagram_id was set!")
			ctx.Next()
			return
		}
		ctx.Next()

	}
}
