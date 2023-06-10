package google

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/Nextasy01/SNS-connections/utils"
	"github.com/gin-gonic/gin"
)

func (gh *GoogleAuthHandler) RevokeAccess(ctx *gin.Context) {
	uid, err := ctx.Cookie("google_id")
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	acc, err := gh.grepo.GetAccById(uid)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	config, err := utils.NewConfig()
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	token := utils.NewToken(acc)

	form := url.Values{}
	form.Set("token", token.RefreshToken)

	resp, err := config.Client(ctx, token).Post("https://oauth2.googleapis.com/revoke", "application/x-www-form-urlencoded", strings.NewReader(form.Encode()))
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		log.Println("Successfully revoked access to Google")

		gh.grepo.DeleteAcc(*acc)

		ctx.Redirect(301, "/view/")
		return
	}
	ctx.AbortWithError(http.StatusInternalServerError, fmt.Errorf("token revokation failed. response status:%s ", resp.Status))

}
