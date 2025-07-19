import request from '@/utils/request'
//获取
export function getStorageList(params){
    return request({
        url:'cloud/storage/vmNas/list',
        method:'post',
        data: params
    })
}
//删除NAS存储
export function deleteNasStorage(params){
    return request({
        url:'cloud/storage/vmNas/delete',
        method:'post',
        data: params
    })
}
//删除NAS存储
export function deleteNasStorageByIds(params){
    return request({
        url:'cloud/storage/vmNas/deleteByIds',
        method:'post',
        data: params
    })
}
//添加NAS存储
export function addNasStorage(params){
    return request({
        url:'cloud/storage/vmNas/create',
        method:'post',
        data: params
    })
}
//禁用存储
export function disableNasStorage(params){
    return request({
        url:'cloud/storage/vmNas/disable',
        method:'post',
        data: params
    })
}
//启用存储
export function enableNasStorage(params){
    return request({
        url:'cloud/storage/vmNas/enable',
        method:'post',
        data: params
    })
}
//宿主机关联NAS存储
export function nasStorageConnHost(params){
    return request({
        url:'cloud/storage/vmNas/vmNasConnHost',
        method:'post',
        data: params
    })
}
//宿主机取消NAS存储关联
export function nasStorageDisConnHost(params){
    return request({
        url:'cloud/storage/vmNas/vmNasDisConnHost',
        method:'post',
        data: params
    })
}
//NAS存储关联的宿主机列表
export function nasStorageConnHostList(params){
    return request({
        url:'cloud/storage/vmNas/vmNasConnHostList',
        method:'post',
        data: params
    })
}
//NAS存储未关联的宿主机列表
export function nasStorageNoConnHostList(params){
    return request({
        url:'cloud/storage/vmNas/vmNasNoConnHostList',
        method:'post',
        data: params
    })
}


//容器NAS存储

//获取容器NAS存储列表
export function getContainerNasList(params){
    return request({
        url:'cloud/storage/containerNas/list',
        method:'post',
        data: params
    })
}
//创建容器NAS存储
export function createContainerNas(params){
    return request({
        url:'cloud/storage/containerNas/create',
        method:'post',
        data: params
    })
}
//删除容器NAS存储
export function deleteContainerNas(params){
    return request({
        url:'cloud/storage/containerNas/delete',
        method:'post',
        data: params
    })
}
//删除容器NAS存储
export function deleteContainerNasByIds(params){
    return request({
        url:'cloud/storage/containerNas/deleteByIds',
        method:'post',
        data: params
    })
}
//获取保密字典列表
export function getSecretList(params){
    return request({
        url:'cloud/storage/containerNas/secretList',
        method:'post',
        data: params
    })
}
//创建保密字典
export function secretCreate(params){
    return request({
        url:'cloud/storage/containerNas/secretCreate',
        method:'post',
        data: params
    })
}
//删除保密字典
export function secretDelete(params){
    return request({
        url:'cloud/storage/containerNas/secretDelete',
        method:'post',
        data: params
    })
}
//获取ceph集群的id
export function getCephClusterId(params){
    return request({
        url:'cloud/storage/containerNas/getCephClusterId',
        method:'post',
        data: params
    })
}
//获取ceph集群的存储池列表
export function getCephPoolsList(params){
    return request({
        url:'cloud/storage/containerNas/getCephPools',
        method:'post',
        data: params
    })
}