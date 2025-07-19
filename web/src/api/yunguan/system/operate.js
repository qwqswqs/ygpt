import request from "@/utils/request";

// 删除日志
export function deleteSysOperate(ID) {
  return request({
    url: "sys/operate/deleteSysOperate",
    method: "delete",
    data: { ID: ID },
  });
}

// 获取日志列表
export function getSysOperateList(params) {
  return request({
    url: "sys/operate/getSysOperateList",
    method: "post",
    data: params,
  });
}
