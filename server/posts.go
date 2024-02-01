package main

import (
	"fmt"
	"time"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
	"github.com/microcosm-cc/bluemonday"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type PostCommentModel struct {
	gorm.Model

	ParentID  uint   `gorm:"not null"`
	CreatorID uint   `gorm:"not null"`
	Content   string `gorm:"not null"`
}

func (comment *PostCommentModel) GetParent() (PostModel, error) {
	return getPost(int(comment.ParentID))
}

type CategoryModel struct {
	gorm.Model

	Name        string `gorm:"not null,uniqueIndex"`
	Description string `gorm:"not null,default:''"`
}

/*
* getCategories returns all categories in the database.
* It returns an error if the categories could not be fetched.
 */
func getCategories() ([]CategoryModel, error) {
	var categories []CategoryModel
	err := database.Find(&categories).Error
	if err != nil {
		return nil, err
	}

	return categories, nil
}

type PostModel struct {
	gorm.Model

	Title      string `gorm:"not null"`
	Markdown   string `gorm:"not null"`
	UserID     uint   `gorm:"not null"`
	User       UserModel
	CategoryID uint `gorm:"not null"`
}

// Get last active (last comment or post) (i.e 1 hour ago)
func (post *PostModel) GetCreatedAt() string {
	duration := time.Since(post.CreatedAt)
	if duration.Minutes() < 1 {
		return fmt.Sprintf("%.0f seconds ago", duration.Seconds())
	} else if duration.Hours() < 1 {
		return fmt.Sprintf("%.0f minutes ago", duration.Minutes())
	} else if duration.Hours() < 24 {
		return fmt.Sprintf("%.0f hours ago", duration.Hours())
	} else {
		return post.CreatedAt.Format("January 2, 2006")
	}
}

func (comment *PostCommentModel) GetCreator() UserModel {
	var user UserModel
	database.First(&user, comment.CreatorID)

	return user
}

func (comment *PostCommentModel) GetCreatedAt() string {
	duration := time.Since(comment.CreatedAt)
	if duration.Minutes() < 1 {
		return fmt.Sprintf("%.0f seconds ago", duration.Seconds())
	} else if duration.Hours() < 1 {
		return fmt.Sprintf("%.0f minutes ago", duration.Minutes())
	} else if duration.Hours() < 24 {
		return fmt.Sprintf("%.0f hours ago", duration.Hours())
	} else {
		return comment.CreatedAt.Format("January 2, 2006")
	}
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

	unsafeHTML := markdown.ToHTML([]byte(post.Markdown), parser, renderer)

	safeHTML := bluemonday.UGCPolicy().SanitizeBytes(unsafeHTML)

	return string(safeHTML)
}

var postCache = make(map[int]PostModel)

func getPost(id int) (PostModel, error) {
	post, ok := postCache[id]
	if ok {
		return post, nil
	}

	var postModel PostModel
	err := database.Preload(clause.Associations).First(&postModel, id).Error
	if err != nil {
		return PostModel{}, err
	}

	postCache[id] = postModel

	return postModel, nil
}

func getPostComments(id int) ([]PostCommentModel, error) {
	var comments []PostCommentModel
	err := database.Where("parent_id = ?", id).Find(&comments).Error
	if err != nil {
		return nil, err
	}

	return comments, nil
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
