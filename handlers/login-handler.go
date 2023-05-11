package handlers

import (
	"net/http"
	"net/url"

	"github.com/Nextasy01/SNS-connections/entity"
	"github.com/Nextasy01/SNS-connections/repository"
	"github.com/gin-gonic/gin"
)

type LoginInput struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginHandler struct {
	urepo repository.UserRepository
}

func NewLoginHandler(ur repository.UserRepository) LoginHandler {
	return LoginHandler{ur}
}

func NewLoginInput(username, password string) LoginInput {
	return LoginInput{username, password}
}

func (l *LoginHandler) Login(c *gin.Context) {

	// _, err := c.Cookie("token")
	// if err == nil {
	// 	c.SetCookie("error", "you need to log out first!", 10, "/view", c.Request.URL.Hostname(), false, true)
	// 	c.Abort()
	// 	return
	// }

	input := NewLoginInput(c.PostForm("username"), c.PostForm("password"))

	if err := c.ShouldBind(&input); err != nil {
		c.HTML(http.StatusBadRequest, "login.html", gin.H{"error": err})
		return
	}

	u := entity.NewUser()

	u.Username = input.Username
	u.Password = input.Password

	token, err := l.urepo.LoginCheck(u.Username, u.Password)

	if err != nil {
		c.HTML(http.StatusBadRequest, "login.html", gin.H{"error": "username or password is incorrect"})
		return
	}
	c.SetCookie("token", token, 24*3600, "/", c.Request.URL.Hostname(), false, true)
	c.Header("Cache-Control", "no-cache, private, max-age=0")
	location := url.URL{Path: "/view/"}
	c.Redirect(http.StatusFound, location.RequestURI())
}

func LoginView(c *gin.Context) {
	// _, err := c.Cookie("token")
	// if err == nil {
	// 	c.SetCookie("error", "you need to log out first!", 10, "/view", c.Request.URL.Hostname(), false, true)
	// 	c.Abort()
	// 	return
	// }
	c.HTML(http.StatusOK, "login.html", nil)
}
