package v1

import (
	"kubespace/server/api/v1/example"
	"kubespace/server/api/v1/kubernetes"
	"kubespace/server/api/v1/system"
)

type ApiGroup struct {
	SystemApiGroup   system.ApiGroup
	ExampleApiGroup  example.ApiGroup
	KubernetesApiGroup  kubernetes.ApiGroup
}

var ApiGroupApp = new(ApiGroup)
