package v1

import (
	"kubespace/server/api/v1/kubernetes"
	"kubespace/server/api/v1/system"
)

type ApiGroup struct {
	SystemApiGroup   system.ApiGroup
	KubernetesApiGroup  kubernetes.ApiGroup
}

var ApiGroupApp = new(ApiGroup)
