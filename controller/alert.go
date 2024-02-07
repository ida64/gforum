package controller

import "github.com/gin-gonic/gin"

const (
	AlertTypeSuccess = "success"
	AlertTypeWarning = "warning"
	AlertTypeError   = "error"
)

type AlertView struct {
	Type       string
	Text       string
	AutoDelete bool
}

func NewAlertView(t string, text string, autoDelete bool) *AlertView {
	return &AlertView{
		Type:       t,
		Text:       text,
		AutoDelete: autoDelete,
	}
}

func renderAlert(c *gin.Context, alert *AlertView) {
	var content = parseHTMLTemplatesFromResources("components/alerts/alert.html")

	err := content.ExecuteTemplate(c.Writer, "componentBody", alert)
	if err != nil {
		c.AbortWithError(500, err)
	}

	c.Status(200)
}

func renderError(c *gin.Context, err error) {
	renderAlert(c, NewAlertView(AlertTypeError, err.Error(), false))
}

func renderSuccess(c *gin.Context, message string, autoDelete bool) {
	renderAlert(c, NewAlertView(AlertTypeSuccess, message, autoDelete))
}
