package clientLogin

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"simple-sso/data"
)

func helloWorldHandler(c *gin.Context) {
	c.JSON(http.StatusOK, data.Response{http.StatusOK, "hello world", nil})
}


