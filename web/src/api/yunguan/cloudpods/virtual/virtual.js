import request from '@/utils/request'

//远程连接
export function GetVirtualSSH(params) {
    return request({
        url: 'cloudpods/virtual_machine/connect',
        method: 'post',
        data:params,
    })
}
//关机
export  function  setHostShutDown(params){
    return request({
        url:'cloudpods/virtual_machine/hostShutDown',
        method:'post',
        data: params
    })
}
// 开机
export  function  setHostStart(params){
    return request({
        url:'cloudpods/virtual_machine/hostStart',
        method:'post',
        data: params
    })
}
//重启
export  function  setHostRestart(params){
    return request({
        url:'cloudpods/virtual_machine/hostRestart',
        method:'post',
        data: params
    })
}
//挂起
export  function  setHostSuspend(params){
    return request({
        url:'cloudpods/virtual_machine/hostSuspend',
        method:'post',
        data: params
    })
}
//恢复
export  function  setHostResume(params){
    return request({
        url:'cloudpods/virtual_machine/hostResume',
        method:'post',
        data: params
    })
}
//创建主机快照
export  function  CreateServerSnap(params){
    return request({
        url:'cloudpods/virtual_machine/createServerSnap',
        method:'post',
        data: params
    })
}