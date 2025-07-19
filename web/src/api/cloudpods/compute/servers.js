import request from '@/utils/request'

//恢复回收主机
export function recoverRecycleHost(params) {
    return request({
        url: 'cloud/compute/servers/recoverRecycleHost',
        method: 'post',
        data: params
    })
}
//清除回收主机
export function clearRecycleHost(params) {
    return request({
        url: 'cloud/compute/servers/clearRecycleHost',
        method: 'post',
        data: params
    })
}
//批量清除回收主机
export function clearRecycleHostByIds(params) {
    return request({
        url: 'cloud/compute/servers/clearRecycleHostByIds',
        method: 'post',
        data: params
    })
}
//获取回收站列表
export function getRecycleServerList(params) {
    return request({
        url: 'cloud/compute/servers/getRecycleServerList',
        method: 'post',
        data: params
    })
}

//获取可使用的GPU信息
export function getGpuList() {
    return request({
        url: 'cloud/compute/servers/getGpuList',
        method: 'post',
    })
}

//云主机列表
export function getHostList(query) {
    return request({
        url: 'cloud/compute/servers/getServerList',
        method: 'post',
        data: query
    })
}
//获取裸金属列表
export function getBaremetalList(query) {
    return request({
        url: 'cloud/compute/servers/getBaremetalList',
        method: 'post',
        data: query
    })
}
export  function  getWebSSHConn(params){
    return request({
        url:'cloud/compute/servers/getWebConsoleSSH',
        method:'post',
        data: params
    })
}
//删除
export  function  setHostDelete(params){
    return request({
        url:'cloud/compute/servers/hostDelete',
        method:'post',
        data: params
    })
}
//关机
export  function  setHostShutDown(params){
    return request({
        url:'cloud/compute/servers/hostShutDown',
        method:'post',
        data: params
    })
}
// 开机
export  function  setHostStart(params){
    return request({
        url:'cloud/compute/servers/hostStart',
        method:'post',
        data: params
    })
}
//重启
export  function  setHostRestart(params){
    return request({
        url:'cloud/compute/servers/hostRestart',
        method:'post',
        data: params
    })
}
//挂起
export  function  setHostSuspend(params){
    return request({
        url:'cloud/compute/servers/hostSuspend',
        method:'post',
        data: params
    })
}
//删除
export  function  deleteHost(params){
    return request({
        url:'cloud/compute/servers/hostDelete',
        method:'post',
        data: params
    })
}
//删除
export  function  deleteHostByIds(params){
    return request({
        url:'cloud/compute/servers/hostDeleteByIds',
        method:'post',
        data: params
    })
}
//恢复
export  function  setHostResume(params){
    return request({
        url:'cloud/compute/servers/hostResume',
        method:'post',
        data: params
    })
}
//创建
export  function  addHost(params){
    return request({
        url:'cloud/compute/servers/hostCreate',
        method:'post',
        data: params
    })
}

//查询套餐列表
export  function  getResSkuList(params){
    return request({
        url:'cloud/compute/servers/getResSkuList',
        method:'post',
        data: params
    })
}
//获取账号密码
export  function getLoginInfo(params){
    return request({
        url:'cloud/compute/servers/getLoginInfo',
        method:'post',
        data: params
    })
}
//获取账号密码
export  function resetPassword(params){
    return request({
        url:'cloud/compute/servers/resetPassword',
        method:'post',
        data: params
    })
}

// 裸金属创建
export function createBakemeats(data){
    return request({
        url:'cloud/compute/servers/baremetalCreate',
        method:'post',
        data:data
    })
}
// 修改名称
export function updateName(data){
    return request({
        url:'cloud/compute/servers/updateName',
        method:'post',
        data:data
    })
}
// 获取实例详情
export function hostDetails(data){
    return request({
        url:'cloud/compute/servers/hostDetails',
        method:'post',
        data:data
    })
}
//获取CPU监控数据
export  function  getMonitorCpuData(params){
    return request({
        url:'cloud/compute/servers/queryCpu',
        method:'post',
        data: params
    })
}
//获取内存监控数据
export  function  getMonitorMemData(params){
    return request({
        url:'cloud/compute/servers/queryMem',
        method:'post',
        data: params
    })
}
//获取磁盘监控数据
export  function  getMonitorDiskData(params){
    return request({
        url:'cloud/compute/servers/queryDisk',
        method:'post',
        data: params
    })
}
//获取磁盘读速率监控数据
export  function  getMonitorDiskReadData(params){
    return request({
        url:'cloud/compute/servers/queryDiskRead',
        method:'post',
        data: params
    })
}
//获取磁盘写速率监控数据
export  function  getMonitorDiskWriteData(params){
    return request({
        url:'cloud/compute/servers/queryDiskWrite',
        method:'post',
        data: params
    })
}
//获取网络入带宽监控数据
export  function  getMonitorBpsRecvData(params){
    return request({
        url:'cloud/compute/servers/queryBpsRecv',
        method:'post',
        data: params
    })
}
//获取网络出带宽监控数据
export  function  getMonitorBpsSentData(params){
    return request({
        url:'cloud/compute/servers/queryBpsSent',
        method:'post',
        data: params
    })
}
//获取网络收包数监控数据
export  function  getMonitorPpsRecvData(params){
    return request({
        url:'cloud/compute/servers/queryPpsRecv',
        method:'post',
        data: params
    })
}
//获取网络发包数监控数据
export  function  getMonitorPpsSentData(params){
    return request({
        url:'cloud/compute/servers/queryPpsSent',
        method:'post',
        data: params
    })
}
//获取GPU使用率
export  function  getMonitorGpuData(params){
    return request({
        url:'cloud/compute/servers/queryGpu',
        method:'post',
        data: params
    })
}
//获取GPU显存使用率
export  function  getMonitorGpuMemData(params){
    return request({
        url:'cloud/compute/servers/queryGpuMem',
        method:'post',
        data: params
    })
}
//获取Gpu温度
export  function  getMonitorGpuTemperatureData(params){
    return request({
        url:'cloud/compute/servers/queryGpuTemperature',
        method:'post',
        data: params
    })
}
//获取云主机可迁移的宿主机列表
export  function  getMigrateForecastHost(params){
    return request({
        url:'cloud/compute/servers/migrateForecast',
        method:'post',
        data: params
    })
}
//迁移
export  function  doMigrate(params){
    return request({
        url:'cloud/compute/servers/migrate',
        method:'post',
        data: params
    })
}
//探测是否有agent
export function forecastAgent(params){
    return request({
        url:'cloud/compute/servers/forecastAgent',
        method:'post',
        data: params
    })
}
//探测是否具备安装agent的条件
export function forecastAgentCanInstall(params){
    return request({
        url:'cloud/compute/servers/forecastAgentCanInstall',
        method:'post',
        data: params
    })
}
//输入ssh信息设置其能安装agent
export function setAgentCanInstall(params){
    return request({
        url:'cloud/compute/servers/setAgentCanInstall',
        method:'post',
        data: params
    })
}
//安装agent
export function installAgent(params){
    return request({
        url:'cloud/compute/servers/installAgent',
        method:'post',
        data: params
    })
}
//获取gpu信息列表
export function getGpuInfo(params){
    return request({
        url:'cloud/compute/servers/getGpuInfo',
        method:'post',
        data: params
    })
}
//获取gpu信息列表
export function manageGpu(params){
    return request({
        url:'cloud/compute/servers/manageGpu',
        method:'post',
        data: params
    })
}