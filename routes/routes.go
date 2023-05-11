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
	googleHandler       google.GoogleAuthHandler       = google.NewGoogleAuthHandler(googleRepository)
	youtubeHandler      google.YouTubeHandler          = google.NewYouTubeHandler(youtubeRepository, &googleHandler)
	instagramHandler    instagram.InstagramAuthHandler = instagram.NewInstagramHandler(instagramRepository)
	driverHandler       *handlers.DriveHandler         = handlers.New(&googleHandler, &instagramHandler)
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
		// 	uname, err := ctx.Cookie("username")
		// 	if err != nil {
		// 		ctx.AbortWithError(http.StatusNotFound, err)
		// 		log.Println("User not found!?")
		// 		return
		// 	}
		// 	videos_ch := make(chan interface{})
		// 	break_ch := make(chan struct{})

		// 	go func() {
		// 		videos_ch <- youtubeHandler.GetVideos(ctx)
		// 	}()

		// 	go func() {
		// 		videos_ch <- instagramHandler.GetVideos(ctx)
		// 	}()

		// 	go func() {
		// 		<-break_ch
		// 		<-break_ch
		// 		close(videos_ch)
		// 		close(break_ch)
		// 	}()

		// 	var youtube_videos []entity.YoutubeCandidate
		// 	var instagram_videos []entity.InstagramCandidate
		// 	log.Println("Starting loop operation")

		// loop:
		// 	for p := range videos_ch {
		// 		switch v := p.(type) {
		// 		case []entity.YoutubeCandidate:
		// 			youtube_videos = append(youtube_videos, v...)
		// 		case []entity.InstagramCandidate:
		// 			instagram_videos = append(instagram_videos, v...)
		// 		default:
		// 			break loop
		// 		}
		// 		break_ch <- struct{}{}
		// 	}

		// 	ctx.HTML(http.StatusOK, "index.html", gin.H{
		// 		"youtube_videos": youtube_videos, "instagram_videos": instagram_videos, "username": uname,
		// 	})
	})
}

func GoogleRoutes(g *gin.RouterGroup) {
	g.GET("/authg", googleHandler.CreateAuth)
	g.GET("/drive", driverHandler.DriveAPICall)
	g.POST("/import", driverHandler.Upload)
}

func InstagramRoutes(g *gin.RouterGroup) {
	g.GET("/auth", instagramHandler.CreateAuth)

}
