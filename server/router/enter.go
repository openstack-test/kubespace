package router

import (
	"kubespace/server/router/example"
	"kubespace/server/router/system"
)

type RouterGroup struct {
	System   system.RouterGroup
	Example  example.RouterGroup
}

var RouterGroupApp = new(RouterGroup)
