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

	var numUsers int64
	database.Model(&UserModel{}).Count(&numUsers)
	view.NumUsers = numUsers

	var numPosts int64
	database.Model(&PostModel{}).Count(&numPosts)
	view.NumPosts = numPosts

	var categories []CategoryModel
	database.Find(&categories)
	view.Categories = categories

	return &view
}

func renderColorModesComponent(c *gin.Context) {
	var content = parseTmplFromResources("components/colorModes.html")

	err := content.ExecuteTemplate(c.Writer, "componentBody", nil)
	if err != nil {
		c.AbortWithError(500, err)
	}

	c.Status(200)
}

type NavigationView struct {
	GenericView    *GenericView
	CategoryModels []CategoryModel
}

func renderNavigationComponent(c *gin.Context) {
	var categories []CategoryModel
	err := database.Find(&categories).Error
	if err != nil {
		renderErrorAlert(c, "Error fetching categories")
		return
	}

	var content = parseTmplFromResources("components/navigation.html")

	err = content.ExecuteTemplate(c.Writer, "componentBody", NavigationView{GenericView: NewGenericView(c), CategoryModels: categories})
	if err != nil {
		c.AbortWithError(500, err)
	}

	c.Status(200)
}

type FeedView struct {
	GenericView    *GenericView
	Posts          []PostModel
	CurrentOffset  int
	NextPageOffset int
	NextPageExists bool
	LastPageExists bool
	LastPageOffset int
}

func renderFeedComponent(c *gin.Context) {
	var offset int = 0

	if c.Param("offset") != "" {
		var err error
		offset, err = strconv.Atoi(c.Param("offset"))
		if err != nil {
			renderErrorAlert(c, "invalid offset")
			return
		}
	}

	posts, err := getRecentPosts(5, offset)
	if err != nil {
		renderErrorAlert(c, "error fetching posts")
		return
	}

	var view = FeedView{
		GenericView:    NewGenericView(c),
		Posts:          posts,
		CurrentOffset:  offset,
		NextPageOffset: offset + 1,
	}

	if len(posts) < 5 {
		view.NextPageExists = false
	} else {
		view.NextPageExists = true
	}

	if offset > 0 {
		view.LastPageExists = true
		view.LastPageOffset = offset - 1
	} else {
		view.LastPageExists = false
	}

	var content = parseTmplFromResources("components/feed.html")

	err = content.ExecuteTemplate(c.Writer, "componentBody", view)
	if err != nil {
		c.AbortWithError(500, err)
	}

	c.Status(200)
}

func renderFeedComponentCategorized(c *gin.Context) {
	var categoryID int

	if c.Param("category") != "" {
		var err error
		categoryID, err = strconv.Atoi(c.Param("category"))
		if err != nil {
			renderErrorAlert(c, "invalid category")
			return
		}
	}

	posts, err := getRecentPostsByCategory(categoryID, 10)
	if err != nil {
		renderErrorAlert(c, "Error fetching posts")
		return
	}

	var content = parseTmplFromResources("components/feed.html")

	err = content.ExecuteTemplate(c.Writer, "componentBody", FeedView{Posts: posts, GenericView: NewGenericView(c)})
	if err != nil {
		c.AbortWithError(500, err)
	}

	c.Status(200)
}

func renderSidebarComponent(c *gin.Context) {
	var content = parseTmplFromResources("components/sidebar.html")

	err := content.ExecuteTemplate(c.Writer, "componentBody", NewGenericView(c))
	if err != nil {
		c.AbortWithError(500, err)
	}

	c.Status(200)
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
