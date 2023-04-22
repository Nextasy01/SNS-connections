package handlers

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func ViewDashboard(c *gin.Context) {
	session := sessions.Default(c)

	user := session.Get("user")
	c.HTML(http.StatusOK, "index.html", gin.H{"username": user})
}
