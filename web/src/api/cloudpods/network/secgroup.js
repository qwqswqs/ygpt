import request from '@/utils/request'
//获取
export function getSecgroupList(params){
    return request({
        url:'cloud/network/secgroup/getSecgroupList',
        method:'post',
        data: params
    })
}
//创建VPC
export function addSecgroup(params){
    return request({
        url:'cloud/network/secgroup/addSecgroup',
        method:'post',
        data: params
    })
}
//删除VPC
export function deleteSecgroup(params){
    return request({
        url:'cloud/network/secgroup/deleteSecgroup',
        method:'post',
        data: params
    })
}
//批量删除VPC
export function deleteSecgroupByIds(params){
    return request({
        url:'cloud/network/secgroup/deleteSecgroupByIds',
        method:'post',
        data: params
    })
}
//获取
export function getSecgroupRuleList(params){
    return request({
        url:'cloud/network/secgroup/getSecgroupRuleList',
        method:'post',
        data: params
    })
}
//创建VPC
export function addSecgroupRule(params){
    return request({
        url:'cloud/network/secgroup/addSecgroupRule',
        method:'post',
        data: params
    })
}
//创建VPC
export function updateSecgroupRule(params){
    return request({
        url:'cloud/network/secgroup/updateSecgroupRule',
        method:'post',
        data: params
    })
}
//删除VPC
export function deleteSecgroupRule(params){
    return request({
        url:'cloud/network/secgroup/deleteSecgroupRule',
        method:'post',
        data: params
    })
}
//删除VPC
export function deleteSecgroupRuleByIds(params){
    return request({
        url:'cloud/network/secgroup/deleteSecgroupRuleByIds',
        method:'post',
        data: params
    })
}
//创建VPC
export function addSeverSecgroupRule(params){
    return request({
        url:'cloud/network/secgroup/addSeverSecgroupRule',
        method:'post',
        data: params
    })
}
//删除VPC
export function revokeSeverSecgroupRule(params){
    return request({
        url:'cloud/network/secgroup/revokeSeverSecgroupRule',
        method:'post',
        data: params
    })
}