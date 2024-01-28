package main

import (
	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
	"gorm.io/gorm"
)

type PostCommentModel struct {
	gorm.Model

	ParentID  uint `gorm:"not null"`
	CreatorID uint `gorm:"not null"`

	Content string `gorm:"not null"`
}

type CategoryModel struct {
	gorm.Model

	Name        string `gorm:"not null,uniqueIndex"`
	Description string `gorm:"not null,default:'none provided'"`
}

func createCategory(name string, description string) error {
	category := CategoryModel{
		Name:        name,
		Description: description,
	}

	err := database.Create(&category).Error
	if err != nil {
		return err
	}

	return nil
}

type PostModel struct {
	gorm.Model

	Title    string `gorm:"not null"`
	Markdown string `gorm:"not null"`

	UserID uint `gorm:"not null"`
	User   UserModel

	CategoryID uint `gorm:"not null"`
}

func (post *PostModel) GetCategoryName() string {
	var category CategoryModel
	database.First(&category, post.CategoryID)

	return category.Name
}

func (post *PostModel) ToHTML() string {
	extensions := parser.CommonExtensions | parser.AutoHeadingIDs

	parser := parser.NewWithExtensions(extensions)

	htmlFlags := html.CommonFlags | html.HrefTargetBlank
	opts := html.RendererOptions{Flags: htmlFlags}

	renderer := html.NewRenderer(opts)

	return string(markdown.ToHTML([]byte(post.Markdown), parser, renderer))
}

func getRecentPosts(limit int, offset int) ([]PostModel, error) {
	var posts []PostModel
	err := database.Preload("User").Order("created_at desc").Limit(limit).Offset(offset).Find(&posts).Error
	if err != nil {
		return nil, err
	}

	return posts, nil
}

func getRecentPostsByCategory(category int, limit int) ([]PostModel, error) {
	var posts []PostModel
	err := database.Preload("User").Where("category_id = ?", category).Order("created_at desc").Limit(limit).Find(&posts).Error
	if err != nil {
		return nil, err
	}

	return posts, nil
}
