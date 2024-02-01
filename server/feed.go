package main

import (
	"log"

	"github.com/gin-gonic/gin"
)

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
	log.Println("renderFeedComponent")
	offset := getParamterInt(c, "offset")

	posts, err := getRecentPosts(5, offset)
	if err != nil {
		renderErrorAlert(c, "error fetching posts")
		return
	}

	view := FeedView{
		GenericView:    NewGenericView(c),
		Posts:          posts,
		CurrentOffset:  offset,
		NextPageOffset: offset + 1,
		NextPageExists: len(posts) >= 5,
		LastPageExists: offset > 0,
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

func renderFeedComponentCategorized(c *gin.Context) {
	var categoryID = getParamterInt(c, "categoryId")

	posts, err := getRecentPostsByCategory(categoryID, 5)
	if err != nil {
		renderErrorAlert(c, "error fetching posts")
		return
	}

	var content = parseHTMLTemplatesFromResources("components/feed.html")

	err = content.ExecuteTemplate(c.Writer, "componentBody", FeedView{Posts: posts, GenericView: NewGenericView(c)})
	if err != nil {
		c.AbortWithError(500, err)
	}

	c.Status(200)
}
