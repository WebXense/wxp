package server

import (
	"github.com/WebXense/env"
	"github.com/WebXense/ginger/ginger"
	"github.com/gin-gonic/gin"

	"github.com/WebXense/wxp/errs"
)

var Engine = ginger.NewEngine()

func Run() {
	Engine.Run(env.String("GIN_HOST"))
}

func OK(ctx *gin.Context, data interface{}, page *ginger.Pagination) {
	ginger.OK(ctx, data, page)
}

func ERR(ctx *gin.Context, err errs.Error, data ...interface{}) {
	if len(data) > 0 {
		ginger.ERR(ctx, err.UUID(), err.Error(), data[0])
		return
	}
	ginger.ERR(ctx, err.UUID(), err.Error(), nil)
}
