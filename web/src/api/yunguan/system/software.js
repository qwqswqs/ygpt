import request from "@/utils/request";

// 添加系统运行工具
export function addSysSoftware(modelData) {
  return request({
    url: "sys/software/addSysSoftware",
    method: "post",
    data: modelData,
  });
}

// 修改系统运行工具
export function updateSysSoftware(data) {
  return request({
    url: "sys/software/updateSysSoftware",
    method: "put",
    data: data,
  });
}

// 获取系统运行工具列表
export function getSysSoftwareList(params) {
  return request({
    url: "sys/software/getSysSoftwareList",
    method: "post",
    data: params,
  });
}

// 删除系统运行工具
export function deleteSysSoftware(ID) {
  return request({
    url: "sys/software/deleteSysSoftware",
    method: "delete",
    data: { ID: ID },
  });
}
