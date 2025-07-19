package renter

import (
	v1 "ygpt/server/api/v1"

	"github.com/gin-gonic/gin"
)

type RenterResRouter struct {
}

func (s *RenterResRouter) InitRenterResRouter(Router *gin.RouterGroup) {
	dataRouter := Router.Group("renter/res")
	supplyInfoApi := v1.ApiGroupApp.YunGuanApiGroup.RenterResApi
	{
		dataRouter.POST("addRenterRes", supplyInfoApi.AddRenterRes)                         // 添加数据
		dataRouter.PUT("updateRenterRes", supplyInfoApi.UpdateRenterRes)                    // 添加数据
		dataRouter.DELETE("deleteRenterRes", supplyInfoApi.DeleteRenterRes)                 // 删除数据
		dataRouter.POST("queryRenterRes", supplyInfoApi.QueryRenterRes)                     // 获取数据列表
		dataRouter.POST("queryRenterResInfo", supplyInfoApi.QueryRenterResInfo)             //获取某个资源信息
		dataRouter.POST("queryRenterResCount", supplyInfoApi.QueryRenterResCountByTicketId) // 获取某个故障工单资源数量
		//dataRouter.POST("queryAllRenterRes", supplyInfoApi.QueryAllRenterRes) // 获取所有数据列表
		dataRouter.PUT("bindResInfo", supplyInfoApi.BindRenterRes)                      // 绑定资源
		dataRouter.PUT("releaseResInfo", supplyInfoApi.ReleaseRenterRes)                // 释放资源
		dataRouter.GET("queryRenterFinance", supplyInfoApi.QueryRenterFinance)          // 获取财务明细
		dataRouter.POST("getRenterResByTicketId", supplyInfoApi.GetRenterResByTicketId) // 获取与工单相关的租户资源列表

		dataRouter.POST("queryRenterResList", supplyInfoApi.QueryRenterResList) // 获取与工单相关的租户资源列表
	}
}
