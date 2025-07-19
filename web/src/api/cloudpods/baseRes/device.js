import request from "@/utils/request";
import service from '@/utils/request';

//获取GPU列表
export function getBaseDeviceList(params) {
    return request({
        url: 'cloud/baseRes/device/getBaseDeviceList',
        method: 'post',
        data: params
    })
}
//获取裸金属GPU列表
export function getBaremetalDeviceList(params) {
    return request({
        url: 'cloud/baseRes/device/getBaremetalGpuList',
        method: 'post',
        data: params
    })
}

//获取首页cloudpods统计数据
export function getcloudpodsStatistic(params) {
    return service({
        url: 'cloudpods/resinfo',
        method: 'get',
        params
      })
}
