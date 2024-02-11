package controller

import (
	"gforum/config"

	"github.com/dchest/captcha"
	"github.com/gin-gonic/gin"
)

type Route struct {
	Method   string
	Path     string
	Handlers []gin.HandlerFunc
}

var rootRoutes = []Route{
	{
		Method:   "GET",
		Path:     "/",
		Handlers: []gin.HandlerFunc{userMiddleware, renderIndexPage},
	},
	{
		Method:   "GET",
		Path:     "/feed/:id",
		Handlers: []gin.HandlerFunc{userMiddleware, restorePageMiddleware, renderFeedByCategory},
	},

	{
		Method:   "GET",
		Path:     "/user/register",
		Handlers: []gin.HandlerFunc{restorePageMiddleware, renderRegister},
	},
	{
		Method:   "GET",
		Path:     "/user/login",
		Handlers: []gin.HandlerFunc{restorePageMiddleware, renderLogin},
	},

	{
		Method:   "GET",
		Path:     "/user/compose",
		Handlers: []gin.HandlerFunc{userMiddleware, userEnforceMiddleware, restorePageMiddleware, renderCompose},
	},
	{
		Method:   "GET",
		Path:     "/user/profile",
		Handlers: []gin.HandlerFunc{userMiddleware, userEnforceMiddleware, restorePageMiddleware, renderProfile},
	},

	{
		Method:   "GET",
		Path:     "/post/:id",
		Handlers: []gin.HandlerFunc{userMiddleware, restorePageMiddleware, renderPost},
	},

	{
		Method:   "GET",
		Path:     "/management",
		Handlers: []gin.HandlerFunc{userMiddleware, userAdministratorEnforceMiddleware, restorePageMiddleware, renderManagement},
	},
	{
		Method:   "GET",
		Path:     "/management/categories",
		Handlers: []gin.HandlerFunc{userMiddleware, userAdministratorEnforceMiddleware, restorePageMiddleware, renderManagementCategories},
	},
}

var componentRoutes = []Route{
	{
		Method:   "GET",
		Path:     "/navigation",
		Handlers: []gin.HandlerFunc{userMiddleware, renderNavigation},
	},
	{
		Method:   "GET",
		Path:     "/feedSelector/:feed_id",
		Handlers: []gin.HandlerFunc{renderFeedSelector},
	},
	{
		Method:   "GET",
		Path:     "/feed/:offset",
		Handlers: []gin.HandlerFunc{renderFeed},
	},
	{
		Method:   "GET",
		Path:     "/user/:user_id/icon",
		Handlers: []gin.HandlerFunc{renderUserIcon},
	},
	{
		Method:   "GET",
		Path:     "/post/:post_id/feed",
		Handlers: []gin.HandlerFunc{renderPostFeed},
	},

	{
		Method:   "GET",
		Path:     "/management/navigation",
		Handlers: []gin.HandlerFunc{userMiddleware, userAdministratorEnforceMiddleware, renderManagementNavigation},
	},
}

var actionRoutes = []Route{
	{
		Method:   "POST",
		Path:     "/user/login",
		Handlers: []gin.HandlerFunc{loginAction},
	},
	{
		Method:   "POST",
		Path:     "/user/register",
		Handlers: []gin.HandlerFunc{registerAction},
	},

	{
		Method:   "PUT",
		Path:     "/user",
		Handlers: []gin.HandlerFunc{userMiddleware, userEnforceMiddleware, updateProfileAction},
	},
	{
		Method:   "PUT",
		Path:     "/user/avatar",
		Handlers: []gin.HandlerFunc{userMiddleware, userEnforceMiddleware, updateAvatarAction},
	},

	{
		Method:   "POST",
		Path:     "/user/compose",
		Handlers: []gin.HandlerFunc{userMiddleware, userEnforceMiddleware, composeAction},
	},

	{
		Method:   "POST",
		Path:     "/post/:id/comment",
		Handlers: []gin.HandlerFunc{userMiddleware, userEnforceMiddleware, addCommentAction},
	},
}

func addRoutes(group *gin.RouterGroup, routes []Route) {
	for _, route := range routes {
		group.Handle(route.Method, route.Path, route.Handlers...)
	}
}

var router *gin.Engine

func init() {
	router = gin.Default()

	addRoutes(router.Group("/"), rootRoutes)
	addRoutes(router.Group("/component"), componentRoutes)
	addRoutes(router.Group("/action"), actionRoutes)

	router.Static("/static/js", "resources/static/js")
	router.Static("/static/css", "resources/static/css")
	router.Static("/static/img", "resources/static/img")

	router.Static("/resources/avatars", "resources/avatars")

	captchaHandler := captcha.Server(captcha.StdWidth, captcha.StdHeight)

	router.GET("/captcha/*all", func(c *gin.Context) {
		captchaHandler.ServeHTTP(c.Writer, c.Request)
	})
}

func ListenServer() error {
	return router.Run(config.LoadedConfig.Server.Host)
}
