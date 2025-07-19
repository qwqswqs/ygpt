import request from "@/utils/request";

// 添加系统运行工具
export function addSysTool(modelData) {
  return request({
    url: "sys/tool/addSysTool",
    method: "post",
    data: modelData,
  });
}

// 修改系统运行工具
export function updateSysTool(data) {
  return request({
    url: "sys/tool/updateSysTool",
    method: "put",
    data: data,
  });
}

// 获取系统运行工具列表
export function getSysToolList(params) {
  return request({
    url: "sys/tool/getSysToolList",
    method: "get",
    params: params,
  });
}

// 删除系统运行工具
export function deleteSysTool(ID) {
  return request({
    url: "sys/tool/deleteSysTool",
    method: "delete",
    data: { ID: ID },
  });
}
