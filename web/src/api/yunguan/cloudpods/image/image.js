import request from '@/utils/request'

//获取系统镜像
export function GetSystemImage() {
    return request({
        url: 'cloudpods/systemimage/list',
        method: 'get',
    })
}
//获取容器镜像
export function GetDockerImage() {
    return request({
        url: 'cloudpods/dockerimage/list',
        method: 'get',
    })
}
