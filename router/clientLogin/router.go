package clientLogin

import "github.com/gin-gonic/gin"


func Routers(e *gin.Engine) {
	e.GET("/", helloWorldHandler)
}