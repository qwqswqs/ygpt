package renter

import (
	"github.com/gin-gonic/gin"
	v1 "ygpt/server/api/v1"
)

type RenterRouter struct {
}

func (s *RenterRouter) InitRenterRouter(Router *gin.RouterGroup) {
	dataRouter := Router.Group("renter/renter")
	supplyInfoApi := v1.ApiGroupApp.YunGuanApiGroup.RenterApi

	{
		dataRouter.POST("addRenter", supplyInfoApi.AddRenter)               // 添加数据
		dataRouter.POST("getRenterList", supplyInfoApi.GetUserList)         // 添加数据
		dataRouter.POST("getAllRenterList", supplyInfoApi.GetAllRenterList) // 添加数据
		dataRouter.POST("getAllUserList", supplyInfoApi.GetAllUserList)     // 添加数据
		dataRouter.DELETE("deleteRenter", supplyInfoApi.DeleteRenter)       // 删除数据
		dataRouter.POST("queryResRenter", supplyInfoApi.QueryResRenter)     // 获取数据列表
		dataRouter.POST("queryRenter", supplyInfoApi.QueryRenter)           // 获取数据列表

		dataRouter.POST("queryRenterInfo", supplyInfoApi.QueryRenterInfo) // 获取数据列表
		dataRouter.POST("updateRenter", supplyInfoApi.UpdateRenter)       // 获取数据列表
	}
}
