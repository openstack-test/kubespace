package kubernetes

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"kubespace/server/controller"
	"kubespace/server/global"
	"kubespace/server/model/common/response"
	"kubespace/server/model/kubernetes"
	"kubespace/server/service/kubernetes/cluster"
	"kubespace/server/service/kubernetes/event"
	"kubespace/server/service/kubernetes/parser"
	"strconv"
)

type ClusterApi struct {}

// CreateK8sCluster 创建k8s集群
// @Tags K8sCluster
// @Summary 创建k8s集群
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body kubernetes.K8SCluster true "创建k8s集群"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /k8s/cluster [post]
func (a *ClusterApi) CreateK8SCluster(c *gin.Context) {
	var K8sCluster kubernetes.K8SCluster
	err := controller.CheckParams(c, &K8sCluster)
	if err != nil {
		return
	}
	client, err := cluster.GetK8sClient(K8sCluster.KubeConfig)
	if err != nil {
		response.FailWithMessage("创建K8s集群错误", c)
		return
	}
	version, err := cluster.GetClusterVersion(client)
	if err != nil {
		response.FailWithMessage("连接集群异常,请检查网络是否畅通！", c)
		return
	}
	K8sCluster.ClusterVersion = version
	number, err := cluster.GetClusterNodesNumber(client)
	if err != nil {
		global.GVA_LOG.Error("获取集群节点数量异常", zap.Any("err", err))
	}
	K8sCluster.NodeNumber = number

	if err := cluster.CreateK8SCluster(K8sCluster); err != nil {
		global.GVA_LOG.Error("创建K8s集群错误", zap.Any("err", err))
		response.FailWithMessage("创建K8s集群错误", c)
		return
	} else {
		response.OkWithMessage("创建K8s集群成功", c)
		return
	}
}

// CreateK8sCluster 获取k8s集群
// @Tags K8sCluster
// @Summary 获取k8s集群
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query kubernetes.K8SCluster true "获取k8s集群"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /k8s/cluster [get]
func (a *ClusterApi) ListK8SCluster(c *gin.Context) {
	query := kubernetes.PaginationQ{}
	if c.ShouldBindQuery(&query) != nil {
		response.FailWithMessage("ShouldBindQuery失败", c)
		return
	}

	var K8sCluster []kubernetes.K8SCluster

	if err := cluster.ListK8SCluster(&query, &K8sCluster); err != nil {
		global.GVA_LOG.Error("获取集群失败", zap.Any("err", err))
		response.FailWithMessage("获取集群失败", c)
	} else {
		response.OkWithDetailed(response.PageResult{
			List:  K8sCluster,
			Total: query.Total,
			PageSize:  query.Size,
			Page:  query.Page,
		}, "获取集群成功", c)
	}
}

// CreateK8sCluster 删除k8s集群
// @Tags K8sCluster
// @Summary 删除k8s集群
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query kubernetes.K8SCluster true "删除k8s集群"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /k8s/cluster/delete [post]
func (a *ClusterApi) DelK8SCluster(c *gin.Context) {
	var id kubernetes.ClusterIds
	err := controller.CheckParams(c, &id)
	if err != nil {
		return
	}
	err2 := cluster.DelCluster(id)
	if err2 != nil {
		username, _ := c.Get("username")
		global.GVA_LOG.Error(fmt.Sprintf("用户：%s, 删除数据失败", username))
		response.FailWithMessage("删除K8s集群失败！", c)
		return
	}
	response.Ok(c)
	return
}

// CreateK8sCluster 获取k8s集群凭证
// @Tags K8sCluster
// @Summary 获取k8s集群凭证
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query kubernetes.K8SCluster true "获取k8s集群凭证"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /k8s/cluster/secret [get]
func (a *ClusterApi) ClusterSecret(c *gin.Context) {
	clusterId := c.DefaultQuery("clusterId", "1")
	clusterIdUint, err := strconv.ParseUint(clusterId, 10, 32)
	clusterConfig, err := cluster.GetK8sCluster(uint(clusterIdUint))
	if err != nil {
		global.GVA_LOG.Error("获取集群失败", zap.Any("err", err))
		response.FailWithMessage("获取集群凭证失败", c)
		return
	}
	data := map[string]interface{}{"secret": clusterConfig.KubeConfig, "name": clusterConfig.ClusterName}
	response.OkWithData(data, c)
	return
}

// CreateK8sCluster 获取k8s集群详情
// @Tags K8sCluster
// @Summary 获取k8s集群详情
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query kubernetes.K8SCluster true "获取k8s集群详情"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /k8s/cluster/detail [get]
func (a *ClusterApi) GetK8SClusterDetail(c *gin.Context) {
	client, err := cluster.ClusterID(c)
	if err != nil {
		response.FailWithMessage("获取K8s集群详情失败", c)
		return
	}
	data := cluster.GetClusterInfo(client)
	response.OkWithData(data, c)
}

// CreateK8sCluster 获取k8s事件
// @Tags K8sCluster
// @Summary 获取k8s事件
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query kubernetes.K8SCluster true "获取k8s事件"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /k8s/events [get]
func (a *ClusterApi) Events(c *gin.Context) {
	namespace := parser.ParseNamespaceParameter(c)
	client, err := cluster.ClusterID(c)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	field := fmt.Sprintf("type=%s", "Warning")
	data, err := event.GetClusterNodeEvent(client, namespace, field)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithData(data, c)
	return
}
