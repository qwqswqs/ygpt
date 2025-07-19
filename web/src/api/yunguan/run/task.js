import request from '@/utils/request'

// 添加元素
export function addRenterTask(modelData) {
    return request({
        url: 'renter/task/addRenterTask',
        method: 'post',
        data: modelData
    })
}

// 修改元素
export function updateRenterTask(params) {
    return request({
        url: 'renter/task/updateRenterTask',
        method: 'put',
        data: params
    })
}

// 获取节点元素列表
export function getRenterTaskList(params) {
    return request({
        url: 'renter/task/queryRenterTask',
        method: "post",
        data:params
    })
}

// 删除元素
export function deleteRenterTask(params) {
    return request({
        url: 'renter/task/deleteRenterTask',
        method: "delete",
        data: {
            id: params
        }
    })
}
