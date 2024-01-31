package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/robfig/go-cache"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserModel struct {
	gorm.Model

	Username string `gorm:"unique;not null" validate:"required,alphanum,min=4,max=32"`
	Email    string `gorm:"unique;not null" validate:"required,email"`
	Password string `gorm:"not null"`

	Token string `gorm:"unique"`

	Avatar string

	Posts []PostModel `gorm:"foreignKey:UserID"`

	IsSuspended bool `gorm:"not null,default=false"`

	IsAdministrator bool `gorm:"not null,default=false"`
}

func (user *UserModel) SetPassword(password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user.Password = string(hashedPassword)

	return nil
}

func (user *UserModel) IssueSession() (string, error) {
	user.Token = generateRandomString(64)

	err := database.Save(user).Error
	if err != nil {
		return "", err
	}

	return user.Token, nil
}

func (user *UserModel) GetPosts() ([]PostModel, error) {
	var posts []PostModel
	err := database.Where("UserID = ?", user.ID).Find(&posts).Error
	if err != nil {
		return nil, err
	}

	return posts, nil
}

func getUserFromContext(c *gin.Context) (*UserModel, bool) {
	user, ok := c.Get("user")
	if !ok {
		return nil, false
	}

	return user.(*UserModel), true
}

func sessionToUserMiddleware(c *gin.Context) {
	token, err := c.Cookie("session")
	if err != nil {
		c.Next()
		return
	}

	var user UserModel
	err = database.Where("token = ?", token).First(&user).Error
	if err != nil {
		c.Next()
		return
	}

	c.Set("user", &user)
}

func userRequiredMiddleware(c *gin.Context) {
	user, ok := getUserFromContext(c)
	if !ok {
		c.Redirect(http.StatusFound, "/user/login")
		c.Abort()
		return
	}

	if user.IsSuspended {
		c.Redirect(http.StatusFound, "/user/login")
		c.Abort()
		return
	}

	c.Next()
}

func adminRequiredMiddleware(c *gin.Context) {
	user, ok := getUserFromContext(c)
	if !ok {
		c.Redirect(http.StatusFound, "/user/login")
		c.Abort()
		return
	}

	if !user.IsAdministrator {
		c.Redirect(http.StatusFound, "/user/login")
		c.Abort()
		return
	}

	c.Next()
}

var userCache = cache.New(5*time.Minute, 10*time.Minute)

func getUser(id int) *UserModel {
	var identifier string = fmt.Sprintf("%d", id)

	user, ok := userCache.Get(identifier)
	if ok {
		return user.(*UserModel)
	}

	var userModel UserModel
	database.First(&userModel, id)

	userCache.Set(identifier, &userModel, 5*time.Minute)

	return &userModel
}
