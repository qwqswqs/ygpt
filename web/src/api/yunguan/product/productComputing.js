import request from "@/utils/request";

// 添加系统运行工具
export function addComputingInfo(modelData) {
    return request({
        url: "product/computing/addComputingInfo",
        method: "post",
        data: modelData,
    });
}
// 同步单个数据
export function syncedComputingInfo(modelData) {
    return request({
        url: "product/computing/syncedComputingInfo",
        method: "post",
        data: modelData,
    });
}
// 同步所有数据
export function syncedAllComputingInfo() {
    return request({
        url: "product/computing/syncedAllComputingInfo",
        method: "post",
    });
}

// 修改系统运行工具
export function updateComputingInfo(data) {
    return request({
        url: "product/computing/updateComputingInfo",
        method: "put",
        data: data,
    });
}

// 获取系统运行工具列表
export function getComputingInfoList(params) {
    return request({
        url: "product/computing/queryComputingInfo",
        method: "post",
        data: params,
    });
}

// 删除系统运行工具
export function deleteComputingInfo(ID) {
    return request({
        url: "product/computing/deleteComputingInfo",
        method: "delete",
        data: { ID: ID },
    });
}

//同步一条算力产品信息
export function syncComputingInfo(ID) {
    return request({
        url: "product/computing/syncComputingInfo",
        method: "post",
        data: { ID: ID },
    });
}