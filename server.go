package main

import (
	//"io"
	//"net/http"
	//"net/http"
	"os"

	//"github.com/Nextasy01/SNS-connections/handlers"
	"github.com/Nextasy01/SNS-connections/middlewares"
	//"github.com/Nextasy01/SNS-connections/repository"
	"github.com/Nextasy01/SNS-connections/routes"
	"github.com/gin-gonic/gin"
	//"github.com/jinzhu/gorm"
	//gindump "github.com/tpkeeper/gin-dump"
)

func main() {
	server := gin.Default()
	server.Static("/view/css", "./templates/css")
	server.Static("/css", "./templates/css")
	server.Static("/view/js", "./templates/js")
	server.Static("/view/svg", "./templates/svg")
	server.Static("/svg", "./templates/svg")

	server.LoadHTMLGlob("templates/*.html")

	public := server.Group("/")
	public.Use(middlewares.AlreadyLoggedIn())
	routes.PublicRoutes(public)

	private := server.Group("/view")
	private.Use(middlewares.AuthorizeJWT(), middlewares.CheckInstagramAuth(), middlewares.CheckGoogleAuth())
	routes.PrivateRoutes(private)

	google := private.Group("/google")
	google.Use(middlewares.CheckGoogleAuth())
	routes.GoogleRoutes(google)

	instagram := private.Group("/instagram")
	instagram.Use(middlewares.CheckInstagramAuth())
	routes.InstagramRoutes(instagram)

	port := os.Getenv("PORT")

	if port == "" {
		port = "9000"
	}

	server.Run(":" + port)

}
