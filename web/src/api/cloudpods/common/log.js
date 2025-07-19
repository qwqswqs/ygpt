import request from "@/utils/request"

//获取宿主机列表
export function getLogsList(params) {
    return request({
        url: 'cloud/common/log/getLogsList',
        method: 'post',
        data: params
    })
}