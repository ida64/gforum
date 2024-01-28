package main

import (
	"html/template"
	"strconv"

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

	var categories []CategoryModel
	database.Find(&categories)

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
	postID := c.Param("id")

	var post PostModel
	err := database.Preload("User").Where("id = ?", postID).First(&post).Error
	if err != nil {
		renderErrorAlert(c, "Error fetching post")
		return
	}

	var content = parseTmplFromResources("components/post.html")

	err = content.ExecuteTemplate(c.Writer, "componentBody", PostView{Content: template.HTML(post.ToHTML()), GenericView: NewGenericView(c), Post: post})
	if err != nil {
		c.AbortWithError(500, err)
	}

	c.Status(200)
}

type UserRegisterView struct {
	CaptchaID string
}

func renderUserRegisterComponent(c *gin.Context) {
	var content = parseTmplFromResources("components/user/register.html")

	err := content.ExecuteTemplate(c.Writer, "componentBody", UserRegisterView{CaptchaID: captcha.New()})
	if err != nil {
		c.AbortWithError(500, err)
	}

	c.Status(200)
}

func renderUserLoginComponent(c *gin.Context) {
	var content = parseTmplFromResources("components/user/login.html")

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
	var categories []CategoryModel
	err := database.Find(&categories).Error
	if err != nil {
		renderErrorAlert(c, "Error fetching categories")
		return
	}

	var content = parseTmplFromResources("components/user/compose.html")

	err = content.ExecuteTemplate(c.Writer, "componentBody", ComposeView{GenericView: NewGenericView(c), CategoryModels: categories, CaptchaID: captcha.New()})
	if err != nil {
		c.AbortWithError(500, err)
	}

	c.Status(200)
}

func renderUserComposeReplyComponent(c *gin.Context) {
	var id int

	if c.Param("id") != "" {
		var err error
		id, err = strconv.Atoi(c.Param("id"))
		if err != nil {
			renderErrorAlert(c, "invalid id")
			return
		}
	}

	var post PostModel
	err := database.Where("id = ?", id).First(&post).Error
	if err != nil {
		renderErrorAlert(c, "invalid post")
		return
	}

	var content = parseTmplFromResources("components/user/composeReply.html")

	err = content.ExecuteTemplate(c.Writer, "componentBody", post)
	if err != nil {
		c.AbortWithError(500, err)
	}

	c.Status(200)
}

func renderUserPostCommentsFeedComponent(c *gin.Context) {
	var id int

	if c.Param("id") != "" {
		var err error
		id, err = strconv.Atoi(c.Param("id"))
		if err != nil {
			renderErrorAlert(c, "invalid id")
			return
		}
	}

	var post PostModel
	err := database.Where("id = ?", id).First(&post).Error
	if err != nil {
		renderErrorAlert(c, "invalid post")
		return
	}

	var comments []PostCommentModel
	err = database.Where("parent_id = ?", id).Find(&comments).Error
	if err != nil {
		renderErrorAlert(c, "invalid post")
		return
	}

	var content = parseTmplFromResources("components/user/postRepliesFeed.html")

	err = content.ExecuteTemplate(c.Writer, "componentBody", comments)
	if err != nil {
		c.AbortWithError(500, err)
	}

	c.Status(200)
}

func renderRulesComponent(c *gin.Context) {
	var content = parseTmplFromResources("components/rules.html")

	err := content.ExecuteTemplate(c.Writer, "componentBody", nil)
	if err != nil {
		c.AbortWithError(500, err)
	}

	c.Status(200)
}

func renderAdministrationMainComponent(c *gin.Context) {
	var content = parseTmplFromResources("components/administration/main.html")

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
	var id int

	if c.Param("id") != "" {
		var err error
		id, err = strconv.Atoi(c.Param("id"))
		if err != nil {
			renderErrorAlert(c, "invalid id")
			return
		}
	}

	var category CategoryModel
	err := database.Where("id = ?", id).First(&category).Error
	if err != nil {
		renderErrorAlert(c, "invalid category")
		return
	}

	var content = parseTmplFromResources("components/administration/editCategory.html")

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
