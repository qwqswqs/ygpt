import request from '@/utils/request'

// 镜像标签修改
export function updateImageTag(params) {
    return request({
        url: 'cloud/image/image/updateImageTag',
        method: 'post',
        data: params
    })
}
// 恢复回收镜像
export function recoverRecycleImage(params) {
    return request({
        url: 'cloud/image/image/recoverRecycleImage',
        method: 'post',
        data: params
    })
}
// 清除回收镜像
export function clearRecycleImage(params) {
    return request({
        url: 'cloud/image/image/clearRecycleImage',
        method: 'post',
        data: params
    })
}
// 获取回收站列表
export function getRecycleImageList(params) {
    return request({
        url: 'cloud/image/image/getRecycleImageList',
        method: 'post',
        data: params
    })
}
// 获取镜像列表
export function getImageList(params) {
    return request({
        url: 'cloud/image/image/getImageList',
        method: 'post',
        data: params
    })
}
// 获取单个镜像
export function getImageByName(params) {
    return request({
        url: 'cloud/image/image/getImageByName',
        method: 'post',
        data: params
    })
}
// 上传镜像
export function uploadImage(params, config = {}) {
    return request({
        url: 'cloud/image/image/uploadImage',
        method: 'post',
        data: params,
        // 添加支持进度监听的配置
        onUploadProgress: config.onUploadProgress,
        cancelToken: config.cancelToken
    })
}
// 上传镜像
export function deleteImage(params) {
    return request({
        url: 'cloud/image/image/deleteImage',
        method: 'delete',
        data: params
    })
}
// 上传镜像
export function deleteImageByIds(params) {
    return request({
        url: 'cloud/image/image/deleteImageByIds',
        method: 'delete',
        data: params
    })
}

export function getImagesByRepository(data) {
    return request({
        url: 'cloud/image/imageRepos/getImageReposImageList',
        method:'post',
        data: data 
    })
}