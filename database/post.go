package database

import (
	"errors"
	utils "gforum/utils"

	"github.com/go-playground/validator/v10"
	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
	"github.com/microcosm-cc/bluemonday"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var ErrPostNotFound = errors.New("post not found")

var postValidator = validator.New(validator.WithRequiredStructEnabled())

type PostModel struct {
	gorm.Model

	Title      string `gorm:"not null" validate:"required,min=4,max=64"`
	Markdown   string `gorm:"not null" validate:"required"`
	UserID     uint   `gorm:"not null"`
	User       UserModel
	CategoryID uint `gorm:"not null"`
}

func (post *PostModel) GetCreatedAt() string {
	return post.CreatedAt.Format("2006-01-02 15:04:05")
}

func (post *PostModel) GetTimeAgo() string {
	return utils.TimeAgo(post.CreatedAt)
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

func GetPost(id int) (*PostModel, error) {
	var postModel *PostModel

	err := Database.Preload(clause.Associations).First(&postModel, id).Error
	if err != nil {
		return nil, err
	}

	return postModel, nil
}

func GetRecentPosts(limit int, offset int) ([]PostModel, error) {
	var posts []PostModel

	err := Database.Preload("User").Order("created_at desc").Limit(limit).Offset(offset).Find(&posts).Error
	if err != nil {
		return nil, err
	}

	return posts, nil
}

func GetRecentPostsByCategory(category int, limit int) ([]PostModel, error) {
	var posts []PostModel

	err := Database.Preload("User").Where("category_id = ?", category).Order("created_at desc").Limit(limit).Find(&posts).Error
	if err != nil {
		return nil, err
	}

	return posts, nil
}

type PostCommentModel struct {
	gorm.Model

	ParentID  uint   `gorm:"not null"`
	CreatorID uint   `gorm:"not null"`
	Content   string `gorm:"not null"`
}

func (comment *PostCommentModel) GetCreator() *UserModel {
	var user UserModel
	Database.First(&user, comment.CreatorID)
	return &user
}

func (comment *PostCommentModel) GetCreatedAt() string {
	return comment.CreatedAt.Format("2006-01-02 15:04:05")
}

func (comment *PostCommentModel) GetTimeAgo() string {
	return utils.TimeAgo(comment.CreatedAt)
}

func GetPostComments(id int) ([]PostCommentModel, error) {
	var comments []PostCommentModel

	err := Database.Where("parent_id = ?", id).Find(&comments).Error
	if err != nil {
		return nil, err
	}

	return comments, nil
}

func CreatePost(post *PostModel) error {
	err := postValidator.Struct(post)
	if err != nil {
		return err
	}

	err = Database.Create(post).Error
	if err != nil {
		return err
	}

	return nil
}

func CreatePostComment(comment *PostCommentModel) error {
	err := postValidator.Struct(comment)
	if err != nil {
		return err
	}

	err = Database.Create(comment).Error
	if err != nil {
		return err
	}

	return nil
}
