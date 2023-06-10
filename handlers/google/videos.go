package google

import (
	"log"
	"strings"
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
		log.Println(err)
		return nil
	}

	config, err := utils.NewConfig()
	if err != nil {
		log.Println(err)
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

	// testVideos, err := filterShortVideos(channelList.Items[0].Id, service)
	// if err != nil {
	// 	log.Println(err)
	// } else {
	// 	log.Println(testVideos)
	// }

	if playlistResponse == nil {
		log.Println("video retrieving error")
		return nil
	}

	//c.HTML(http.StatusOK, "index.html", gin.H{"youtube_videos": "true", "username": uname, "videos": yt.SaveNewVideos(playlistResponse.Items, channelList.Items[0].Id, acc.ID)})
	videos, err := yt.SaveNewVideos(playlistResponse.Items, channelList.Items[0].Id, acc.ID, service)
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

func (yt *YouTubeHandler) SaveNewVideos(items []*youtube.PlaylistItem, channelId string, accId uuid.UUID, service *youtube.Service) ([]entity.YoutubeCandidate, error) {
	videos := []entity.YoutubeCandidate{}
	videosFromDb, err := yt.ytrepo.GetVideosByAcc(accId.String())
	if err != nil {
		log.Println(err)
		return nil, err
	}
	for i, item := range items {
		call := service.Videos.List([]string{"contentDetails"}).Id(item.Snippet.ResourceId.VideoId)
		resp, err := call.Do()
		if err != nil {
			log.Println("Couldn't get the duration of the video: ", err)
		} else {
			if strings.Contains(resp.Items[0].ContentDetails.Duration, "H") || strings.Contains(resp.Items[0].ContentDetails.Duration, "M") {
				continue
			}
		}
		videos = append(videos, entity.YoutubeCandidate{
			Title:                 item.Snippet.Title,
			Description:           item.Snippet.Description,
			VideoId:               item.Snippet.ResourceId.VideoId,
			IsImported:            false,
			IsImportedToInstagram: false,
			ChannelId:             channelId,
			CreatorId:             accId,
		})
		if videos[i].ID, err = uuid.NewRandom(); err != nil {
			return nil, err
		}
		if videos[i].PublishedAt, err = time.Parse(time.RFC3339, item.Snippet.PublishedAt); err != nil {
			return nil, err
		}
	}

	if len(*videosFromDb) == 0 {
		log.Println("Saving new YouTube Videos")
		yt.ytrepo.SaveVideos(&videos)
		return videos, nil
	}

	checkImports(&videos, videosFromDb)

	log.Println("Filtering New YouTube videos and Database videos")
	newVideos := difference(videos, *videosFromDb)

	if len(newVideos) > 0 {
		yt.ytrepo.SaveVideos(&newVideos)
		videos = append(videos, newVideos...)
		return videos, nil
	}

	return videos, nil
}

func (yt *YouTubeHandler) UpdateVideos(videoId string) error {
	if err := yt.ytrepo.UpdateByYouTubeVideoId(videoId); err != nil {
		return err
	}
	return nil
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

func checkImports(fromAPI, fromDB *[]entity.YoutubeCandidate) {
	DBvideos := make(map[string]bool, len(*fromAPI))
	for _, x := range *fromDB {
		DBvideos[x.VideoId] = x.IsImported
	}
	for i, x := range *fromAPI {
		if _, ok := DBvideos[x.VideoId]; !ok {
			continue
		}
		if DBvideos[x.VideoId] {
			(*fromAPI)[i].IsImported = true
		}
	}
}
