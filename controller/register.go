package controller

import (
	database "gforum/database"

	"github.com/dchest/captcha"
	"github.com/gin-gonic/gin"
)

type RegisterView struct {
	Captcha *Captcha
}

func NewRegisterView() *RegisterView {
	return &RegisterView{Captcha: NewCaptcha()}
}

func renderRegister(c *gin.Context) {
	var content = parseHTMLTemplatesFromResources("components/user/register.html")

	err := content.ExecuteTemplate(c.Writer, "componentBody", NewRegisterView())
	if err != nil {
		renderError(c, err)
		return
	}

	c.Status(200)
}

type RegisterForm struct {
	Username             string `form:"inputUsername" binding:"required" validate:"required,min=4,max=32"`
	Password             string `form:"inputPassword" binding:"required" validate:"required,min=8,max=32"`
	PasswordConfirmation string `form:"inputPasswordConfirmation" binding:"required" validate:"required,eqfield=Password"`
	Email                string `form:"inputEmail" binding:"required" validate:"required,email"`

	CaptchaSolution string `form:"inputCaptchaSolution" binding:"required" validate:"required"`
	CaptchaID       string `form:"inputCaptchaID" binding:"required" validate:"required"`
}

func registerAction(c *gin.Context) {
	var form RegisterForm
	if err := c.ShouldBind(&form); err != nil {
		renderError(c, err)
		return
	}

	err := userFormValidator.Struct(&form)
	if err != nil {
		renderError(c, err)
		return
	}

	if !captcha.VerifyString(form.CaptchaID, form.CaptchaSolution) {
		renderError(c, ErrInvalidCaptchaSolution)
		return
	}

	var user = database.UserModel{
		Username: form.Username,
		Email:    form.Email,
	}

	if err := database.CreateUser(&user, form.Password); err != nil {
		renderError(c, err)
		return
	}

	if session, err := user.IssueSession(); err == nil {
		c.SetCookie("session", session, database.UserSessionCookieExpiration, "/", "", false, true)
	}

	c.Header("HX-Refresh", "true")

	c.Status(200)
}
