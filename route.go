package main

import (
	"github.com/dchest/captcha"
	"github.com/gin-gonic/gin"
)

type Route struct {
	Method   string
	Path     string
	Handlers []gin.HandlerFunc
}

var Routes = []Route{
	{
		Method:   "GET",
		Path:     "/",
		Handlers: []gin.HandlerFunc{renderRootIndexPage},
	},
	{
		Method:   "GET",
		Path:     "/rules",
		Handlers: []gin.HandlerFunc{restorePageMiddleware, renderRulesComponent},
	},
	{
		Method:   "GET",
		Path:     "/post/:id",
		Handlers: []gin.HandlerFunc{sessionToUserMiddleware, restorePageMiddleware, renderPostComponent},
	},
	{
		Method:   "GET",
		Path:     "/user/register",
		Handlers: []gin.HandlerFunc{restorePageMiddleware, renderUserRegisterComponent},
	}, {
		Method:   "GET",
		Path:     "/user/login",
		Handlers: []gin.HandlerFunc{restorePageMiddleware, renderUserLoginComponent},
	},
	{
		Method:   "GET",
		Path:     "/user/profile",
		Handlers: []gin.HandlerFunc{sessionToUserMiddleware, userRequiredMiddleware, restorePageMiddleware, renderUserProfileComponent},
	},
	{
		Method:   "GET",
		Path:     "/user/compose",
		Handlers: []gin.HandlerFunc{sessionToUserMiddleware, userRequiredMiddleware, restorePageMiddleware, renderUserComposeComponent},
	},
	{
		Method:   "GET",
		Path:     "/feed/:category",
		Handlers: []gin.HandlerFunc{restorePageMiddleware, renderFeedComponentCategorized},
	},
	{
		Method:   "GET",
		Path:     "/siteAdministration",
		Handlers: []gin.HandlerFunc{sessionToUserMiddleware, userRequiredMiddleware, adminRequiredMiddleware, restorePageMiddleware, renderAdministrationMainComponent},
	},
}

var Components = []Route{
	{
		Method:   "GET",
		Path:     "/sidebar",
		Handlers: []gin.HandlerFunc{sessionToUserMiddleware, renderSidebarComponent},
	},
	{
		Method:   "GET",
		Path:     "/colorModes",
		Handlers: []gin.HandlerFunc{renderColorModesComponent},
	},
	{
		Method:   "GET",
		Path:     "/navigation",
		Handlers: []gin.HandlerFunc{sessionToUserMiddleware, renderNavigationComponent},
	},
	{
		Method:   "GET",
		Path:     "/feed/:offset",
		Handlers: []gin.HandlerFunc{restorePageMiddleware, renderFeedComponent},
	},
	{
		Method:   "GET",
		Path:     "/user/logout",
		Handlers: []gin.HandlerFunc{sessionToUserMiddleware, userRequiredMiddleware, handleUserLogoutAction},
	},
	{
		Method:   "GET",
		Path:     "/user/composeReply/:id",
		Handlers: []gin.HandlerFunc{sessionToUserMiddleware, userRequiredMiddleware, renderUserComposeReplyComponent},
	},
	{
		Method:   "GET",
		Path:     "/administration/editCategory/:id",
		Handlers: []gin.HandlerFunc{sessionToUserMiddleware, userRequiredMiddleware, adminRequiredMiddleware, renderAdministratorEditCategoryComponent},
	},
}

var Actions = []Route{
	{
		Method:   "POST",
		Path:     "/user/register",
		Handlers: []gin.HandlerFunc{logActionMiddleware, handleUserRegisterAction},
	},
	{
		Method:   "POST",
		Path:     "/user/login",
		Handlers: []gin.HandlerFunc{logActionMiddleware, handleUserLoginAction},
	},
	{
		Method:   "POST",
		Path:     "/user/compose",
		Handlers: []gin.HandlerFunc{sessionToUserMiddleware, userRequiredMiddleware, logActionMiddleware, handleUserComposeAction},
	},
	{
		Method:   "POST",
		Path:     "/user/uploadAvatar",
		Handlers: []gin.HandlerFunc{sessionToUserMiddleware, userRequiredMiddleware, logActionMiddleware, handleUserUploadAvatarAction},
	},
	{
		Method:   "POST",
		Path:     "/user/updateProfile",
		Handlers: []gin.HandlerFunc{sessionToUserMiddleware, userRequiredMiddleware, logActionMiddleware, handleUserUpdateProfileAction},
	},
	{
		Method:   "POST",
		Path:     "/administration/category",
		Handlers: []gin.HandlerFunc{sessionToUserMiddleware, userRequiredMiddleware, adminRequiredMiddleware, logActionMiddleware, handleAdministrationAddCategoryAction},
	},
	{
		Method:   "PUT",
		Path:     "/administration/category/:id",
		Handlers: []gin.HandlerFunc{sessionToUserMiddleware, userRequiredMiddleware, adminRequiredMiddleware, logActionMiddleware, handleAdministrationEditCategoryAction},
	},
	{
		Method:   "DELETE",
		Path:     "/administration/category/:id",
		Handlers: []gin.HandlerFunc{sessionToUserMiddleware, userRequiredMiddleware, adminRequiredMiddleware, logActionMiddleware, handleAdministrationDeleteCategoryAction},
	},
}

func addRouters(group *gin.RouterGroup) {
	for _, route := range Routes {
		group.Handle(route.Method, route.Path, route.Handlers...)
	}

	components := group.Group("/components")
	for _, route := range Components {
		components.Handle(route.Method, route.Path, route.Handlers...)
	}

	actions := group.Group("/actions")
	for _, route := range Actions {
		actions.Handle(route.Method, route.Path, route.Handlers...)
	}
}

var router *gin.Engine

func init() {
	router = gin.Default()

	addRouters(router.Group("/"))

	router.Static("/static/js", "resources/static/js")
	router.Static("/static/css", "resources/static/css")
	router.Static("/static/img", "resources/static/img")

	router.Static("/resources/avatars", "resources/avatars")

	captchaHandler := captcha.Server(captcha.StdWidth, captcha.StdHeight)

	router.GET("/captcha/*all", func(c *gin.Context) {
		captchaHandler.ServeHTTP(c.Writer, c.Request)
	})
}

func listenServer() error {
	return router.Run(loadedConfig.Server.Host)
}
