import request from '@/utils/request'

//添加
export function addIp(modelData){
    return request({
        url:'network/ip/addIp',
        method: 'post',
        data: modelData
    })
}

//修改
export function  updateIp(params){
    return request({
        url:'network/ip/updateIp',
        method: 'put',
        data: params
    })
}


//获取
export function  getIpList(params){
    return request({
        url:'network/ip/queryIp',
        method:'post',
        data: params
    })
}

//获取
export function  getAllIpList(){
    return request({
        url:'network/ip/queryAllIp',
        method:'post',
    })
}

//删除
export function  deleteIp(params){
    return request({
        url:'network/ip/deleteIp',
        method:'delete',
        data: {
                id: params.ID
        }
    })
}
