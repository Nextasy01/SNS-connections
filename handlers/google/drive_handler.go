package google

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
)

func (gh *GoogleAuthHandler) DriveAPICall(c *gin.Context) {
	srv, err := drive.NewService(c, option.WithCredentialsFile("key.json"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	res, err := srv.Files.List().Do()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	var files []string

	for _, file := range res.Files {
		files = append(files, file.Name)
	}

	c.JSON(http.StatusOK, gin.H{"files": files})
}
