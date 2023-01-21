package statics

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func StaticsMiddleware(c *gin.Context) {
	if c.Request.URL.Path == "/favicon.ico" {
		c.AbortWithStatus(http.StatusNoContent)
		return
	}

	c.Next()
}
