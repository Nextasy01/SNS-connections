package handlers

import (
	"log"
	"net/http"

	"github.com/Nextasy01/SNS-connections/entity"
	"github.com/Nextasy01/SNS-connections/handlers/google"
	"github.com/Nextasy01/SNS-connections/handlers/instagram"
	"github.com/gin-gonic/gin"
)

func GetVideos(ctx *gin.Context, youtubeHandler *google.YouTubeHandler, instagramHandler *instagram.InstagramAuthHandler) {
	uname, err := ctx.Cookie("username")
	if err != nil {
		ctx.AbortWithError(http.StatusNotFound, err)
		log.Println("User not found!?")
		return
	}
	videos_ch := make(chan interface{})
	break_ch := make(chan struct{})

	go func() {
		videos_ch <- youtubeHandler.GetVideos(ctx)
	}()

	go func() {
		videos_ch <- instagramHandler.GetVideos(ctx)
	}()

	go func() {
		<-break_ch
		<-break_ch
		close(videos_ch)
		close(break_ch)
	}()

	var youtube_videos []entity.YoutubeCandidate
	var instagram_videos []entity.InstagramCandidate
	log.Println("Starting loop operation")

loop:
	for p := range videos_ch {
		switch v := p.(type) {
		case []entity.YoutubeCandidate:
			youtube_videos = append(youtube_videos, v...)
		case []entity.InstagramCandidate:
			instagram_videos = append(instagram_videos, v...)
		default:
			break loop
		}
		break_ch <- struct{}{}
	}

	ctx.HTML(http.StatusOK, "index.html", gin.H{
		"youtube_videos": youtube_videos, "instagram_videos": instagram_videos, "username": uname,
	})
}
