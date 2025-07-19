import request from "@/utils/request"

//获取宿主机列表
export function getBaseHostList(params) {
    return request({
        url: 'cloud/baseRes/baseHost/getBaseHostList',
        method: 'post',
        data: params
    })
}
//获取物理机列表
export function getBaseBareHostList(params) {
    return request({
        url: 'cloud/baseRes/baseHost/getBaseBareHostList',
        method: 'post',
        data: params
    })
}
//获取IMPI信息
export function getImpiInfo(params) {
    return request({
        url: 'cloud/baseRes/baseHost/getImpiInfo',
        method: 'post',
        data: params
    })
}
//获取物理机登录信息
export function getLoginInfo(params) {
    return request({
        url: 'cloud/baseRes/baseHost/getLoginInfo',
        method: 'post',
        data: params
    })
}
//修改物理机状态
export function updateBareHostStatus(params) {
    return request({
        url: 'cloud/baseRes/baseHost/updateBareHostStatus',
        method: 'post',
        data: params
    })
}
//修改物理机备注
export function updateBareHostDescription(params) {
    return request({
        url: 'cloud/baseRes/baseHost/updateBareHostDescription',
        method: 'post',
        data: params
    })
}
//删除物理机
export function deleteBareHost(params) {
    return request({
        url: 'cloud/baseRes/baseHost/deleteBareHost',
        method: 'post',
        data: params
    })
}
//批量删除物理机
export function deleteBareHostByIds(params) {
    return request({
        url: 'cloud/baseRes/baseHost/deleteBareHostByIds',
        method: 'post',
        data: params
    })
}
//新增物理机
export function addBareHost(params) {
    return request({
        url: 'cloud/baseRes/baseHost/addBareHost',
        method: 'post',
        data: params
    })
}
//远程连接物理机
export function getWebConsoleSSH(params) {
    return request({
        url: 'cloud/baseRes/baseHost/getWebConsoleSSH',
        method: 'post',
        data: params
    })
}
//获取CPU监控数据
export  function  getMonitorHostCpuData(params){
    return request({
        url:'cloud/baseRes/baseHost/queryCpu',
        method:'post',
        data: params
    })
}
//获取内存监控数据
export  function  getMonitorHostMemData(params){
    return request({
        url:'cloud/baseRes/baseHost/queryMem',
        method:'post',
        data: params
    })
}
//获取磁盘监控数据
export  function  getMonitorHostDiskData(params){
    return request({
        url:'cloud/baseRes/baseHost/queryDisk',
        method:'post',
        data: params
    })
}
//获取磁盘读速率监控数据
export  function  getMonitorHostDiskReadData(params){
    return request({
        url:'cloud/baseRes/baseHost/queryDiskRead',
        method:'post',
        data: params
    })
}
//获取磁盘写速率监控数据
export  function  getMonitorHostDiskWriteData(params){
    return request({
        url:'cloud/baseRes/baseHost/queryDiskWrite',
        method:'post',
        data: params
    })
}
//获取网络入带宽监控数据
export  function  getMonitorHostBpsRecvData(params){
    return request({
        url:'cloud/baseRes/baseHost/queryBpsRecv',
        method:'post',
        data: params
    })
}
//获取网络出带宽监控数据
export  function  getMonitorHostBpsSentData(params){
    return request({
        url:'cloud/baseRes/baseHost/queryBpsSent',
        method:'post',
        data: params
    })
}