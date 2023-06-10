package google

import (
	"log"
	"os"

	"github.com/Nextasy01/SNS-connections/utils"
	"github.com/gin-gonic/gin"
	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
)

func (yt *YouTubeHandler) Import(ctx *gin.Context, title string, file *os.File) error {
	uid, err := ctx.Cookie("google_id")
	if err != nil {
		log.Println("Cookie not present?")
		return err
	}
	acc, err := yt.gh.grepo.GetAccById(uid)
	if err != nil {
		log.Println(err)
		return err
	}

	config, err := utils.NewConfig()
	if err != nil {
		log.Println(err)
		return err
	}

	token := utils.NewToken(acc)

	service, err := youtube.NewService(ctx, option.WithTokenSource(config.TokenSource(ctx, token)))
	if err != nil {
		return err
	}

	upload := &youtube.Video{
		Snippet: &youtube.VideoSnippet{
			Title:       title,
			Description: title,
			CategoryId:  "22",
			Tags:        []string{"test", "instagram"},
		},
		Status: &youtube.VideoStatus{PrivacyStatus: "private"},
	}
	log.Println("Preparing to upload video to YouTube")
	uploadCall := service.Videos.Insert([]string{"snippet,status"}, upload)

	log.Println("Starting to upload video to YouTube")
	response, err := uploadCall.Media(file).Do()
	if err != nil {
		return err
	}
	log.Println("Successfully uploaded to YouTube!")
	log.Println(response)
	return nil
}
