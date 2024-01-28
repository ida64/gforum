package main

import (
	"crypto/rand"
	"fmt"
	"net/http"

	"math/big"
	"strings"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func generateRandomString(length int) string {
	const alphanumeric = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	var builder strings.Builder
	builder.Grow(length)

	for i := 0; i < length; i++ {
		randomIndex, _ := rand.Int(rand.Reader, big.NewInt(int64(len(alphanumeric))))
		builder.WriteByte(alphanumeric[randomIndex.Int64()])
	}

	return builder.String()
}

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
		fmt.Println(err)
		c.Next()
		return
	}

	c.Set("user", &user)
}

func userRequiredMiddleware(c *gin.Context) {
	user, ok := getUserFromContext(c)
	if !ok {
		c.Redirect(http.StatusFound, "/login")
		c.Abort()
		return
	}

	if user.IsSuspended {
		c.Redirect(http.StatusFound, "/login")
		c.Abort()
		return
	}

	c.Next()
}

func adminRequiredMiddleware(c *gin.Context) {
	user, ok := getUserFromContext(c)
	if !ok {
		c.Redirect(http.StatusFound, "/login")
		c.Abort()
		return
	}

	if !user.IsAdministrator {
		c.Redirect(http.StatusFound, "/login")
		c.Abort()
		return
	}

	c.Next()
}
