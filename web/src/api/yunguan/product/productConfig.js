import request from "@/utils/request";

// 添加系统运行工具
export function addConfigInfo(modelData) {
    return request({
        url: "product/config/addConfigInfo",
        method: "post",
        data: modelData,
    });
}

// 修改系统运行工具
export function updateConfigInfo(data) {
    return request({
        url: "product/config/updateConfigInfo",
        method: "put",
        data: data,
    });
}

// 获取系统运行工具列表
export function getConfigInfoList(params) {
    return request({
        url: "product/config/queryConfigInfo",
        method: "post",
        data: params,
    });
}
// 获取系统运行工具列表
export function getAllConfigInfoList() {
    return request({
        url: "product/config/queryAllConfigInfo",
        method: "post",
    });
}

// 删除系统运行工具
export function deleteConfigInfo(ID) {
    return request({
        url: "product/config/deleteConfigInfo",
        method: "delete",
        data: { ID: ID },
    });
}
