package kubernetes

import (
	"github.com/gin-gonic/gin"
	v1 "kubespace/server/api/v1"
	"kubespace/server/middleware"
)

type ClusterRouter struct {}

func (c ClusterRouter) InitClusterRouter(r *gin.RouterGroup) {
	K8sClusterRouter := r.Group("k8s").Use(middleware.OperationRecord())
	k8s := v1.ApiGroupApp.KubernetesApiGroup.ClusterApi
	{
		K8sClusterRouter.POST("cluster", k8s.CreateK8SCluster)
		K8sClusterRouter.GET("cluster/secret", k8s.ClusterSecret)
		K8sClusterRouter.POST("cluster/delete", k8s.DelK8SCluster)
		K8sClusterRouter.GET("cluster", k8s.ListK8SCluster)
		K8sClusterRouter.GET("cluster/detail", k8s.GetK8SClusterDetail)
		K8sClusterRouter.GET("events", k8s.Events)
	}
}

