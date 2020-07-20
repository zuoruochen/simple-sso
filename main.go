package main

import (
	"github.com/gin-gonic/gin"
	"io"
	"log"
	"os"
	"simple-sso/application"
	"simple-sso/middleware"
	"simple-sso/router"
	"simple-sso/router/clientLogin"
	"simple-sso/router/hello"
	"simple-sso/router/login"
	"simple-sso/util"
)

var Debug = true
func main() {
	if !Debug {
		f, _ := os.Create("./log/gin.log")
		gin.DefaultWriter = io.MultiWriter(f)
		log.SetOutput(f)
	}

	applications := application.Default()

	// init cas application
	casRouter := router.New()
	casRouter.RegisterRouter(login.Routers, hello.Routers)
	casApp := application.NewApp("CAS-Server", util.CAS_HOST, casRouter)
	casApp.LoadHTMLGlob("template/*")

	// init client application

	clientRouter := router.New()
	clientRouter.RegisterRouter(clientLogin.Routers)
	clientApp := application.NewApp("client hello", util.HELLO_CLIENT_HOST, clientRouter, middleware.SSOCheck())

	clientRouter2 := router.New()
	clientRouter2.RegisterRouter(clientLogin.Routers)
	clientApp2 := application.NewApp("client2 hello", util.HELLO_CLIENT2_HOST, clientRouter, middleware.SSOCheck())

	// add app to applications
	applications.Add(casApp, clientApp, clientApp2)

	//run all apps in applications
	applications.Run()

}
