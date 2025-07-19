import request from "@/utils/request";

// 添加系统运行工具
export function addSysTicket(modelData) {
  return request({
    url: "sys/ticket/addSysTicket",
    method: "post",
    data: modelData,
  });
}

// 修改系统运行工具
export function updateSysTicket(data) {
  return request({
    url: "sys/ticket/updateSysTicket",
    method: "put",
    data: data,
  });
}
// 分配工单
export function assignSysTicket(data) {
  return request({
    url: "sys/ticket/assignSysTicket",
    method: "put",
    data: data,
  });
}

// 获取系统运行工具列表
export function getSysTicketList(params) {
  return request({
    url: "sys/ticket/getSysTicketList",
    method: "post",
    data: params,
  });
}

// 删除系统运行工具
export function deleteSysTicket(ID) {
  return request({
    url: "sys/ticket/deleteSysTicket",
    method: "delete",
    data: { ID: ID },
  });
}

//报价
export function postQuotation(params) {
  return request({
    url: "sys/ticket/postQuotation",
    method: "post",
    data: params,
  });
}

//发票
export function uploadInvoice(params) {
  return request({
    url: "sys/ticket/uploadInvoice",
    method: "post",
    data: params,
  });
}

//合同
export function uploadContract(params) {
  return request({
    url: "sys/ticket/uploadContract",
    method: "post",
    data: params,
  });
}
//获取所有资源列表
export function getAllRes(params) {
  return request({
    url: "sys/ticket/getAllRes",
    method: "post",
    data: params,
  });
}
//获取未被分配资源列表
export function getNoAssignRes(params) {
  return request({
    url: "sys/ticket/getNoAssignRes",
    method: "post",
    data: params,
  });
}