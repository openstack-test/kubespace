package router

import (
	"kubespace/server/router/example"
	"kubespace/server/router/kubernetes"
	"kubespace/server/router/system"
)

type RouterGroup struct {
	System   system.RouterGroup
	Example  example.RouterGroup
	Kubernetes kubernetes.RouterGroup
}

var RouterGroupApp = new(RouterGroup)
