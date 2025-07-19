import request from '@/utils/request'
//获取
export function getVmEvsList(params){
    return request({
        url:'cloud/storage/vmEvs/list',
        method:'post',
        data: params
    })
}
//获取回收站硬盘列表
export function getRecycleVmEvsList(params){
    return request({
        url:'cloud/storage/vmEvs/getRecycleVmEvsList',
        method:'post',
        data: params
    })
}
//清除回收站硬盘
export function clearRecycleVmEvs(params){
    return request({
        url:'cloud/storage/vmEvs/clearRecycleVmEvs',
        method:'post',
        data: params
    })
}
//恢复回收站硬盘
export function recoverRecycleVmEvs(params){
    return request({
        url:'cloud/storage/vmEvs/recoverRecycleVmEvs',
        method:'post',
        data: params
    })
}
//获取
export function getContainerEvsList(params){
    return request({
        url:'cloud/storage/containerEvs/list',
        method:'post',
        data: params
    })
}
//获取
export function createVmEvs(params){
    return request({
        url:'cloud/storage/vmEvs/create',
        method:'post',
        data: params
    })
}
//获取
export function deleteVmEvs(params){
    return request({
        url:'cloud/storage/vmEvs/delete',
        method:'post',
        data: params
    })
}
//获取
export function deleteVmEvsByIds(params){
    return request({
        url:'cloud/storage/vmEvs/deleteByIds',
        method:'post',
        data: params
    })
}
//获取
export function createContainerEvs(params){
    return request({
        url:'cloud/storage/containerEvs/create',
        method:'post',
        data: params
    })
}
//获取
export function deleteContainerEvs(params){
    return request({
        url:'cloud/storage/containerEvs/delete',
        method:'post',
        data: params
    })
}
//获取
export function deleteContainerEvsByIds(params){
    return request({
        url:'cloud/storage/containerEvs/deleteByIds',
        method:'post',
        data: params
    })
}
//获取
export function resizeVmEvs(params){
    return request({
        url:'cloud/storage/vmEvs/resize',
        method:'post',
        data: params
    })
}
//获取
export function attachVmEvs(params){
    return request({
        url:'cloud/storage/vmEvs/attach',
        method:'post',
        data: params
    })
}
//获取
export function detachVmEvs(params){
    return request({
        url:'cloud/storage/vmEvs/detach',
        method:'post',
        data: params
    })
}
//获取
export function getAttachableVms(params){
    return request({
        url:'cloud/storage/vmEvs/getAttachableVms',
        method:'post',
        data: params
    })
}