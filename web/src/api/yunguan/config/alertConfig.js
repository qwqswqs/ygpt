import request from '@/utils/request'
// 添加元素
export function addAlertConfig(modelData) {
    return request({
        url: 'config/alertConfig/addAlertConfig',
        method: 'post',
        data: modelData
    })
}
// 添加元素
export function updateAlertConfig(modelData) {
    return request({
        url: 'config/alertConfig/updateAlertConfig',
        method: 'post',
        data: modelData
    })
}
// 添加元素
export function deleteAlertConfig(modelData) {
    return request({
        url: 'config/alertConfig/deleteAlertConfig',
        method: 'post',
        data: modelData
    })
}
// 添加元素
export function queryAlertConfig(modelData) {
    return request({
        url: 'config/alertConfig/queryAlertConfig',
        method: 'post',
        data: modelData
    })
}