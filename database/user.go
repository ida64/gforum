package database

import (
	"errors"
	utils "gforum/utils"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var ErrUserNotFound = errors.New("user not found")

var userFormValidator = validator.New(validator.WithRequiredStructEnabled())

const UserSessionCookieExpiration = 60 * 60 * 24 * 30 // 30 days

type UserModel struct {
	gorm.Model

	Username        string `gorm:"unique;not null" validate:"required,alphanum,min=4,max=32"`
	Email           string `gorm:"unique;not null" validate:"required,email"`
	Password        string `gorm:"not null"`
	Token           string `gorm:"unique"`
	Avatar          string
	Posts           []PostModel `gorm:"foreignKey:UserID"`
	IsSuspended     bool        `gorm:"not null,default=false"`
	IsAdministrator bool        `gorm:"not null,default=false"`
}

func (user *UserModel) SetPassword(password string) error {
	// TODO: validate password

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user.Password = string(hashedPassword)

	return nil
}

func (user *UserModel) IssueSession() (string, error) {
	user.Token = utils.GenerateRandomString(64)

	err := Database.Save(user).Error
	if err != nil {
		return "", err
	}

	return user.Token, nil
}

func (user *UserModel) GetPosts() ([]PostModel, error) {
	var posts []PostModel
	err := Database.Where("UserID = ?", user.ID).Find(&posts).Error
	if err != nil {
		return nil, err
	}

	return posts, nil
}

func GetUserByValue(c *gin.Context) (*UserModel, bool) {
	token, err := c.Cookie("session")
	if err != nil {
		return nil, false
	}

	var user UserModel
	Database.Where("token = ?", token).First(&user)

	if user.ID == 0 {
		return nil, false
	}

	return &user, true
}

func GetUser(id int) *UserModel {
	var user UserModel
	Database.First(&user, id)

	return &user
}

func GetUserByUsername(username string) (*UserModel, bool) {
	var user UserModel
	Database.Where("username = ?", username).First(&user)

	if user.ID == 0 {
		return nil, false
	}

	return &user, true
}

func CreateUser(user *UserModel, password string) error {
	if err := user.SetPassword(password); err != nil {
		return err
	}

	if err := userFormValidator.Struct(user); err != nil {
		return err
	}

	err := Database.Create(user).Error
	if err != nil {
		return err
	}

	return nil
}
