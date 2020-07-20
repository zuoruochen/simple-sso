package hello

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func helloWorldHandler(c *gin.Context) {
	c.String(http.StatusOK, "hello World!")
}