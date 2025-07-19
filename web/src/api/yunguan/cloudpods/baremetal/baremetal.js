import request from '@/utils/request'

//获取物理机列表
export function GetBareMetalList(params) {
    return request({
        url: 'cloudpods/bareMetal/list',
        method: 'post',
        data:params,
    })
}