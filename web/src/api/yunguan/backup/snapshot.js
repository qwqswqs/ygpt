import request from '@/utils/request'
// 添加元素
export function addSnapshot(modelData) {
    return request({
        url: 'backup/snapshot/addSnapshot',
        method: 'post',
        data: modelData
    })
}

// 修改元素
export function updateSnapshot(params) {
    return request({
        url: 'backup/snapshot/updateSnapshot',
        method: 'put',
        data: params
    })
}

// 获取节点元素列表
export function getSnapshotList(params) {
    return request({
        url: 'backup/snapshot/querySnapshot',
        method: "post",
        data:params
    })
}
// 获取节点元素列表
export function getUserSnapshotList(params) {
    return request({
        url: 'backup/snapshot/queryUserSnapshot',
        method: "post",
        data:params
    })
}

// 删除元素
export function deleteSnapshot(params) {
    return request({
        url: 'backup/snapshot/deleteSnapshot',
        method: "delete",
        data: {
            id: params
        }
    })
}
