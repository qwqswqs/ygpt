import request from '@/utils/request'

// 添加元素
export function addInfo(modelData) {
    return request({
        url: 'res/info/addInfo',
        method: 'post',
        data: modelData
    })
}

// 修改元素
export function updateInfo(params) {
    return request({
        url: 'res/info/updateInfo',
        method: 'put',
        data: params
    })
}

// 获取节点元素列表
export function getInfoList(params) {
    return request({
        url: 'res/info/queryInfo',
        method: "post",
        data:params
    })
}

// 获取租户资源列表
export function getRenterInfoList(params) {
    return request({
        url: 'res/info/queryRenterInfo',
        method: "post",
        data:params
    })
}

// 获取租户资源列表
export function getRenterResList(params) {
    return request({
        url: 'renter/res/queryRenterRes',
        method: "post",
        data:params
    })
}
// 获取租户财务明细列表
export function getRenterFinanceList(params) {
    return request({
        url: 'renter/res/queryRenterFinance',
        method: "get",
        params
    })
}
// 删除元素
export function deleteInfo(params) {
    return request({
        url: 'res/info/deleteInfo',
        method: "delete",
        data: {
            id: params
        }
    })
}
