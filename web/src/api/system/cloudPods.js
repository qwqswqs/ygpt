import service from "@/utils/request";



/**========================================================================
 *                           CloudPods相关APi
 *========================================================================**/


export const getCloudPodsConsoleUrl = ()=>{
    return service({
        url: '/cloudpods/jump',
        method: 'get',
    })
}