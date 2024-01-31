package main

import (
	"html/template"

	"github.com/dchest/captcha"
	"github.com/gin-gonic/gin"
)

type GenericView struct {
	CurrentUser *UserModel
	Branding    Branding
	NumUsers    int64
	NumPosts    int64
	Categories  []CategoryModel
}

func NewGenericView(c *gin.Context) *GenericView {
	var view = GenericView{}

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
	GenericView *GenericView
	Post        PostModel
	Content     template.HTML

	IsRestore   bool
	RestorePath string
}

func renderPostComponent(c *gin.Context) {
	var id int = getParamterInt(c, "id")

	post, err := getPost(id)
	if err != nil {
		renderErrorAlert(c, "invalid post")
		return
	}

	var content = parseHTMLTemplatesFromResources("components/post.html")

	err = content.ExecuteTemplate(c.Writer, "componentBody", PostView{
		Content:     template.HTML(post.ToHTML()),
		GenericView: NewGenericView(c),
		Post:        post,
	})

	if err != nil {
		c.AbortWithError(500, err)
	}

	c.Status(200)
}

type UserRegisterView struct {
	CaptchaID string
}

func renderUserRegisterComponent(c *gin.Context) {
	var content = parseHTMLTemplatesFromResources("components/user/register.html")

	err := content.ExecuteTemplate(c.Writer, "componentBody", UserRegisterView{CaptchaID: captcha.New()})
	if err != nil {
		c.AbortWithError(500, err)
	}

	c.Status(200)
}

func renderUserLoginComponent(c *gin.Context) {
	var content = parseHTMLTemplatesFromResources("components/user/login.html")

	err := content.ExecuteTemplate(c.Writer, "componentBody", nil)
	if err != nil {
		c.AbortWithError(500, err)
	}

	c.Status(200)
}

func renderUserProfileComponent(c *gin.Context) {
	var content = parseTextTemplatesFromResources("components/user/profile.html")

	err := content.ExecuteTemplate(c.Writer, "componentBody", NewGenericView(c))
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
	GenericView    *GenericView
	CategoryModels []CategoryModel
	CaptchaID      string
}

func renderUserComposeComponent(c *gin.Context) {
	categories, err := getCategories()
	if err != nil {
		renderError(c, ErrCategoryNotFound)
		return
	}

	var content = parseHTMLTemplatesFromResources("components/user/compose.html")

	err = content.ExecuteTemplate(c.Writer, "componentBody", ComposeView{GenericView: NewGenericView(c), CategoryModels: categories, CaptchaID: captcha.New()})
	if err != nil {
		c.AbortWithError(500, err)
	}

	c.Status(200)
}

func renderUserPostComposeReplyComponent(c *gin.Context) {
	var id int = getParamterInt(c, "id")

	post, err := getPost(id)
	if err != nil {
		renderError(c, ErrPostNotFound)
		return
	}

	var content = parseHTMLTemplatesFromResources("components/user/composeReply.html")

	err = content.ExecuteTemplate(c.Writer, "componentBody", post)
	if err != nil {
		c.AbortWithError(500, err)
	}

	c.Status(200)
}

type UserPostCommentsFeedView struct {
	GV       *GenericView
	Comments []PostCommentModel
}

func renderUserPostCommentsFeedComponent(c *gin.Context) {
	var id int = getParamterInt(c, "id")

	comments, err := getPostComments(id)
	if err != nil {
		renderError(c, ErrPostNotFound)
		return
	}

	var content = parseHTMLTemplatesFromResources("components/user/postRepliesFeed.html")

	err = content.ExecuteTemplate(c.Writer, "componentBody", UserPostCommentsFeedView{GV: NewGenericView(c), Comments: comments})
	if err != nil {
		c.AbortWithError(500, err)
	}

	c.Status(200)
}

func renderAdministrationMainComponent(c *gin.Context) {
	var content = parseHTMLTemplatesFromResources("components/administration/main.html")

	err := content.ExecuteTemplate(c.Writer, "componentBody", NewGenericView(c))
	if err != nil {
		c.AbortWithError(500, err)
	}

	c.Status(200)
}

type AdministratorEditCategoryComponentView struct {
	GenericView *GenericView
	Category    CategoryModel
}

func renderAdministratorEditCategoryComponent(c *gin.Context) {
	var id int = getParamterInt(c, "id")

	var category CategoryModel
	err := database.Where("id = ?", id).First(&category).Error
	if err != nil {
		renderError(c, ErrCategoryNotFound)
		return
	}

	var content = parseHTMLTemplatesFromResources("components/administration/editCategory.html")

	var view = AdministratorEditCategoryComponentView{
		GenericView: NewGenericView(c),
		Category:    category,
	}

	err = content.ExecuteTemplate(c.Writer, "componentBody", view)
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
