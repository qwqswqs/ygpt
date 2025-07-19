import request from "@/utils/request";

// 添加缓存信息
export function addSysCache(modelData) {
  return request({
    url: "sys/cache/addSysCache",
    method: "post",
    data: modelData,
  });
}

// 删除缓存信息
export function deleteSysCache(ID) {
  return request({
    url: "sys/cache/deleteSysCache",
    method: "delete",
    data: { ID: ID },
  });
}

// 获取缓存列表
export function getSysCacheList(params) {
  return request({
    url: "sys/cache/getSysCacheList",
    method: "get",
    params: params,
  });
}

// 修改缓存信息
export function updateSysCache(data) {
    return request({
      url: "sys/cache/updateSysCache",
      method: "put",
      data:data,
    });
  }

