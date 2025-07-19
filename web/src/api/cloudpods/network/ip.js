import request from '@/utils/request'
//获取
export function getNetIpList(params){
    return request({
        url:'cloud/network/ip/getNetIpList',
        method:'post',
        data: params
    })
}
//新增
export function addNetIP(params){
    return request({
        url:'cloud/network/ip/addNetIP',
        method:'post',
        data: params
    })
}
//删除
export function deleteNetIP(params){
    return request({
        url:'cloud/network/ip/deleteNetIP',
        method:'post',
        data: params
    })
}
//批量删除
export function deleteNetIpByIds(params){
    return request({
        url:'cloud/network/ip/deleteNetIpByIds',
        method:'post',
        data: params
    })
}
//编辑
export function updateNetIP(params){
    return request({
        url:'cloud/network/ip/updateNetIP',
        method:'post',
        data: params
    })
}