import request from '@/utils/request'

// 添加元素
export function addConfig(modelData) {
    return request({
        url: 'sys/config/addSysConfig',
        method: 'post',
        data: modelData
    })
}

// 修改元素
export function updateConfig(params) {
    return request({
        url: 'sys/config/updateSysConfig',
        method: 'put',
        data: params
    })
}

// 获取节点元素列表
export function getConfigList(params) {
    return request({
        url: 'sys/config/getSysConfigList',
        method: "post",
        data:params
    })
}

// 删除元素
export function deleteConfig(params) {
    return request({
        url: 'sys/config/deleteSysConfig',
        method: "delete",
        data: {
            id: params
        }
    })
}
