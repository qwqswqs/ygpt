import request from '@/utils/request'

export function getHostList(query) {
    return request({
        url: 'host/server/getServerList',
        method: 'get',
        params: query
    })
}
export  function  getWebSSHConn(params){
    return request({
        url:'host/server/getWebConsoleSSH',
        method:'post',
        data: params
    })
}
//关机
export  function  setHostShutDown(params){
    return request({
        url:'host/server/hostShutDown',
        method:'post',
        data: params
    })
}
// 开机
export  function  setHostStart(params){
    return request({
        url:'host/server/hostStart',
        method:'post',
        data: params
    })
}
//重启
export  function  setHostRestart(params){
    return request({
        url:'host/server/hostRestart',
        method:'post',
        data: params
    })
}
//挂起
export  function  setHostSuspend(params){
    return request({
        url:'host/server/hostSuspend',
        method:'post',
        data: params
    })
}
//删除
export  function  deleteHost(params){
    return request({
        url:'host/server/hostDelete',
        method:'post',
        data: params
    })
}
//恢复
export  function  setHostResume(params){
    return request({
        url:'host/server/hostResume',
        method:'post',
        data: params
    })
}
//创建
export  function  addHost(params){
    return request({
        url:'host/server/hostCreate',
        method:'post',
        data: params
    })
}
