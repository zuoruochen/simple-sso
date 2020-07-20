package router

import "github.com/gin-gonic/gin"

type RouterHandler func(engine *gin.Engine)

type Routers []RouterHandler

func (r *Routers)RegisterRouter(router...RouterHandler) {
	*r = append(*r, router...)
}

func New() Routers{
	return make(Routers, 0)
}
