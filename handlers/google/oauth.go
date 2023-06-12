package google

import (
	"net/http"
	"time"

	"github.com/Nextasy01/SNS-connections/entity"
	"github.com/Nextasy01/SNS-connections/repository"
	"github.com/Nextasy01/SNS-connections/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	goauth "google.golang.org/api/oauth2/v2"
	"google.golang.org/api/option"
)

var CodeCh = make(chan string)
var UserCh = make(chan uuid.UUID)
var UserPicCh = make(chan string)

type GoogleAuthHandler struct {
	grepo repository.GoogleRepository
}

type GoogleServiceAccount struct {
	Email            string
	ShortAccessToken string
	Lifetime         time.Duration
}

func NewGoogleServiceAccount(email, access_token string) *GoogleServiceAccount {
	return &GoogleServiceAccount{email, access_token, time.Second * 3600}
}

func NewGoogleAuthHandler(gr repository.GoogleRepository) GoogleAuthHandler {
	return GoogleAuthHandler{gr}
}

func (gh *GoogleAuthHandler) CreateAuth(c *gin.Context) {
	g := new(utils.GoogleEnvReader)
	client_id, secret_key, _, _, err := g.ReadFromEnv()

	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	app_env, err := g.GetAppEnv()

	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	redirect_uri := ""

	if app_env == "local" {
		redirect_uri = "http://localhost:9000/view"
	} else {
		redirect_uri = "https://sns-service.onrender.com/view"
	}
	conf := &oauth2.Config{
		ClientID:     client_id,
		ClientSecret: secret_key,
		Scopes:       []string{"email", "profile", "https://www.googleapis.com/auth/youtube", "https://www.googleapis.com/auth/youtube.upload", "https://www.googleapis.com/auth/youtube.readonly"},
		Endpoint:     google.Endpoint,
		RedirectURL:  redirect_uri,
	}

	url := conf.AuthCodeURL("youtube", oauth2.AccessTypeOffline)

	go gh.ExchangeToken(c, CodeCh, UserCh, conf)

	c.Redirect(http.StatusTemporaryRedirect, url)

}

func (gh *GoogleAuthHandler) ExchangeToken(c *gin.Context, code <-chan string, user <-chan uuid.UUID, conf *oauth2.Config) {

	tok, err := conf.Exchange(c, <-code)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "token was not found"})
		return
	}

	oAUTH2service, err := goauth.NewService(c, option.WithTokenSource(conf.TokenSource(c, tok)))
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	userInfo, err := oAUTH2service.Userinfo.Get().Do()
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	acc := entity.NewGoogleAccount()

	if acc.ID, err = uuid.NewRandom(); err != nil {
		panic(err)
	}
	acc.UserID = <-user
	acc.AccessToken = tok.AccessToken
	acc.RefreshToken = tok.RefreshToken
	acc.TokenType = tok.TokenType
	acc.ExpiresAt = tok.Expiry
	acc.Email = userInfo.Email
	acc.Username = userInfo.Name
	acc.ProfilePic = userInfo.Picture

	UserPicCh <- acc.ProfilePic

	gh.grepo.SaveAcc(*acc)

}
