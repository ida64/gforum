package controller

import (
	database "gforum/database"
	utils "gforum/utils"

	"github.com/gin-gonic/gin"
)

func renderUserIcon(c *gin.Context) {
	var userId = utils.GetParamterInt(c, "id")

	var content = parseTextTemplatesFromResources("components/user/userIcon.html")

	err := content.ExecuteTemplate(c.Writer, "componentBody", database.GetUser(userId))
	if err != nil {
		c.AbortWithError(500, err)
	}

	c.Status(200)
}

type ProfileView struct {
	GlobalView *GlobalView
}

func NewProfileView(c *gin.Context) *ProfileView {
	return &ProfileView{
		GlobalView: NewGlobalView(c),
	}
}

func renderProfile(c *gin.Context) {
	var content = parseHTMLTemplatesFromResources("components/user/profile.html")

	err := content.ExecuteTemplate(c.Writer, "componentBody", NewProfileView(c))
	if err != nil {
		c.AbortWithError(500, err)
	}

	c.Status(200)
}

type ProfileForm struct {
	Username string `form:"inputUsername" binding:"required" validate:"required,min=4,max=32"`
	Email    string `form:"inputEmail" binding:"required" validate:"required,email"`
	Password string `form:"inputPassword" binding:"required" validate:"required,min=8,max=64"`
}

func updateProfileAction(c *gin.Context) {
	user, ok := database.GetUserByValue(c)
	if !ok {
		renderError(c, database.ErrUserNotFound)
		return
	}

	var form ProfileForm
	if err := c.ShouldBind(&form); err != nil {
		renderError(c, err)
		return
	}

	err := userFormValidator.Struct(&form)
	if err != nil {
		renderError(c, err)
		return
	}

	user.Username = form.Username
	user.Email = form.Email

	if form.Password != "donotchange" {
		err = user.SetPassword(form.Password)
		if err != nil {
			renderError(c, err)
			return
		}
	}

	if err := userFormValidator.Struct(user); err != nil {
		renderError(c, err)
		return
	}

	err = database.Database.Save(user).Error
	if err != nil {
		renderError(c, err)
		return
	}

	renderSuccess(c, "profile updated", true)
}

func userMiddleware(c *gin.Context) {
	user, ok := database.GetUserByValue(c)
	if ok {
		c.Set("user", user)
	}

	c.Next()
}

func userEnforceMiddleware(c *gin.Context) {
	_, ok := database.GetUserByValue(c)
	if !ok {
		c.AbortWithStatus(401)
		return
	}

	c.Next()
}

func userAdministratorEnforceMiddleware(c *gin.Context) {
	user, ok := database.GetUserByValue(c)
	if !ok {
		sendErrorWithPage(c, 401, "user not found")
		return
	}

	if !user.IsAdministrator {
		sendErrorWithPage(c, 401, "user is not an administrator")
		return
	}

	c.Next()
}
