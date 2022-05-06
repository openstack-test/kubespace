import service from '@/utils/request'

export const createK8sCluster = (data) => {
     return service({
         url: "/k8s/cluster",
         method: 'post',
         data
     })
 }


 export const getK8sClusterList = (data) => {
     return service({
         url: "/k8s/cluster",
         method: 'get',
         data
     })
 }


 export const findK8sCluster = (data) => {
     return service({
         url: "/k8s/cluster/detail",
         method: 'get',
         data
     })
 }

 export const deleteK8sCluster = (data) => {
    return service({
        url: "/k8s/cluster/detail",
        method: 'get',
        data
    })
}

export const deleteK8sClusterByIds = (data) => {
    return service({
        url: "/k8s/cluster/detail",
        method: 'get',
        data
    })
}

export const updateK8sCluster = (data) => {
    return service({
        url: "/k8s/cluster/detail",
        method: 'get',
        data
    })
}
