package routes

import (
	"log"
	"net/http"

	"github.com/Nextasy01/SNS-connections/handlers"
	"github.com/Nextasy01/SNS-connections/handlers/google"
	"github.com/Nextasy01/SNS-connections/handlers/instagram"
	"github.com/Nextasy01/SNS-connections/repository"
	"github.com/gin-gonic/gin"
)

var (
	DB                  repository.Database            = repository.NewDatabase()
	userRepository      repository.UserRepository      = repository.NewUserRepository(&DB)
	googleRepository    repository.GoogleRepository    = repository.NewGoogleRepository(&DB)
	youtubeRepository   repository.YouTubeRepository   = repository.NewYouTubeRepository(&DB)
	instagramRepository repository.InstagramRepository = repository.NewInstagramRepository(&DB)
	postRepository      repository.PostRepository      = repository.NewPostRepository(&DB)
	googleHandler       google.GoogleAuthHandler       = google.NewGoogleAuthHandler(googleRepository)
	youtubeHandler      google.YouTubeHandler          = google.NewYouTubeHandler(youtubeRepository, &googleHandler)
	instagramHandler    instagram.InstagramAuthHandler = instagram.NewInstagramHandler(instagramRepository)
	driverHandler       *handlers.DriveHandler         = handlers.New(&youtubeHandler, &instagramHandler, postRepository)
	registerHandler     handlers.RegisterHandler       = handlers.NewRegisterHandler(userRepository)
	loginHandler        handlers.LoginHandler          = handlers.NewLoginHandler(userRepository)
	currentUserHandler  handlers.CurrentUserHandler    = handlers.NewCurrentUserHandler(userRepository)
)

func PublicRoutes(g *gin.RouterGroup) {
	g.GET("/", handlers.LoginView)
	g.GET("/login", handlers.LoginView)
	g.GET("/register", func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "register.html", nil)
	})
	g.POST("/login", loginHandler.Login)
	g.POST("/register", registerHandler.Register)
}

func PrivateRoutes(g *gin.RouterGroup) {
	g.GET("/", currentUserHandler.CurrentUser)
	g.POST("/logout", func(ctx *gin.Context) {
		ctx.SetCookie("token", "", -1, "/", ctx.Request.URL.Hostname(), false, true)
		ctx.SetCookie("instagram_id", "", -1, "/", ctx.Request.URL.Hostname(), false, true)
		ctx.SetCookie("google_id", "", -1, "/", ctx.Request.URL.Hostname(), false, true)
		ctx.Redirect(http.StatusMovedPermanently, "/login")
	})

	g.GET("/view-google", googleHandler.ViewGoogleCredentials)
	//g.GET("youtube-videos", youtubeHandler.GetVideos)
	//g.GET("/instagram-videos", instagramHandler.GetVideos)

	g.GET("/videos", func(ctx *gin.Context) {
		log.Println("Getting Videos from SNS")
		handlers.GetVideos(ctx, &youtubeHandler, &instagramHandler)
	})
}

func GoogleRoutes(g *gin.RouterGroup) {
	g.GET("/authg", googleHandler.CreateAuth)
	g.GET("/drive", driverHandler.DriveAPICall)

	g.POST("/revoke", googleHandler.RevokeAccess)
	g.POST("/import", driverHandler.Upload)
}

func InstagramRoutes(g *gin.RouterGroup) {
	g.GET("/auth", instagramHandler.CreateAuth)

}
