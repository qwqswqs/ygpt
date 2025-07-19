import service from '@/utils/request'
// 租户注册
export const addRenter = (data) => {
    return service({
        url: '/renter/renter/addRenter',
        method: 'post',
        data: data
    })
}
export const getRenterList = (data) => {
    return service({
        url: '/renter/renter/queryRenter',
        method: 'post',
        data: data
    })
}
export const getResRenterList = (data) => {
    return service({
        url: '/renter/renter/queryResRenter',
        method: 'post',
        data: data
    })
}
export const getUserList = (data) => {
    return service({
        url: '/renter/renter/getRenterList',
        method: 'post',
        data: data
    })
}
export const getAllRenterList = () => {
    return service({
        url: '/renter/renter/getAllRenterList',
        method: 'post',
    })
}
export const getAllUserList = () => {
    return service({
        url: '/renter/renter/getAllUserList',
        method: 'post',
    })
}

export const deleteRenter = (data) => {
    return service({
        url: '/renter/renter/deleteRenter',
        method: 'delete',
        data: data
    })
}
// 分配资源
export const bindResInfo = (modelData) =>{
    return service({
        url: "renter/res/bindResInfo",
        method: "put",
        data: modelData,
    });
}
// 绑定租户和资源
export const addRenterRes = (modelData) =>{
    return service({
        url: "renter/res/addRenterRes",
        method: "post",
        data: modelData,
    });
}
// 绑定租户和资源
export const deleteRenterRes = (modelData) =>{
    return service({
        url: "renter/res/deleteRenterRes",
        method: "delete",
        data: modelData,
    });
}
// 修改绑定租户和资源
export const updateRenterRes = (modelData) =>{
    return service({
        url: "renter/res/updateRenterRes",
        method: "put",
        data: modelData,
    });
}
// 获取所有租户资源列表
export function queryAllRenterRes  () {
    return service({
        url: "renter/res/queryAllRenterRes",
        method: "post",
    });
}
// 获取某个资源信息
export const queryRenterResInfo = (modelData) =>{
    return service({
        url: "renter/res/queryRenterResInfo",
        method: "post",
        data: modelData,
    });
}
// 获取与订单相关的租户资源列表
export const getRenterResByOrderId = (modelData) =>{
    return service({
        url: "renter/res/getRenterResByOrderId",
        method: "post",
        data: modelData,
    });
}
// 获取与订单相关的租户资源列表
export const getRenterResByTicketId = (modelData) =>{
    return service({
        url: "renter/res/getRenterResByTicketId",
        method: "post",
        data: modelData,
    });
}
// 获取与订单相关的租户资源列表
export const queryRenterInfo = (modelData) =>{
    return service({
        url: "renter/renter/queryRenterInfo",
        method: "post",
        data: modelData,
    });
}
// 获取与订单相关的租户资源列表
export const updateRenter = (modelData) =>{
    return service({
        url: "renter/renter/updateRenter",
        method: "post",
        data: modelData,
    });
}
// 获取与订单相关的租户资源列表
export const queryRenterResList = (modelData) =>{
    return service({
        url: "renter/res/queryRenterResList",
        method: "post",
        data: modelData,
    });
}