import request from '@/utils/request'

//添加
export function addVpc(modelData){
    return request({
        url:'network/vpc/addVpc',

        method: 'post',
        data: modelData
    })
}

//修改
export function  updateVpc(params){
    return request({
        url:'network/vpc/updateVpc',
        method: 'put',
        data: params
    })
}


//获取
export function  getVpcList(params){
    return request({
        url:'network/vpc/getVpcList',
        method:'post',
        data:params
    })
}
//获取
export function  getAllVpcList(){
    return request({
        url:'network/vpc/getAllVpcList',
        method:'post',
    })
}

//删除
export function  deleteVpc(params){
    return request({
        url:'network/vpc/deleteVpc',
        method:'delete',
        data: {
            id: params.ID
        }
    })
}
