package controller

import (
	database "gforum/database"
	utils "gforum/utils"

	"html/template"

	"github.com/gin-gonic/gin"
)

type PostView struct {
	GlobalView *GlobalView
	Post       *database.PostModel
	Content    template.HTML

	IsRestore   bool
	RestorePath string
}

func NewPostView(post *database.PostModel, globalView *GlobalView) *PostView {
	return &PostView{
		GlobalView: globalView,
		Post:       post,
		Content:    template.HTML(post.ToHTML()),
	}
}

func renderPost(c *gin.Context) {
	var id int = utils.GetParamterInt(c, "id")

	var post, err = database.GetPost(id)
	if err != nil {
		renderError(c, err)
		return
	}

	var content *template.Template = parseHTMLTemplatesFromResources("components/post.html")

	err = content.ExecuteTemplate(c.Writer, "componentBody", NewPostView(post, NewGlobalView(c)))
	if err != nil {
		c.AbortWithError(500, err)
	}

	c.Status(200)
}

type FeedView struct {
	GlobalView      *GlobalView
	FeedName        string
	Posts           []database.PostModel
	CurrentOffset   int
	NextPageOffset  int
	NextPageExists  bool
	LastPageExists  bool
	LastPageOffset  int
	CurrentCategory int
}

func renderFeed(c *gin.Context) {
	offset := utils.GetParamterInt(c, "offset")

	posts, err := database.GetRecentPosts(5, offset)
	if err != nil {
		renderError(c, err)
		return
	}

	view := FeedView{
		GlobalView:      NewGlobalView(c),
		FeedName:        "Recent Posts",
		Posts:           posts,
		CurrentOffset:   offset,
		NextPageOffset:  offset + 1,
		NextPageExists:  len(posts) >= 5,
		LastPageExists:  offset > 0,
		CurrentCategory: 0,
	}

	if view.LastPageExists {
		view.LastPageOffset = offset - 1
	}

	if offset == 0 {
		view.LastPageExists = false
	}

	content := parseHTMLTemplatesFromResources("components/feed.html")

	if err := content.ExecuteTemplate(c.Writer, "componentBody", view); err != nil {
		c.AbortWithError(500, err)
		return
	}

	c.Status(200)
}

func renderFeedByCategory(c *gin.Context) {
	var id int = utils.GetParamterInt(c, "id")

	posts, err := database.GetRecentPostsByCategory(id, 5)
	if err != nil {
		renderError(c, err)
		return
	}

	var feedName string = "Recent Posts"
	if len(posts) > 0 {
		if categoryName, err := database.GetCategoryName(id); err == nil {
			feedName = categoryName
		}
	}

	var content = parseHTMLTemplatesFromResources("components/feed.html")

	if err := content.ExecuteTemplate(c.Writer, "componentBody", FeedView{
		GlobalView:      NewGlobalView(c),
		FeedName:        feedName,
		Posts:           posts,
		CurrentCategory: id,
	}); err != nil {
		c.AbortWithError(500, err)
		return
	}

	c.Status(200)
}

type PostFeedView struct {
	Comments     []database.PostCommentModel
	OriginalPost *database.PostModel
}

func NewPostFeedView(c *gin.Context) *PostFeedView {
	var id int = utils.GetParamterInt(c, "post_id")

	comments, err := database.GetPostComments(id)
	if err != nil {
		renderError(c, database.ErrPostNotFound)
		return nil
	}

	post, err := database.GetPost(id)
	if err != nil {
		renderError(c, database.ErrPostNotFound)
		return nil
	}

	return &PostFeedView{Comments: comments, OriginalPost: post}
}

func renderPostFeed(c *gin.Context) {
	var content = parseHTMLTemplatesFromResources("components/user/commentFeed.html")

	err := content.ExecuteTemplate(c.Writer, "componentBody", NewPostFeedView(c))
	if err != nil {
		c.AbortWithError(500, err)
	}

	c.Status(200)
}

type CommentForm struct {
	Content string `form:"inputContent" binding:"required" validate:"required,min=4,max=1024"`
}

func addCommentAction(c *gin.Context) {
	var form CommentForm
	if err := c.ShouldBind(&form); err != nil {
		renderError(c, err)
		return
	}

	if err := userFormValidator.Struct(&form); err != nil {
		renderError(c, err)
		return
	}

	var id int = utils.GetParamterInt(c, "id")

	if _, err := database.GetPost(id); err != nil {
		renderError(c, err)
		return
	}

	user, ok := database.GetUserByValue(c)
	if !ok {
		renderError(c, database.ErrUserNotFound)
		return
	}

	comment := database.PostCommentModel{
		ParentID:  uint(id),
		CreatorID: user.ID,
		Content:   form.Content,
	}

	if err := database.CreatePostComment(&comment); err != nil {
		renderError(c, err)
		return
	}

	renderPost(c)
}
