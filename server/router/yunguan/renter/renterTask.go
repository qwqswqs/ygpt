package renter

import (
	"github.com/gin-gonic/gin"
	v1 "ygpt/server/api/v1"
)

type RenterTaskRouter struct {
}

func (s *RenterTaskRouter) InitRenterTaskRouter(Router *gin.RouterGroup) {
	dataRouter := Router.Group("renter/task")
	supplyInfoApi := v1.ApiGroupApp.YunGuanApiGroup.RenterTaskApi

	{
		dataRouter.POST("addRenterTask", supplyInfoApi.AddRenterTask)         // 添加数据
		dataRouter.PUT("updateRenterTask", supplyInfoApi.UpdateRenterTask)    // 添加数据
		dataRouter.DELETE("deleteRenterTask", supplyInfoApi.DeleteRenterTask) // 删除数据
		dataRouter.POST("queryRenterTask", supplyInfoApi.QueryRenterTask)     // 获取数据列表
	}
}
