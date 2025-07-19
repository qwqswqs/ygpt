import request from "@/utils/request";

// 添加系统运行工具
export function addCategoryInfo(modelData) {
    return request({
        url: "product/Category/addCategoryInfo",
        method: "post",
        data: modelData,
    });
}


// 获取系统运行工具列表
export function getCategoryInfoList() {
    return request({
        url: "product/Category/queryCategoryInfo",
        method: "post",
    });
}

// 删除系统运行工具
export function deleteCategoryInfo(modelData) {
    return request({
        url: "product/Category/deleteCategoryInfo",
        method: "delete",
        data: modelData,
    });
}
