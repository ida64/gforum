package controller

import "github.com/gin-gonic/gin"

type HttpErrorView struct {
	ErrorCode int
	Message   string
}

func (view *HttpErrorView) Render(c *gin.Context) {
	var content = parseHTMLTemplatesFromResources("httpError.html")

	err := content.ExecuteTemplate(c.Writer, "componentBody", view)
	if err != nil {
		c.AbortWithError(500, err)
	}

	c.Status(200)
}

func NewHttpErrorView(globalView *GlobalView, errorCode int, message string) *HttpErrorView {
	return &HttpErrorView{
		ErrorCode: errorCode,
		Message:   message,
	}
}

func sendErrorWithPage(c *gin.Context, errorCode int, message string) {
	var view = NewHttpErrorView(NewGlobalView(c), errorCode, message)
	view.Render(c)
	c.Abort()
}
