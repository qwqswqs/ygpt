import request from '@/utils/request'

//获取远程连接
export function GetContainSSH(params) {
    return request({
        url: 'cloudpods/container/connect',
        method: 'post',
        data:params,
    })
}