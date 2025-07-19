import request from '@/utils/request'

// 添加元素
export function addDevice(modelData) {
    return request({
        url: 'res/device/addDevice',
        method: 'post',
        data: modelData
    })
}

// 修改元素
export function updateDevice(params) {
    return request({
        url: 'res/device/updateDevice',
        method: 'put',
        data: params
    })
}

// 获取节点元素列表
export function getDeviceList(params) {
    return request({
        url: 'res/device/queryDevice',
        method: "post",
        data:params
    })
}
// 获取节点元素列表
export function getAllServerList() {
    return request({
        url: 'res/device/queryAllServer',
        method: "post",
    })
}

// 删除元素
export function deleteDevice(params) {
    return request({
        url: 'res/device/deleteDevice',
        method: "delete",
        data: {
            id: params
        }
    })
}
