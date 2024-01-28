package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"github.com/dchest/captcha"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type HXLocation struct {
	Path   string `json:"path"`
	Target string `json:"target"`
}

func handleUserLogoutAction(c *gin.Context) {
	user, ok := getUserFromContext(c)
	if !ok {
		renderErrorAlert(c, "You are not logged in.")
		return
	}

	_, err := user.IssueSession()
	if err != nil {
		renderErrorAlert(c, err.Error())
		return
	}

	c.SetCookie("session", "", -1, "/", "", false, true)

	renderSuccessAlert(c, "You have been logged out.", true)
}

var userFormValidator = validator.New(validator.WithRequiredStructEnabled())

type UserRegisterActionForm struct {
	Username        string `form:"inputUsername" binding:"required" validate:"required,min=4,max=32"`
	Email           string `form:"inputEmail" binding:"required" validate:"required,email"`
	Password        string `form:"inputPassword" binding:"required" validate:"required,min=8,max=32"`
	CaptchaSolution string `form:"inputCaptchaSolution" binding:"required" validate:"required"`
	CaptchaID       string `form:"inputCaptchaID" binding:"required" validate:"required"`
}

func handleUserRegisterAction(c *gin.Context) {
	var form UserRegisterActionForm
	if err := c.ShouldBind(&form); err != nil {
		renderErrorAlert(c, err.Error())
		return
	}

	err := userFormValidator.Struct(&form)
	if err != nil {
		renderErrorAlert(c, err.Error())
		return
	}

	if !captcha.VerifyString(form.CaptchaID, form.CaptchaSolution) {
		renderErrorAlert(c, "invalid captcha solution.")
		return
	}

	var user = UserModel{
		Username: form.Username,
		Email:    form.Email,
	}

	err = user.SetPassword(form.Password)
	if err != nil {
		renderErrorAlert(c, err.Error())
		return
	}

	err = userFormValidator.Struct(&user)
	if err != nil {
		renderErrorAlert(c, err.Error())
		return
	}

	err = database.Create(&user).Error
	if err != nil {
		renderErrorAlert(c, err.Error())
		return
	}

	token, err := user.IssueSession()
	if err != nil {
		renderErrorAlert(c, err.Error())
		return
	}

	c.SetCookie("session", token, 0, "/", "", false, true)

	c.Header("HX-Refresh", "true")
}

type UserLoginActionForm struct {
	Username string `form:"inputUsername" binding:"required" validate:"required,min=4,max=32"`
	Password string `form:"inputPassword" binding:"required" validate:"required,min=8,max=32"`
}

func handleUserLoginAction(c *gin.Context) {
	var form UserLoginActionForm
	if err := c.ShouldBind(&form); err != nil {
		renderErrorAlert(c, err.Error())
		return
	}

	err := userFormValidator.Struct(&form)
	if err != nil {
		renderErrorAlert(c, err.Error())
		return
	}

	var user UserModel
	err = database.Where("username = ?", form.Username).First(&user).Error
	if err != nil {
		renderErrorAlert(c, "invalid username or password.")
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(form.Password))
	if err != nil {
		renderErrorAlert(c, "invalid username or password.")
		return
	}

	token, err := user.IssueSession()
	if err != nil {
		renderErrorAlert(c, err.Error())
		return
	}

	c.SetCookie("session", token, 0, "/", "", false, true)

	c.Header("HX-Refresh", "true")

	c.Abort()
}

type UserCreatePostActionForm struct {
	Title    string `form:"inputTitle" binding:"required" validate:"required,min=1,max=128"`
	Category string `form:"inputCategory" binding:"required" validate:"required,min=1,max=32"`

	Markdown string `form:"inputMarkdown" binding:"required" validate:"required,min=1,max=1024"`

	CaptchaSolution string `form:"inputCaptchaSolution" binding:"required" validate:"required"`
	CaptchaID       string `form:"inputCaptchaID" binding:"required" validate:"required"`
}

func handleUserComposeAction(c *gin.Context) {
	user, ok := getUserFromContext(c)
	if !ok {
		renderErrorAlert(c, "you are not logged in.")
		return
	}

	var form UserCreatePostActionForm
	if err := c.ShouldBind(&form); err != nil {
		renderErrorAlert(c, err.Error())
		return
	}

	err := userFormValidator.Struct(&form)
	if err != nil {
		renderErrorAlert(c, err.Error())
		return
	}

	if !captcha.VerifyString(form.CaptchaID, form.CaptchaSolution) {
		renderErrorAlert(c, "invalid captcha solution.")
		return
	}

	var category CategoryModel
	err = database.Where("name = ?", form.Category).First(&category).Error
	if err != nil {
		renderErrorAlert(c, "invalid category.")
		return
	}

	var post = PostModel{
		Title:      form.Title,
		Markdown:   form.Markdown,
		CategoryID: category.ID,
		UserID:     user.ID,
	}

	err = database.Create(&post).Error
	if err != nil {
		renderErrorAlert(c, err.Error())
		return
	}

	c.Header("HX-Redirect", fmt.Sprintf("/post/%d", post.ID))
}

const (
	avatarPath      = "resources/avatars/"
	maxFileSize     = 1024 * 1024 * 2 // 2MB
	avatarExtension = ".png"
)

func handleUserUploadAvatarAction(c *gin.Context) {
	user, ok := getUserFromContext(c)
	if !ok {
		renderErrorAlert(c, "you are not logged in.")
		return
	}

	file, err := c.FormFile("inputAvatar")
	if err != nil {
		renderErrorAlert(c, err.Error())
		return
	}

	if file.Size > maxFileSize {
		renderErrorAlert(c, "file is too large.")
		return
	}

	fileName := filepath.Join(avatarPath, generateRandomString(48)+avatarExtension)

	err = c.SaveUploadedFile(file, fileName)
	if err != nil {
		renderErrorAlert(c, err.Error())
		return
	}

	user.Avatar = "\\" + fileName

	err = database.Save(user).Error
	if err != nil {
		renderErrorAlert(c, err.Error())
		return
	}

	renderSuccessAlert(c, "avatar updated.", true)
}

type UserUpdateProfileActionForm struct {
	Username string `form:"inputUsername" binding:"required" validate:"required,alphanum,min=4,max=32"`
	Email    string `form:"inputEmail" binding:"required" validate:"required,email"`
	Password string `form:"inputPassword" binding:"required" validate:"required,min=8,max=32"`
}

func handleUserUpdateProfileAction(c *gin.Context) {
	user, ok := getUserFromContext(c)
	if !ok {
		renderErrorAlert(c, "you are not logged in.")
		return
	}

	var form UserUpdateProfileActionForm
	if err := c.ShouldBind(&form); err != nil {
		renderErrorAlert(c, err.Error())
		return
	}

	user.Email = form.Email
	user.Username = filepath.Clean(form.Username)

	var err error

	if form.Password != "donotchange" {
		err = user.SetPassword(form.Password)
		if err != nil {
			renderErrorAlert(c, err.Error())
			return
		}
	}

	if form.Username != user.Username {
		var count int64
		err = database.Model(&UserModel{}).Where("username = ?", form.Username).Count(&count).Error
		if err != nil {
			renderErrorAlert(c, err.Error())
			return
		}

		if count > 0 {
			renderErrorAlert(c, "username is already taken.")
			return
		}

		err = os.Rename(filepath.Join(avatarPath, user.Username+avatarExtension), filepath.Join(avatarPath, form.Username+avatarExtension))
		if err != nil {
			renderErrorAlert(c, err.Error())
			return
		}

		user.Avatar = form.Username + avatarExtension
	}

	err = userFormValidator.Struct(user)
	if err != nil {
		renderErrorAlert(c, err.Error())
		return
	}

	err = database.Save(user).Error
	if err != nil {
		renderErrorAlert(c, err.Error())
		return
	}

	renderSuccessAlert(c, "profile updated.", true)
}

type ActionModel struct {
	gorm.Model

	Method string `gorm:"not null"`
	Url    string `gorm:"not null"`
	IP     string `gorm:"not null"`

	UserID uint
}

func logActionMiddleware(c *gin.Context) {
	var action = ActionModel{
		Method: c.Request.Method,
		Url:    c.Request.URL.String(),
		IP:     c.ClientIP(),
	}

	user, ok := getUserFromContext(c)
	if ok {
		action.UserID = user.ID
	}

	err := database.Create(&action).Error
	if err != nil {
		c.AbortWithError(500, err)
		return
	}

	c.Next()
}

type AdministrationAddCategoryActionForm struct {
	Name        string `form:"inputName" binding:"required" validate:"required,min=1,max=32"`
	Description string `form:"inputDescription" binding:"required" validate:"required,min=1,max=128"`
}

func handleAdministrationAddCategoryAction(c *gin.Context) {
	var form AdministrationAddCategoryActionForm
	if err := c.ShouldBind(&form); err != nil {
		renderErrorAlert(c, err.Error())
		return
	}

	err := userFormValidator.Struct(&form)
	if err != nil {
		renderErrorAlert(c, err.Error())
		return
	}

	var category = CategoryModel{
		Name:        form.Name,
		Description: form.Description,
	}

	err = database.Create(&category).Error
	if err != nil {
		renderErrorAlert(c, err.Error())
		return
	}

	c.Header("HX-Redirect", "/siteAdministration")
}

func handleAdministrationEditCategoryAction(c *gin.Context) {
	var id int

	if c.Param("id") != "" {
		var err error
		id, err = strconv.Atoi(c.Param("id"))
		if err != nil {
			renderErrorAlert(c, "invalid id")
			return
		}
	}

	var form AdministrationAddCategoryActionForm
	if err := c.ShouldBind(&form); err != nil {
		renderErrorAlert(c, err.Error())
		return
	}

	err := userFormValidator.Struct(&form)
	if err != nil {
		renderErrorAlert(c, err.Error())
		return
	}

	var category CategoryModel
	err = database.Where("id = ?", id).First(&category).Error
	if err != nil {
		renderErrorAlert(c, "invalid category")
		return
	}

	category.Name = form.Name
	category.Description = form.Description

	err = database.Save(&category).Error
	if err != nil {
		renderErrorAlert(c, err.Error())
		return
	}

	c.Header("HX-Redirect", "/siteAdministration")
}

func handleAdministrationDeleteCategoryAction(c *gin.Context) {
	var id int

	if c.Param("id") != "" {
		var err error
		id, err = strconv.Atoi(c.Param("id"))
		if err != nil {
			renderErrorAlert(c, "invalid id")
			return
		}
	}

	var category CategoryModel
	err := database.Where("id = ?", id).First(&category).Error
	if err != nil {
		renderErrorAlert(c, "invalid category")
		return
	}

	err = database.Where("category_id = ?", category.ID).Delete(&PostModel{}).Error
	if err != nil {
		renderErrorAlert(c, err.Error())
		return
	}

	err = database.Delete(&category).Error
	if err != nil {
		renderErrorAlert(c, err.Error())
		return
	}

	c.Header("HX-Redirect", "/siteAdministration")
}
