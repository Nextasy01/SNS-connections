package instagram

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/gin-gonic/gin"
)

type responseObject struct {
	ID string `json:"id"`
}

type responseStatusObject struct {
	Status     string `json:"status"`
	StatusCode string `json:"status_code"`
	ID         string `json:"id"`
}

func (ih *InstagramAuthHandler) Import(ctx *gin.Context, link, caption string) {
	params := url.Values{}

	uid, err := ctx.Cookie("instagram_id")
	if err != nil {
		log.Println("Instagram Id cookie is not present!", err)
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	acc, err := ih.igrepo.GetInstaAccById(uid)
	if err != nil {
		log.Println(err)
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	params.Set("access_token", acc.AccessToken)
	params.Set("media_type", "REELS")
	params.Set("video_url", link)
	params.Set("caption", caption)
	params.Set("thumb_offset", "0")
	params.Set("share_to_feed", "true")

	location := url.URL{Path: fmt.Sprintf("https://graph.facebook.com/%s/%s/media", facebookGraphAPIVersion, acc.InstagramBusinessID), RawQuery: params.Encode()}

	response, err := http.Post(location.RequestURI(), "application/json", nil)
	if err != nil {
		log.Print(err)
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	responseData, err := io.ReadAll(response.Body)
	if err != nil {
		log.Print(err)
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	var responseObj responseObject

	log.Println("Unmarshalling JSON output")
	err = json.Unmarshal(responseData, &responseObj)
	if err != nil {
		log.Print(err)
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	log.Println(string(responseData))

	for {
		status, err := ih.checkStatus(responseObj.ID, acc.AccessToken)
		if err != nil {
			log.Print(err)
			ctx.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		if status == "FINISHED" {
			break
		}
		if status == "IN_PROGRESS" {
			time.Sleep(3 * time.Second)
		} else {
			log.Println("Status: " + status)
			break
		}
	}

	err = ih.publish(responseObj.ID, acc.AccessToken, acc.InstagramBusinessID)
	if err != nil {
		log.Print(err)
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

}

func (ih *InstagramAuthHandler) publish(creation_id, access_token, business_id string) error {
	params := url.Values{}
	params.Set("access_token", access_token)
	params.Set("creation_id", creation_id)

	location := url.URL{Path: fmt.Sprintf("https://graph.facebook.com/%s/%s/media_publish", facebookGraphAPIVersion, business_id), RawQuery: params.Encode()}

	response, err := http.Post(location.RequestURI(), "application/json", nil)
	if err != nil {
		return err
	}

	responseData, err := io.ReadAll(response.Body)
	if err != nil {
		return err
	}

	log.Println(string(responseData))
	return nil
}

func (ih *InstagramAuthHandler) checkStatus(creation_id, access_token string) (string, error) {
	var responseObjStatus responseStatusObject

	params := url.Values{}
	params.Set("access_token", access_token)
	params.Set("fields", "status,status_code")

	location := url.URL{Path: fmt.Sprintf("https://graph.facebook.com/%s/%s", facebookGraphAPIVersion, creation_id), RawQuery: params.Encode()}

	response, err := http.Get(location.RequestURI())
	if err != nil {
		return "", err
	}

	responseData, err := io.ReadAll(response.Body)
	if err != nil {
		return "", err
	}

	err = json.Unmarshal(responseData, &responseObjStatus)
	if err != nil {
		return "", err
	}

	return responseObjStatus.StatusCode, nil
}
