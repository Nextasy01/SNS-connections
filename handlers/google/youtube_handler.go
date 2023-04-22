package google

import (
	"log"
	"time"

	"github.com/Nextasy01/SNS-connections/entity"
	"github.com/Nextasy01/SNS-connections/repository"
	"github.com/Nextasy01/SNS-connections/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
)

type YouTubeHandler struct {
	ytrepo repository.YouTubeRepository
	gh     *GoogleAuthHandler
}

func NewYouTubeHandler(ytrepo repository.YouTubeRepository, gh *GoogleAuthHandler) YouTubeHandler {
	return YouTubeHandler{ytrepo: ytrepo, gh: gh}
}

func (yt *YouTubeHandler) GetVideos(c *gin.Context) []entity.YoutubeCandidate {

	// uname, err := c.Cookie("username")
	// if err != nil {
	// 	c.AbortWithError(http.StatusNotFound, err)
	// 	log.Println("User not found!?")
	// 	return nil
	// }

	uid, err := c.Cookie("google_id")
	if err != nil {
		log.Println("Cookie not present?")
		return nil
	}
	acc, err := yt.gh.grepo.GetAccById(uid)
	if err != nil {
		return nil
	}

	config, err := utils.NewConfig()
	if err != nil {
		return nil
	}

	token := utils.NewToken(acc)

	service, err := youtube.NewService(c, option.WithTokenSource(config.TokenSource(c, token)))
	if err != nil {
		return nil
	}

	listCall := service.Channels.List([]string{"contentDetails"})
	listCall = listCall.Mine(true)

	channelList, err := listCall.Do()
	if err != nil {
		log.Println("video retrieving error")
		return nil
	}

	for _, channel := range channelList.Items {
		log.Println(channel.Id)
	}

	playlistResponse := playlistItemsList(service, []string{"snippet"}, channelList.Items[0].ContentDetails.RelatedPlaylists.Uploads)

	if playlistResponse == nil {
		log.Println("video retrieving error")
		return nil
	}

	//c.HTML(http.StatusOK, "index.html", gin.H{"youtube_videos": "true", "username": uname, "videos": yt.SaveNewVideos(playlistResponse.Items, channelList.Items[0].Id, acc.ID)})
	videos, err := yt.SaveNewVideos(playlistResponse.Items, channelList.Items[0].Id, acc.ID)
	if err != nil {
		return nil
	}
	return videos
}

func playlistItemsList(service *youtube.Service, part []string, playlistId string) *youtube.PlaylistItemListResponse {
	call := service.PlaylistItems.List(part)
	call = call.PlaylistId(playlistId)

	response, err := call.Do()
	if err != nil {
		return nil
	}
	return response
}

func (yt *YouTubeHandler) SaveNewVideos(items []*youtube.PlaylistItem, channelId string, accId uuid.UUID) ([]entity.YoutubeCandidate, error) {
	videos := []entity.YoutubeCandidate{}
	videosFromDb, err := yt.ytrepo.GetVideosByAcc(accId.String())
	if err != nil {
		log.Println(err)
		return nil, err
	}
	for i, item := range items {
		videos = append(videos, entity.YoutubeCandidate{
			Title:       item.Snippet.Title,
			Description: item.Snippet.Description,
			VideoId:     item.Snippet.ResourceId.VideoId,
			ChannelId:   channelId,
			CreatorId:   accId,
		})
		if videos[i].ID, err = uuid.NewRandom(); err != nil {
			return nil, err
		}
		if videos[i].PublishedAt, err = time.Parse(time.RFC3339, item.Snippet.PublishedAt); err != nil {
			return nil, err
		}
	}
	newVideos := difference(videos, *videosFromDb)

	if len(newVideos) > 0 {
		yt.ytrepo.SaveVideos(&newVideos)
	}

	return videos, nil
}

func difference(fromAPI, fromDB []entity.YoutubeCandidate) []entity.YoutubeCandidate {
	DBvideos := make(map[string]struct{}, len(fromAPI))
	for _, x := range fromDB {
		DBvideos[x.VideoId] = struct{}{}
	}
	var diff []entity.YoutubeCandidate
	for _, x := range fromAPI {
		if _, ok := DBvideos[x.VideoId]; !ok {
			diff = append(diff, x)
		}
	}
	return diff
}
