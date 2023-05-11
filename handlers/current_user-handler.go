package handlers

import (
	"log"
	"net/http"
	"strings"

	"github.com/Nextasy01/SNS-connections/handlers/google"
	"github.com/Nextasy01/SNS-connections/handlers/instagram"
	"github.com/Nextasy01/SNS-connections/repository"
	"github.com/Nextasy01/SNS-connections/utils"
	"github.com/gin-gonic/gin"
)

type CurrentUserHandler struct {
	urepo repository.UserRepository
}

func NewCurrentUserHandler(ur repository.UserRepository) CurrentUserHandler {
	return CurrentUserHandler{ur}
}

func (cuh *CurrentUserHandler) CurrentUser(c *gin.Context) {
	user_id, err := utils.ExtractTokenID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	u, err := cuh.urepo.GetUserByID(user_id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if _, err := c.Cookie("username"); err != nil {
		c.SetCookie("username", u.Username, 24*3600, "/view", c.Request.URL.Hostname(), false, true)
	}

	if err := c.Query("code"); err != "" && strings.Contains(c.Query("state"), "youtube") {
		google.CodeCh <- c.Query("code")
		google.UserCh <- u.ID
		c.HTML(http.StatusAccepted, "index.html", gin.H{"google": "You authenticated your google account",
			"username": u.Username, "isIndex": true})

		close(google.CodeCh)
		close(google.UserCh)

		return
	}

	if err := c.Query("code"); err != "" && strings.Contains(c.Query("state"), "instagram") {
		log.Println("Proceeding with instagram OAuth")
		instagram.InstaCodeCh <- c.Query("code")
		instagram.InstaUserCh <- u.ID

		c.HTML(http.StatusAccepted, "index.html", gin.H{"instagram": "You authenticated your instagram account",
			"username": u.Username, "isIndex": true})
		return
	}

	if msg, err := c.Cookie("error"); err == nil {
		c.HTML(http.StatusOK, "index.html", gin.H{"username": u.Username, "error": msg, "isIndex": true})
		return
	}

	c.HTML(http.StatusOK, "index.html", gin.H{"username": u.Username, "isIndex": true})
}
