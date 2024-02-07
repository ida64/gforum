package controller

import (
	"errors"
	database "gforum/database"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"golang.org/x/crypto/bcrypt"
)

type LoginView struct {
	GlobalView *GlobalView
}

func NewLoginView(c *gin.Context) *LoginView {
	var view = LoginView{
		GlobalView: NewGlobalView(c),
	}

	return &view
}

func renderLogin(c *gin.Context) {
	var content = parseHTMLTemplatesFromResources("components/user/login.html")

	err := content.ExecuteTemplate(c.Writer, "componentBody", NewLoginView(c))
	if err != nil {
		renderError(c, err)
		return
	}

	c.Status(200)
}

var userFormValidator = validator.New(validator.WithRequiredStructEnabled())

var ErrInvalidCredentials = errors.New("invalid username or password")

type LoginForm struct {
	Username string `form:"inputUsername" binding:"required" validate:"required,min=4,max=32"`
	Password string `form:"inputPassword" binding:"required" validate:"required,min=8,max=32"`
}

func loginAction(c *gin.Context) {
	var form LoginForm
	if err := c.ShouldBind(&form); err != nil {
		renderError(c, err)
		return
	}

	err := userFormValidator.Struct(&form)
	if err != nil {
		renderError(c, err)
		return
	}

	user, ok := database.GetUserByUsername(form.Username)
	if !ok {
		renderError(c, ErrInvalidCredentials)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(form.Password))
	if err != nil {
		renderError(c, ErrInvalidCredentials)
		return
	}

	token, err := user.IssueSession()
	if err != nil {
		renderError(c, err)
		return
	}

	c.SetCookie("session", token, 0, "/", "", false, true)

	c.Header("HX-Refresh", "true")

	c.Status(200)
}
