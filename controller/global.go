package controller

import (
	"gforum/config"
	database "gforum/database"

	"github.com/gin-gonic/gin"
)

type GlobalView struct {
	CurrentUser *database.UserModel
	Branding    config.Branding
	NumUsers    int64
	NumPosts    int64
	Categories  []database.CategoryModel
}

func NewGlobalView(c *gin.Context) *GlobalView {
	var view = GlobalView{}

	view.Branding = config.LoadedConfig.Branding

	view.CurrentUser, _ = database.GetUserByValue(c)

	categories, err := database.GetCategories()
	if err != nil {
		return nil
	}

	view.NumUsers = int64(len(categories))
	view.NumPosts = int64(len(categories))

	view.Categories = categories

	return &view
}
