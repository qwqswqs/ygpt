import request from "@/utils/request";

// 添加系统运行工具
export function addSupplyInfo(modelData) {
    return request({
        url: "product/supply/addSupplyInfo",
        method: "post",
        data: modelData,
    });
}

// 修改系统运行工具
export function updateSupplyInfo(data) {
    return request({
        url: "product/supply/updateSupplyInfo",
        method: "put",
        data: data,
    });
}

// 获取系统运行工具列表
export function getSupplyInfoList(params) {
    return request({
        url: "product/supply/querySupplyInfo",
        method: "post",
        data: params,
    });
}

//同步单个数据
export function SyncSupplyDemandInfo(params) {
    return request({
        url: 'product/supply/syncSupplyInfo',
        method: "post",
        data: params
    })
}

// 删除系统运行工具
export function deleteSupplyInfo(ID) {
    return request({
        url: "product/supply/deleteSupplyInfo",
        method: "delete",
        data: { ID: ID },
    });
}

export const createProduct = (data) => {
    return request({
        url: '/product/supply/createProduct',
        method: 'post',
        data
    })
}

export const deleteProduct = (params) => {
    return request({
        url: '/product/supply/deleteProduct',
        method: 'delete',
        params
    })
}

export const deleteProductByIds = (params) => {
    return request({
        url: '/product/supply/deleteProductByIds',
        method: 'delete',
        params
    })
}

export const getSupplyOrDemandCityList = () => {
    return request({
        url: '/infoRegion/getInfoRegionList',
        method: 'get',
    })
}
//获取配置列表
export const getProductConfigList = (params) => {
    return request({
        url: '/sysType/getSysTypeList',
        method: 'get',
        params: params
    })
}

export const getProductList = (params) => {
    return request({
        url: '/product/supply/getProductList2',
        method: 'get',
        params: params
    })
}
//供需编辑接口
export const updateProduct = (data) => {
    return request({
        url: '/product/supply/updateProduct',
        method: 'put',
        data
    })
}
