import request from '@/utils/request'
// 添加元素
export function deleteSystemAlert(modelData) {
    return request({
        url: 'sys/alert/deleteSystemAlert',
        method: 'post',
        data: modelData
    })
}
// 添加元素
export function updateSystemAlert(modelData) {
    return request({
        url: 'sys/alert/updateSystemAlert',
        method: 'post',
        data: modelData
    })
}
// 添加元素
export function getSystemAlertList(modelData) {
    return request({
        url: 'sys/alert/getSystemAlertList',
        method: 'post',
        data: modelData
    })
}