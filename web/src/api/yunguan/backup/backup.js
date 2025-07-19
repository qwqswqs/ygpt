import request from '@/utils/request'
// 添加元素
export function addBack(modelData) {
    return request({
        url: 'backup/back/addBack',
        headers: {
            "Content-type": "multipart/form-data"
        },
        method: 'post',
        data: modelData
    })
}

// 修改元素
export function updateBack(params) {
    return request({
        url: 'backup/back/updateBack',
        method: 'put',
        data: params
    })
}

// 获取节点元素列表
export function getBackList(params) {
    return request({
        url: 'backup/back/queryBack',
        method: "post",
        data:params
    })
}
// 获取租户网盘列表
export function getUserBackList(params) {
    return request({
        url: 'backup/back/queryUserBack',
        method: "post",
        data:params
    })
}

// 删除元素
export function deleteBack(params) {
    return request({
        url: 'backup/back/deleteBack',
        method: "delete",
        data: {
            id: params.ID
        }
    })
}
