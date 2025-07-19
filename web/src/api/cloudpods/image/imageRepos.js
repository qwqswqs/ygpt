import request from '@/utils/request'

// 获取镜像仓库列表
export function getImageReposList(params) {
    return request({
        url: 'cloud/image/imageRepos/getImageReposList',
        method: 'post',
        data: params
    })
}
// 新增镜像仓库
export function addImageRepos(params) {
    return request({
        url: 'cloud/image/imageRepos/addImageRepos',
        method: 'post',
        data: params
    })
}
// 删除镜像仓库
export function deleteImageRepos(params) {
    return request({
        url: 'cloud/image/imageRepos/deleteImageRepos',
        method: 'post',
        data: params
    })
}
// 批量删除镜像仓库
export function deleteImageReposByIds(params) {
    return request({
        url: 'cloud/image/imageRepos/deleteImageReposByIds',
        method: 'post',
        data: params
    })
}
// 获取镜像仓库的镜像列表
export function getImageReposImageList(params) {
    return request({
        url: 'cloud/image/imageRepos/getImageReposImageList',
        method: 'post',
        data: params
    })
}