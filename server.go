package wxp

import (
	"github.com/WebXense/env"
	"github.com/WebXense/ginger/ginger"
	"github.com/gin-gonic/gin"
)

var Engine = ginger.NewEngine()

func RunServer() {
	generateTypeScript()
	Engine.Run(env.String("GIN_HOST"))
}

func OK(ctx *gin.Context, data interface{}, page *ginger.Pagination) {
	ginger.OK(ctx, data, page)
}

func ERR(ctx *gin.Context, err Error, data ...interface{}) {
	if len(data) > 0 {
		ginger.ERR(ctx, err.UUID(), err.Error(), data[0])
		return
	}
	ginger.ERR(ctx, err.UUID(), err.Error(), nil)
}
