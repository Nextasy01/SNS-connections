package instagram

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/Nextasy01/SNS-connections/entity"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const (
	facebookGraphAPIVersion = "v16.0"
	instagramVideoFields    = "id,caption,comments_count,like_count,is_comment_enabled,is_shared_only,media_product_type,media_type,media_url,permalink,shortcode,thumbnail_url,timestamp"
)

type VideoResponse struct {
	Data []struct {
		ID           string `json:"id"`
		Caption      string `json:"caption"`
		MediaUrl     string `json:"media_url"`
		ThumbnailUrl string `json:"thumbnail_url"`
		MediaType    string `json:"media_type"`
		Permalink    string `json:"permalink"`
		Shortcode    string `json:"shortcode"`
		Timestamp    string `json:"timestamp"`
	} `json:"data"`
}

func (ih *InstagramAuthHandler) GetVideos(c *gin.Context) []entity.InstagramCandidate {
	params := url.Values{}

	uid, err := c.Cookie("instagram_id")
	if err != nil {
		log.Println("Cookie not present?")
		return nil
	}
	acc, err := ih.igrepo.GetInstaAccById(uid)
	if err != nil {
		log.Println(err)
		return nil
	}

	params.Set("access_token", acc.AccessToken)
	params.Set("fields", instagramVideoFields)

	location := url.URL{Path: fmt.Sprintf("https://graph.facebook.com/%s/%s/media", facebookGraphAPIVersion, acc.InstagramBusinessID), RawQuery: params.Encode()}

	response, err := http.Get(location.RequestURI())
	if err != nil {
		log.Println(err)
		return nil
	}

	responseData, err := io.ReadAll(response.Body)
	if err != nil {
		return nil
	}

	var responseObject VideoResponse

	if err = json.Unmarshal(responseData, &responseObject); err != nil {
		return nil
	}
	defer response.Body.Close()
	//c.HTML(http.StatusOK, "index.html", gin.H{"instagram_videos": "true", "username": uname, "videos": responseObject.Data})
	videos, err := ih.SaveNewVideos(responseObject, acc.ID)
	if err != nil {
		return nil
	}
	return videos
}

func (ih *InstagramAuthHandler) UpdateVideos(videoId string) error {
	if err := ih.igrepo.UpdateByVideoId(videoId); err != nil {
		return err
	}
	return nil
}

func (ih *InstagramAuthHandler) SaveNewVideos(items VideoResponse, acc uuid.UUID) ([]entity.InstagramCandidate, error) {
	videos := []entity.InstagramCandidate{}
	videosFromDb, err := ih.igrepo.GetInstaVideosByAcc(acc.String())
	if err != nil {
		log.Println(err)
	}
	log.Println("Processing API videos")
	for i, item := range items.Data {
		if item.MediaType != "VIDEO" {
			continue
		}

		videos = append(videos, entity.InstagramCandidate{
			Caption:             item.Caption,
			MediaUrl:            item.MediaUrl,
			Permalink:           item.Permalink,
			ShortCode:           item.Shortcode,
			VideoId:             item.ID,
			IsImported:          false,
			IsImportedToYouTube: false,
			CreatorId:           acc,
		})
		if videos[i].ID, err = uuid.NewRandom(); err != nil {
			log.Println(err)
			return nil, err
		}
		if videos[i].Timestamp, err = time.Parse("2006-01-02T15:04:05+0000", item.Timestamp); err != nil {
			log.Println(err)
			return nil, err
		}
	}

	if len(*videosFromDb) == 0 {
		log.Println("Saving new Instagram Videos")
		ih.igrepo.SaveInstaVideos(&videos)
		return videos, nil
	}

	checkImports(&videos, videosFromDb)

	log.Println("Filtering New Instagram videos and Database videos")
	newVideos := difference(videos, *videosFromDb)

	if len(newVideos) > 0 {
		ih.igrepo.SaveInstaVideos(&newVideos)
		videos = append(videos, newVideos...)
		return videos, nil
	}
	return videos, nil
}

func difference(fromAPI, fromDB []entity.InstagramCandidate) []entity.InstagramCandidate {
	DBvideos := make(map[string]struct{}, len(fromAPI))
	for _, x := range fromDB {
		DBvideos[x.VideoId] = struct{}{}
	}
	var diff []entity.InstagramCandidate
	for _, x := range fromAPI {
		if _, ok := DBvideos[x.VideoId]; !ok {
			diff = append(diff, x)
		}
	}
	return diff
}

func checkImports(fromAPI, fromDB *[]entity.InstagramCandidate) {
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
