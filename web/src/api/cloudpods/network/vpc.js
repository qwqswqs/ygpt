import request from '@/utils/request'
//获取
export function getVpcList(params){
    return request({
        url:'cloud/network/vpc/getVpcList',
        method:'post',
        data: params
    })
}
//创建VPC
export function addVpc(params){
    return request({
        url:'cloud/network/vpc/addVpc',
        method:'post',
        data: params
    })
}
//删除VPC
export function deleteVpc(params){
    return request({
        url:'cloud/network/vpc/deleteVpc',
        method:'post',
        data: params
    })
}
//批量删除VPC
export function deleteVpcByIds(params){
    return request({
        url:'cloud/network/vpc/deleteVpcByIds',
        method:'post',
        data: params
    })
}