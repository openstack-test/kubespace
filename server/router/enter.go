package router

import (
	"kubespace/server/router/kubernetes"
	"kubespace/server/router/system"
)

type RouterGroup struct {
	System   system.RouterGroup
	Kubernetes kubernetes.RouterGroup
}

var RouterGroupApp = new(RouterGroup)
