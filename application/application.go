package application

import (
	"github.com/gin-gonic/gin"
	"golang.org/x/sync/errgroup"
	"log"
	"simple-sso/router"
)

type App struct {
	*gin.Engine
	addr string
	name string
}

func NewApp(name, addr string, routers router.Routers, middleware ...gin.HandlerFunc) *App {
	engine := gin.Default()
	// add middleware firstly
	if len(middleware) > 0 {
		engine.Use(middleware...)
	}
	// add router
	for _, opt := range routers {
		opt(engine)
	}
	return &App{
		engine,
		addr,
		name,
	}
}


type Applications map[string]*App

func Default() Applications {
	return make(Applications)
}

func (a Applications) Add(apps ...*App) {
	for _, app := range apps {
		a[app.name] = app
	}
}

func (a Applications) Run() {
	g := new(errgroup.Group)
	for _, value := range a {
		copyValue := value
		g.Go(func() error {
			log.Printf("application %s run at : %s\n", copyValue.name, copyValue.addr)
			return copyValue.Run(copyValue.addr)
		})
	}
	if err :=g.Wait(); err!= nil {
		log.Fatal(err)
	}
}