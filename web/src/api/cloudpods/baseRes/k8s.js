import request from "@/utils/request"

//获取集群列表
export function getBaseK8SList(params) {
    return request({
        url: 'cloud/baseRes/k8s/getBaseK8SList',
        method: 'post',
        data: params
    })
}
//导入集群
export function importBaseCluster(params) {
    return request({
        url: 'cloud/baseRes/k8s/importBaseCluster',
        method: 'post',
        data: params
    })
}
//删除集群
export function deleteBaseCluster(params) {
    return request({
        url: 'cloud/baseRes/k8s/deleteBaseCluster',
        method: 'post',
        data: params
    })
}
//删除集群
export function deleteBaseClusterByIds(params) {
    return request({
        url: 'cloud/baseRes/k8s/deleteBaseClusterByIds',
        method: 'post',
        data: params
    })
}
//修改集群
export function updateBaseCluster(params) {
    return request({
        url: 'cloud/baseRes/k8s/updateBaseCluster',
        method: 'post',
        data: params
    })
}
