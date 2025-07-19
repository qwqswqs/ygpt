//     let formData = new FormData();
//     formData.append('id', modelData.id);
//     formData.append('name', modelData.name);
// // 假设 modelData.usage 是一个数组，例如：["item1", "item2"]
//     modelData.usage.forEach((item, index) => {
//         formData.append('usage[' + index + ']', item);
//     });
//     formData.append('description', modelData.description);
//     if (modelData.file) {
//         formData.append('file', modelData.file);
//     }
//     formData.append('type', modelData.type);
//     formData.append('price', modelData.price);
import request from '@/utils/request'

// 添加元素
export function addElement(modelData) {
    return request({
        url: 'product/elem/addElemInfo',
        headers: {
            "Content-type": "multipart/form-data"
        },
        method: 'post',
        data: modelData
    })
}

// 修改元素文件
export function updateElementFileInfo(modelData) {
    return request({
        url: 'product/elem/updateElemFileInfo',
        headers: {
            "Content-type": "multipart/form-data"
        },
        method: 'put',
        data: modelData
    })
}

// 修改元素
export function updateElement(params) {
    return request({
        url: 'product/elem/updateElemInfo',
        method: 'put',
        data: params
    })
}
// 同步单个元素
export function syncElement(params) {
    return request({
        url: 'product/elem/syncElemInfo',
        method: 'post',
        data:  params
    })
}

// 获取节点元素列表
export function getNodeElementList(params) {
    return request({
        url: 'product/elem/queryElemInfo',
        method: "post",
        data:params
    })
}

// 获取租户元素列表
export function getUserElementList(params) {
    return request({
        url: 'product/elem/queryUserElemInfo',
        method: "post",
        data: params
    })
}
// 获取共享元素列表
export function getShareElementList(params) {
    return request({
        url: 'product/elem/queryShareElemInfo',
        method: "post",
        data: params
    })
}

// 获取所有已上架模型列表
export function getAllModels() {
    return request({
        url: 'product/elem/getAllModels',
        method: 'get'
    })
}

// 获取所有已上架数据集列表
export function getAllDataSets() {
    return request({
        url: 'product/elem/getAllDataSets',
        method: 'get'
    })
}

// 获取共享元素列表
export function getAllShareElementList() {
    return request({
        url: 'product/elem/queryAllShareElemInfo',
        method: "post",
    })
}

// 删除元素
export function deleteElement(params) {
    return request({
        url: 'product/elem/deleteElemInfo',
        method: "delete",
        data: {
            id: params
        }
    })
}

//获取下载元素列表
export function GetDownloadElement(params) {
    return request({
        url: 'product/elemDownload/listDownloadedElement',
        method: "post",
        data: params
    })
}
//新增产品
export function createProduct (data) {
    return request({
        url: '/product/elem/createProduct',
        method: 'post',
        data
    })
}
//新增产品
export function updateProduct (data) {
    return request({
        url: '/product/elem/updateProduct',
        method: 'put',
        data
    })
}
//新增产品
export function getElementProductList (params) {
    return request({
        url: '/product/elem/getEleListInfo',
        method: 'get',
        params
    })
}
//新增产品
export function getUserSupplyList (params) {
    return request({
        url: '/product/elem/getUserSupplyList',
        method: 'get',
        params
    })
}
//单个产品删除借口
export const deleteProduct = (params) => {
    return request({
        url: '/product/elem/deleteProduct',
        method: 'delete',
        params
    })
}
export const deleteProductByIds = (params) => {
    return request({
        url: '/product/elem/deleteProductByIds',
        method: 'delete',
        params
    })
}