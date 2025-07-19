import request from '@/utils/request'
//获取
export function getSnapshotList(params){
    return request({
        url:'cloud/storage/snapshot/getSnapshotList',
        method:'post',
        data: params
    })
}
//回滚主机
export function snapshotReset(params){
    return request({
        url:'cloud/storage/snapshot/snapshotReset',
        method:'post',
        data: params
    })
}
//删除主机快照
export function deleteSnapshot(params){
    return request({
        url:'cloud/storage/snapshot/deleteSnapshot',
        method:'post',
        data: params
    })
}
//批量删除主机快照
export function deleteSnapshotByIds(params){
    return request({
        url:'cloud/storage/snapshot/deleteSnapshotByIds',
        method:'post',
        data: params
    })
}