package instagram

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
)

func GetFacebookPageId() {

	params := url.Values{}

	params.Set("access_token", <-tokenCh)

	location := url.URL{Path: "https://graph.facebook.com/v16.0/me/accounts", RawQuery: params.Encode()}

	response, err := http.Get(location.RequestURI())
	if err != nil {
		panic(err)
	}

	responseData, err := io.ReadAll(response.Body)
	if err != nil {
		panic(err)
	}
	log.Println(string(responseData))

	var responseObject struct {
		Data []struct {
			Name  string `json:"name"`
			ID    string `json:"id"`
			Token string `json:"access_token"`
		} `json:"data"`
	}
	if json.Unmarshal(responseData, &responseObject) != nil {
		panic("Unable to process json")
	}

	log.Println(responseObject)
	InstaUserInfo <- responseObject.Data[0].ID
	GetInstagramUserId(responseObject.Data[0].Token, responseObject.Data[0].ID)

}

func GetInstagramUserId(token, facebook_id string) {
	params := url.Values{}

	params.Set("access_token", token)
	params.Set("fields", "instagram_business_account")

	location := url.URL{Path: "https://graph.facebook.com/v16.0/" + facebook_id, RawQuery: params.Encode()}

	response, err := http.Get(location.RequestURI())
	if err != nil {
		panic(err)
	}

	responseData, err := io.ReadAll(response.Body)
	if err != nil {
		panic(err)
	}
	var responseObject struct {
		InstagramBusinessAccount struct {
			ID string `json:"id"`
		} `json:"instagram_business_account"`
	}
	if json.Unmarshal(responseData, &responseObject) != nil {
		panic("Unable to process the json")
	}
	InstaUserInfo <- responseObject.InstagramBusinessAccount.ID
	GetInstagramUserInfo(token, responseObject.InstagramBusinessAccount.ID)
}

func GetInstagramUserInfo(token, instagram_id string) {
	params := url.Values{}

	params.Set("access_token", token)
	params.Set("fields", "biography,id,ig_id,followers_count,follows_count,name,media_count,profile_picture_url,username,website")
	location := url.URL{Path: "https://graph.facebook.com/v16.0/" + instagram_id, RawQuery: params.Encode()}

	response, err := http.Get(location.RequestURI())
	if err != nil {
		panic(err)
	}

	responseData, err := io.ReadAll(response.Body)
	if err != nil {
		panic(err)
	}

	var responseObject struct {
		IgId              int    `json:"ig_id"`
		Username          string `json:"username"`
		ProfilePictureUrl string `json:"profile_picture_url"`
	}
	if json.Unmarshal(responseData, &responseObject) != nil {
		panic("Unable to process the json")
	}

	InstaUserInfo <- fmt.Sprint(responseObject.IgId)
	InstaUserInfo <- responseObject.Username
	InstaUserInfo <- responseObject.ProfilePictureUrl
}
