//* Deployment（无状态） 相关API
import request from "@/utils/request"

/**
 * @description 获取Deployment列表
 * @param {*} params 请求参数
 * @returns AjaxResult
 */
export function listDeployment(params) {
    return request({
        url: 'cloud/k8s/deployment/list',
        method: 'post',
        data: params
    })
}

/**
 * @description 添加Deployment
 * @param {*} data 表单数据
 * @returns AjaxResult
 */
export function addDeployment(data) {
    return request({
        url: 'cloud/k8s/deployment/create',
        method: 'post',
        data: data
    })
}

/**
 * @description 删除Deployment
 * @param {*} data
 * @returns
 */
export function deleteDeployment(data) {
    return request({
        url: 'cloud/k8s/deployment/delete',
        method: 'delete',
        data: data
    })
}

/**
 * @description 获取WebTerminal连接信息
 * @param {int} params
 */
export function webConnectDeployment(params) {
    return request({
        url: 'cloud/k8s/deployment/webConsole',
        method: 'get',
        params: params
    })
}

export function getDeploymentPods(params){
    return request({
        url: 'cloud/k8s/deployment/pods',
        method: 'get',
        params: params
    })
}
//获取集群节点列表
export function getK8sNodeList(params){
    return request({
        url: 'cloud/k8s/deployment/getK8sNodeList',
        method: 'post',
        data: params
    })
}
//修改集群相关节点调度状态
export function modK8sNodeStatus(params){
    return request({
        url: 'cloud/k8s/deployment/modK8sNodeStatus',
        method: 'post',
        data: params
    })
}
//监控集群容器
export function deploymentMonitor(params){
    return request({
        url: 'cloud/k8s/deployment/deploymentMonitor',
        method: 'post',
        data: params
    })
}