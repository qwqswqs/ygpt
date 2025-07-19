import request from '@/utils/request'

// 添加元素
export function addDisk(modelData) {
    return request({
        url: 'res/disk/addDisk',
        method: 'post',
        data: modelData
    })
}

// 修改元素
export function updateDisk(params) {
    return request({
        url: 'res/disk/updateDisk',
        method: 'put',
        data: params
    })
}

// 获取节点元素列表
export function getDiskList(params) {
    return request({
        url: 'res/disk/queryDisk',
        method: "post",
        data:params
    })
}

// 删除元素
export function deleteDisk(params) {
    return request({
        url: 'res/disk/deleteDisk',
        method: "delete",
        data: {
            id: params
        }
    })
}
