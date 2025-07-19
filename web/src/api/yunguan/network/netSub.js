import request from '@/utils/request'

//添加
export function addSub(modelData){
    return request({
        url:'network/sub/addSub',

        method: 'post',
        data: modelData
    })
}

//修改
export function  updateSub(params){
    return request({
        url:'network/sub/updateSub',
        method: 'put',
        data: params
    })
}


//获取
export function  getSubList(params){
    return request({
        url:'network/sub/querySub',
        method:'post',
        data:params
    })
}
//获取
export function  getAllSubList(){
    return request({
        url:'network/sub/queryAllSub',
        method:'post',
    })
}

//删除
export function  deleteSub(params){
    return request({
        url:'network/sub/deleteSub',
        method:'delete',
        data: {
            id: params.ID
        }
    })
}
