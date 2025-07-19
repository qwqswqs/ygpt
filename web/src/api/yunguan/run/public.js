import request from '@/utils/request'

// 添加元素
export function addPublic(modelData) {
    return request({
        url: 'ops/public/addPublic',
        method: 'post',
        data: modelData
    })
}

// 修改元素
export function updatePublic(params) {
    return request({
        url: 'ops/public/updatePublic',
        method: 'put',
        data: params
    })
}

// 获取节点元素列表
export function getPublicList(params) {
    return request({
        url: 'ops/public/queryPublic',
        method: "post",
        data:params
    })
}

// 删除元素
export function deletePublic(params) {
    return request({
        url: 'ops/public/deletePublic',
        method: "delete",
        data: {
            id: params
        }
    })
}
