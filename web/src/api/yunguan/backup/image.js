import request from '@/utils/request'
// 添加元素
export function addBackImage(modelData) {
    return request({
        url: 'backup/image/addBackImage',
        method: 'post',
        data: modelData
    })
}
// 添加元素
export function addBackImageFile(modelData) {
    return request({
        url: 'backup/image/addBackImageFile',
        headers: {
            "Content-type": "multipart/form-data"
        },
        method: 'post',
        data: modelData
    })
}

// 修改元素
export function updateBackImage(params) {
    return request({
        url: 'backup/image/updateBackImage',
        method: 'put',
        data: params
    })
}

// 获取节点元素列表
export function getBackImageList(params) {
    return request({
        url: 'backup/image/queryBackImage',
        method: "post",
        data:params
    })
}
// 获取节点元素列表
export function getBackImageAllList(params) {
    return request({
        url: 'backup/image/queryBackImageAll',
        method: "post",
    })
}
// 删除元素
export function deleteBackImage(params) {
    return request({
        url: 'backup/image/deleteBackImage',
        method: "delete",
        data: {
            id: params.ID
        }
    })
}
