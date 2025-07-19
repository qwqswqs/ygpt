
import request from '@/utils/request'

// 获取资源列表
export function getSoftwareInfoList(params) {
    return request({
        url: 'soft/info/getSoftwareInfoList',
        method: "post",
        data: params
    })
}

//新增资源
export function addSoftwareInfo(params) {
    return request({
        url: 'soft/info/addSoftwareInfo',
        method: "post",
        data: params
    })
}

//编辑资源
export function updateSoftwareInfo(params) {
    return request({
        url: 'soft/info/updateSoftwareInfo',
        method: "put",
        data: params
    })
}

//删除资源
export function deleteSoftwareInfo(params) {
    return request({
        url: 'soft/info/deleteSoftwareInfo',
        method: "delete",
        data: params
    })
}




