import request from "@/utils/request";

// 添加系统运行工具
export function addSysType(modelData) {
  return request({
    url: "sys/type/addSysType",
    method: "post",
    data: modelData,
  });
}

// 修改系统运行工具
export function updateSysType(data) {
  return request({
    url: "sys/type/updateSysType",
    method: "put",
    data: data,
  });
}

// 获取系统运行工具列表
export function getSysTypeList(params) {
  return request({
    url: "sys/type/getSysTypeList",
    method: "post",
    data: params,
  });
}

// 删除系统运行工具
export function deleteSysType(params) {
  return request({
    url: "sys/type/deleteSysType",
    method: "delete",
    params
  });
}

export const deleteSysTypeByIds = (params) => {
  return request({
    url: 'sys/type/deleteSysTypeByIds',
    method: 'delete',
    params
  })
}

export const getCategoryList = (params) => {
  return request({
    url: '/sysType/getCategoryList',
    method: 'get',
    params
  })
}
export const getElemSysTypeList = (params) => {
  return request({
    url: '/sysType/getElemSysTypeList',
    method: 'get',
    params
  })
}

export const updateCategory = (data) => {
  return request({
    url: '/sysType/updateCategory',
    method: 'put',
    data
  })
}

export const addCategory = (data) => {
  return request({
    url: '/sysType/addCategory',
    method: 'post',
    data
  })
}
export const deleteCategoryByIds = (params) => {
  return request({
    url: '/sysType/deleteCategoryByIds',
    method: 'delete',
    params
  })
}
export const deleteCategory = (params) => {
  return request({
    url: '/sysType/deleteCategory',
    method: 'delete',
    params
  })
}
