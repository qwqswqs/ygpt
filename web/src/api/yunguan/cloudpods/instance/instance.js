import request from '@/utils/request'

//监控
export function monitorInstance(query) {
    return request({
        url: 'udp/instance/monitor',
        method: 'post',
        data: query
    })
}
//获取任务列表
export function listTaskProgress(query) {
    return request({
        url: 'udp/instance/listTask',
        method: 'post',
        data: query
    })
}
//添加任务
export function addRenterTask(query) {
    return request({
        url: 'udp/instance/addTask',
        method: 'post',
        data: query
    })
}

export function manageRenterTask(privateIp,id,manageMethod) {
    return request({
        url: 'udp/instance/manageTask',
        method: 'post',
        data: {privateIp, id, manageMethod}
    })
}
//下载元素
export function DownloadElement(query) {
    return request({
        url: 'udp/instance/download',
        method: 'post',
        data: query
    })
}
