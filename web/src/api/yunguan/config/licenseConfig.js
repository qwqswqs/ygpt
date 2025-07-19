import request from '@/utils/request'
// 添加元素
export function updateLicense(modelData) {
    return request({
        url: 'config/licenseConfig/updateLicense',
        method: 'post',
        data: modelData
    })
}
// 添加元素
export function updateLicensePwd(modelData) {
    return request({
        url: 'config/licenseConfig/updateLicensePwd',
        method: 'post',
        data: modelData
    })
}
// 添加元素
export function queryLicense(modelData) {
    return request({
        url: 'config/licenseConfig/queryLicense',
        method: 'post',
        data: modelData
    })
}
// 添加元素
export function loginLicense(modelData) {
    return request({
        url: 'config/licenseConfig/loginLicense',
        method: 'post',
        data: modelData
    })
}