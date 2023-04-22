package google

import (
	"net/http"

	"github.com/Nextasy01/SNS-connections/utils"
	"github.com/gin-gonic/gin"
)

func (gh *GoogleAuthHandler) ViewGoogleCredentials(c *gin.Context) {
	user_id, err := utils.ExtractTokenID(c)
	if err != nil {
		c.Redirect(http.StatusTemporaryRedirect, "/login")
		return
	}
	acc, err := gh.grepo.GetAccByUserId(user_id)
	if err != nil {
		c.HTML(http.StatusForbidden, "index.html", gin.H{"error": "turns out you don't have google account authenticated"})
		return
	}

	c.HTML(http.StatusOK, "google-info.html", gin.H{"Username": acc.Username, "Email": acc.Email, "ProfilePic": acc.ProfilePic})
}
