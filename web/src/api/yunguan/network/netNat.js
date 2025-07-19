import request from '@/utils/request'

//添加
export function addNat(modelData){
    return request({
        url:'network/nat/addNat',

        method: 'post',
        data: modelData
    })
}

//修改
export function  updateNat(params){
    return request({
        url:'network/nat/updateNat',
        method: 'put',
        data: params
    })
}


//获取
export function  getNatList(params){
    return request({
        url:'network/nat/queryNat',
        method:'post',
        data:params
    })
}

//删除
export function  deleteNat(params){
    return request({
        url:'network/nat/deleteNat',
        method:'delete',
        data: {
            id: params.ID
        }
    })
}