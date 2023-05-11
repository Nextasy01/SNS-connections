package handlers

import (
	"net/http"

	"github.com/Nextasy01/SNS-connections/entity"
	"github.com/Nextasy01/SNS-connections/repository"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type RegisterInput struct {
	Email    string `json:"email" binding:"required,email"`
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type RegisterHandler struct {
	urepo repository.UserRepository
}

func NewRegisterHandler(ur repository.UserRepository) RegisterHandler {
	return RegisterHandler{ur}
}

func NewRegisterInput(email, username, password string) RegisterInput {
	return RegisterInput{email, username, password}
}

func (r *RegisterHandler) Register(c *gin.Context) {

	var err error
	// _, err = c.Cookie("token")
	// if err == nil {
	// 	c.SetCookie("error", "you need to log out first!", 10, "/view", c.Request.URL.Hostname(), false, true)
	// 	c.Abort()
	// 	return
	// }

	input := NewRegisterInput(c.PostForm("email"), c.PostForm("username"), c.PostForm("password"))

	if err := c.ShouldBind(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	u := entity.NewUser()

	if u.ID, err = uuid.NewRandom(); err != nil {
		panic(err)
	}
	u.Email = input.Email
	u.Username = input.Username
	u.Password = input.Password

	if err := r.urepo.SaveUser(*u); err != nil {
		c.HTML(http.StatusBadRequest, "register.html", gin.H{"reg_failure": "failed to register!"})
	}

	c.HTML(http.StatusCreated, "register.html", gin.H{"reg_success": "complete!"})

}
