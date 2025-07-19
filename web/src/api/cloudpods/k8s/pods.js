//* 容器组相关API
import request from "@/utils/request"


/**
 * @description 获取集群对应的命名空间列表
 * @param {*} params 
 */
export function getNameSpaceList(params){
    return request({
        url: 'cloud/k8s/pod/listNameSpace',
        method:'get',
        params:params
    })
}

/**
 * @description 获取集群下的容器组列表
 * @param {*} params
 */
export function getPodsList (params) {
    return request({
        url:'cloud/k8s/pod/listPods',
        method:'get',
        params:params
    })
}