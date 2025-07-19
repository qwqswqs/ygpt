import request from "@/utils/request";

// 删除日志
export function deleteSysLog(ID) {
  return request({
    url: "sys/log/deleteSysLog",
    method: "delete",
    data: { ID: ID },
  });
}

// 获取日志列表
export function getSysLogList(params) {
  return request({
    url: "sys/log/getSysLogList",
    method: "post",
    data: params,
  });
}
