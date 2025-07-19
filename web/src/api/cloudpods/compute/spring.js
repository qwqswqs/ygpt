import request from '@/utils/request'

//伸缩配置列表
export function getScalingConfigList(params) {
    return request({
        url: 'cloud/scaling/config/list',
        method: 'post',
        data: params
    })
}
//伸缩配置创建
export function addScalingConfig(params) {
    return request({
        url: 'cloud/scaling/config/create',
        method: 'post',
        data: params
    })
}
//伸缩配置删除
export function deleteScalingConfig(params) {
    return request({
        url: 'cloud/scaling/config/delete',
        method: 'post',
        data: params
    })
}
//伸缩配置删除
export function deleteScalingConfigByIds(params) {
    return request({
        url: 'cloud/scaling/config/deleteByIds',
        method: 'post',
        data: params
    })
}
//弹性伸缩列表
export function getScalingGroupList(params) {
    return request({
        url: 'cloud/scaling/group/list',
        method: 'post',
        data: params
    })
}
//弹性伸缩创建
export function addScalingGroup(params) {
    return request({
        url: 'cloud/scaling/group/create',
        method: 'post',
        data: params
    })
}
//弹性伸缩删除
export function deleteScalingGroup(params) {
    return request({
        url: 'cloud/scaling/group/delete',
        method: 'post',
        data: params
    })
}
//弹性伸缩删除
export function deleteScalingGroupByIds(params) {
    return request({
        url: 'cloud/scaling/group/deleteByIds',
        method: 'post',
        data: params
    })
}
//弹性伸缩启用
export function enableScalingGroup(params) {
    return request({
        url: 'cloud/scaling/group/enable',
        method: 'post',
        data: params
    })
}
//弹性伸缩禁用
export function disableScalingGroup(params) {
    return request({
        url: 'cloud/scaling/group/disable',
        method: 'post',
        data: params
    })
}
//弹性伸缩移除实例
export function detachScalingGroup(params) {
    return request({
        url: 'cloud/scaling/group/detach',
        method: 'post',
        data: params
    })
}
//伸缩策略列表
export function getScalingPolicyList(params) {
    return request({
        url: 'cloud/scaling/policy/list',
        method: 'post',
        data: params
    })
}
//伸缩策略创建
export function addScalingPolicy(params) {
    return request({
        url: 'cloud/scaling/policy/create',
        method: 'post',
        data: params
    })
}
//伸缩策略删除
export function deleteScalingPolicy(params) {
    return request({
        url: 'cloud/scaling/policy/delete',
        method: 'post',
        data: params
    })
}
//伸缩策略启用
export function enableScalingPolicy(params) {
    return request({
        url: 'cloud/scaling/policy/enable',
        method: 'post',
        data: params
    })
}
//伸缩策略禁用
export function disableScalingPolicy(params) {
    return request({
        url: 'cloud/scaling/policy/disable',
        method: 'post',
        data: params
    })
}