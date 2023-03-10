package wxp

import (
	"github.com/WebXense/env"
	"github.com/WebXense/ginger/ginger"
	"github.com/gin-gonic/gin"
)

var Server = server{
	Engine: ginger.NewEngine(),
}

type server struct {
	Engine *gin.Engine
}

func (s *server) Run() {
	s.Engine.Run(env.String("GIN_HOST"))
}

func (s *server) OK(ctx *gin.Context, data interface{}, page *ginger.Pagination) {
	ginger.OK(ctx, data, page)
}

func (s *server) ERR(ctx *gin.Context, err Error, data ...interface{}) {
	if len(data) > 0 {
		ginger.ERR(ctx, err.UUID(), err.Error(), data[0])
		return
	}
	ginger.ERR(ctx, err.UUID(), err.Error(), nil)
}
