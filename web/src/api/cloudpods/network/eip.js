import request from '@/utils/request'
//获取
export function getNetEipList(params){
    return request({
        url:'cloud/network/eip/getNetEipList',
        method:'post',
        data: params
    })
}
//新增
export function addNetEip(params){
    return request({
        url:'cloud/network/eip/addNetEip',
        method:'post',
        data: params
    })
}
//删除
export function deleteNetEip(params){
    return request({
        url:'cloud/network/eip/deleteNetEip',
        method:'post',
        data: params
    })
}
//批量删除
export function deleteNetEipByIds(params){
    return request({
        url:'cloud/network/eip/deleteNetEipByIds',
        method:'post',
        data: params
    })
}
//绑定
export function bindEipServer(params){
    return request({
        url:'cloud/network/eip/bindEipServer',
        method:'post',
        data: params
    })
}
//解绑
export function unbindEipServer(params){
    return request({
        url:'cloud/network/eip/unbindEipServer',
        method:'post',
        data: params
    })
}
//获取可绑定云主机
export function getBindServerList(params){
    return request({
        url:'cloud/network/eip/getBindServerList',
        method:'post',
        data: params
    })
}
//修改带宽
export function modEipBandWidth(params){
    return request({
        url:'cloud/network/eip/modEipBandWidth',
        method:'post',
        data: params
    })
}