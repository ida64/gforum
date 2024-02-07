package controller

import (
	"fmt"
	"gforum/database"

	"github.com/dchest/captcha"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

var composeFormValidator = validator.New(validator.WithRequiredStructEnabled())

type ComposeView struct {
	GlobalView *GlobalView
	Captcha    *Captcha
}

func NewComposeView(c *gin.Context) *ComposeView {
	return &ComposeView{
		GlobalView: NewGlobalView(c),
		Captcha:    NewCaptcha(),
	}
}

func renderCompose(c *gin.Context) {
	content := parseHTMLTemplatesFromResources("components/user/compose.html")

	if err := content.ExecuteTemplate(c.Writer, "componentBody", NewComposeView(c)); err != nil {
		c.AbortWithError(500, err)
		return
	}

	c.Status(200)
}

type ComposeForm struct {
	Title    string `form:"inputTitle" binding:"required" validate:"required,min=4,max=64"`
	Category string `form:"inputCategory" binding:"required" validate:"required"`

	Markdown string `form:"inputMarkdown" binding:"required" validate:"required"`

	CaptchaSolution string `form:"inputCaptchaSolution" binding:"required" validate:"required"`
	CaptchaID       string `form:"inputCaptchaID" binding:"required" validate:"required"`
}

func composeAction(c *gin.Context) {
	var form ComposeForm
	if err := c.ShouldBind(&form); err != nil {
		renderError(c, err)
		return
	}

	err := composeFormValidator.Struct(&form)
	if err != nil {
		renderError(c, err)
		return
	}

	if !captcha.VerifyString(form.CaptchaID, form.CaptchaSolution) {
		renderError(c, ErrInvalidCaptchaSolution)
		return
	}

	category, err := database.GetCategoryByName(form.Category)
	if err != nil {
		renderError(c, err)
		return
	}

	user, ok := database.GetUserByValue(c)
	if !ok {
		renderError(c, database.ErrUserNotFound)
		return
	}

	post := database.PostModel{
		Title:      form.Title,
		Markdown:   form.Markdown,
		UserID:     user.ID,
		User:       *user, // TODO(paging): is there a flag to skip a field?
		CategoryID: category.ID,
	}

	if err := database.CreatePost(&post); err != nil {
		renderError(c, err)
		return
	}

	c.Header("HX-Redirect", fmt.Sprintf("/post/%d", post.ID))

	c.Status(200)
}
