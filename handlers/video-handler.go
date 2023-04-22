package handlers

import (
	"net/http"

	"github.com/Nextasy01/SNS-connections/entity"
	"github.com/gin-gonic/gin"
)

func GetVideos(c *gin.Context) {
	data := gin.H{
		"title":  "Login Page",
		"videos": entity.Candidates,
	}
	c.HTML(http.StatusOK, "index.html", data)
}
