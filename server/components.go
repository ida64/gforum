package main

import (
	"html/template"

	"github.com/dchest/captcha"
	"github.com/gin-gonic/gin"
)

type GlobalView struct {
	CurrentUser *UserModel
	Branding    Branding
	NumUsers    int64
	NumPosts    int64
	Categories  []CategoryModel
}

func NewGlobalView(c *gin.Context) *GlobalView {
	var view = GlobalView{}

	user, ok := c.Get("user")
	if ok {
		view.CurrentUser = user.(*UserModel)
	}

	view.Branding = loadedConfig.Branding

	categories, err := getCategories()
	if err != nil {
		renderErrorAlert(c, "error fetching categories")
		return nil
	}

	view.NumUsers = int64(len(categories))
	view.NumPosts = int64(len(categories))

	view.Categories = categories

	return &view
}

type PostView struct {
	GlobalView *GlobalView
	Post       PostModel
	Content    template.HTML

	IsRestore   bool
	RestorePath string
}

func NewPostView(c *gin.Context) *PostView {
	var id int = getParamterInt(c, "id")

	post, err := getPost(id)
	if err != nil {
		renderError(c, ErrPostNotFound)
		return nil
	}

	return &PostView{
		GlobalView: NewGlobalView(c),
		Post:       post,
		Content:    template.HTML(post.ToHTML()),
	}
}

func renderPostComponent(c *gin.Context) {
	var content = parseHTMLTemplatesFromResources("components/post.html")

	err := content.ExecuteTemplate(c.Writer, "componentBody", NewPostView(c))
	if err != nil {
		c.AbortWithError(500, err)
	}

	c.Status(200)
}

type UserRegisterView struct {
	CaptchaID string
}

func NewUserRegisterView() *UserRegisterView {
	return &UserRegisterView{CaptchaID: captcha.New()}
}

func renderUserRegisterComponent(c *gin.Context) {
	var content = parseHTMLTemplatesFromResources("components/user/register.html")

	err := content.ExecuteTemplate(c.Writer, "componentBody", NewUserRegisterView())
	if err != nil {
		c.AbortWithError(500, err)
	}

	c.Status(200)
}

func renderUserLoginComponent(c *gin.Context) {
	var content = parseHTMLTemplatesFromResources("components/user/login.html")

	err := content.ExecuteTemplate(c.Writer, "componentBody", NewGlobalView(c))
	if err != nil {
		c.AbortWithError(500, err)
	}

	c.Status(200)
}

func renderUserProfileComponent(c *gin.Context) {
	var content = parseTextTemplatesFromResources("components/user/profile.html")

	err := content.ExecuteTemplate(c.Writer, "componentBody", NewGlobalView(c))
	if err != nil {
		c.AbortWithError(500, err)
	}

	c.Status(200)
}

func renderUserAvatarImageComponent(c *gin.Context) {
	var userId = getParamterInt(c, "id")

	var content = parseTextTemplatesFromResources("components/user/userIcon.html")

	err := content.ExecuteTemplate(c.Writer, "componentBody", getUser(userId))
	if err != nil {
		c.AbortWithError(500, err)
	}

	c.Status(200)
}

type ComposeView struct {
	GlobalView     *GlobalView
	CategoryModels []CategoryModel
	CaptchaID      string
}

func NewComposeView(c *gin.Context) *ComposeView {
	categories, err := getCategories()
	if err != nil {
		renderError(c, ErrCategoryNotFound)
		return nil
	}

	return &ComposeView{
		GlobalView:     NewGlobalView(c),
		CategoryModels: categories,
		CaptchaID:      captcha.New(),
	}
}

func renderUserComposeComponent(c *gin.Context) {
	var content = parseHTMLTemplatesFromResources("components/user/compose.html")

	err := content.ExecuteTemplate(c.Writer, "componentBody", NewComposeView(c))
	if err != nil {
		c.AbortWithError(500, err)
	}

	c.Status(200)
}

type PostFeedView struct {
	Comments     []PostCommentModel
	OriginalPost PostModel
}

func NewPostFeedView(c *gin.Context) *PostFeedView {
	var id int = getParamterInt(c, "id")

	comments, err := getPostComments(id)
	if err != nil {
		renderError(c, ErrPostNotFound)
		return nil
	}

	post, err := getPost(id)
	if err != nil {
		renderError(c, ErrPostNotFound)
		return nil
	}

	return &PostFeedView{Comments: comments, OriginalPost: post}
}

func renderPostFeedComponent(c *gin.Context) {
	var content = parseHTMLTemplatesFromResources("components/user/commentFeed.html")

	err := content.ExecuteTemplate(c.Writer, "componentBody", NewPostFeedView(c))
	if err != nil {
		c.AbortWithError(500, err)
	}

	c.Status(200)
}

func renderAdministrationMainComponent(c *gin.Context) {
	var content = parseHTMLTemplatesFromResources("components/administration/main.html")

	err := content.ExecuteTemplate(c.Writer, "componentBody", NewGlobalView(c))
	if err != nil {
		c.AbortWithError(500, err)
	}

	c.Status(200)
}

type EditCategoryView struct {
	Category CategoryModel
}

func NewEditCategoryView(c *gin.Context) *EditCategoryView {
	var id int = getParamterInt(c, "id")

	var category CategoryModel
	err := database.Where("id = ?", id).First(&category).Error
	if err != nil {
		renderError(c, ErrCategoryNotFound)
		return nil
	}

	return &EditCategoryView{
		Category: category,
	}
}

func renderAdministratorEditCategoryComponent(c *gin.Context) {
	var content = parseHTMLTemplatesFromResources("components/administration/editCategory.html")

	err := content.ExecuteTemplate(c.Writer, "componentBody", NewEditCategoryView(c))
	if err != nil {
		c.AbortWithError(500, err)
	}

	c.Status(200)
}

func renderRulesComponent(c *gin.Context) {
	var content = parseHTMLTemplatesFromResources("components/rules.html")

	err := content.ExecuteTemplate(c.Writer, "componentBody", nil)
	if err != nil {
		c.AbortWithError(500, err)
	}

	c.Status(200)
}
