package instagram

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/Nextasy01/SNS-connections/entity"
	"github.com/Nextasy01/SNS-connections/repository"
	"github.com/Nextasy01/SNS-connections/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

var (
	InstaCodeCh   = make(chan string)
	InstaUserCh   = make(chan uuid.UUID)
	InstaUserInfo = make(chan string)
	tokenCh       = make(chan string)
	redirect_uri  = ""
)

type InstagramAuthHandler struct {
	igrepo repository.InstagramRepository
}

type Response struct {
	Token     string `json:"access_token"`
	TokenType string `json:"token_type"`
	ExpiresAt int    `json:"expires_in"`
}

func NewInstagramHandler(ig repository.InstagramRepository) InstagramAuthHandler {
	return InstagramAuthHandler{ig}
}

func (ih *InstagramAuthHandler) CreateAuth(c *gin.Context) {
	f := new(utils.FacebookEnvReader)
	app_id, app_secret, err := f.ReadFromEnv()
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	app_env, err := f.GetAppEnv()

	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	if app_env == "local" {
		redirect_uri = "http://localhost:9000/view/"
	} else {
		redirect_uri = "https://sns-service.onrender.com/view/"
	}

	params := url.Values{}

	params.Set("client_id", app_id)
	params.Set("redirect_uri", redirect_uri)
	params.Set("scope", "instagram_basic,pages_show_list,pages_read_engagement,instagram_manage_insights,instagram_content_publish")
	params.Set("state", "instagram")

	location := url.URL{Path: "https://www.facebook.com/v16.0/dialog/oauth", RawQuery: params.Encode()}

	go ih.ExchangeToken(c, InstaCodeCh, InstaUserCh, app_id, app_secret)
	go GetFacebookPageId()

	c.Redirect(http.StatusTemporaryRedirect, location.RequestURI())
}

func (ih *InstagramAuthHandler) ExchangeToken(c *gin.Context, code <-chan string, user <-chan uuid.UUID, app_id, app_secret string) {
	defer close(InstaUserInfo)
	defer close(InstaCodeCh)
	defer close(InstaUserCh)

	params := url.Values{}
	temp := <-code
	log.Println(temp)
	params.Set("client_id", app_id)
	params.Set("redirect_uri", redirect_uri)
	params.Set("client_secret", app_secret)
	params.Set("grant_type", "authorization_code")
	params.Set("code", temp)

	location := url.URL{Path: "https://graph.facebook.com/v16.0/oauth/access_token", RawQuery: params.Encode()}

	response, err := http.Get(location.RequestURI())
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
	}

	responseData, err := io.ReadAll(response.Body)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
	}
	log.Println(string(responseData))
	var responseObject Response
	if err = json.Unmarshal(responseData, &responseObject); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	log.Println(responseObject.Token)
	tokenCh <- responseObject.Token

	acc := entity.NewInstagramAccount()

	if acc.ID, err = uuid.NewRandom(); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	acc.UserID = <-InstaUserCh
	acc.InstagramPrivateID = <-InstaUserInfo
	acc.InstagramBusinessID = <-InstaUserInfo
	acc.InstagramUserID = <-InstaUserInfo
	acc.Username = <-InstaUserInfo
	acc.ProfilePic = <-InstaUserInfo
	acc.AccessToken = responseObject.Token
	acc.TokenType = responseObject.TokenType
	acc.ExpiresAt = time.Now().Add(time.Second * time.Duration(responseObject.ExpiresAt))

	ih.igrepo.SaveInstaAcc(*acc)
}
