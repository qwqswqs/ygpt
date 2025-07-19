package res

import (
	"github.com/gin-gonic/gin"
	v1 "ygpt/server/api/v1"
)

type ResInfoRouter struct {
}

func (s *ResInfoRouter) InitResInfoRouter(Router *gin.RouterGroup) {
	dataRouter := Router.Group("res/info")
	diskApi := v1.ApiGroupApp.YunGuanApiGroup.ResApi.ResInfoApi

	{
		dataRouter.POST("addInfo", diskApi.AddResInfo)                  // 添加数据
		dataRouter.PUT("updateInfo", diskApi.UpdateResInfo)             // 修改数据
		dataRouter.DELETE("deleteInfo", diskApi.DeleteResInfo)          // 删除数据
		dataRouter.POST("queryInfo", diskApi.QueryResInfo)              // 获取数据列表
		dataRouter.POST("queryRenterInfo", diskApi.QueryRenterResInfo)  // 获取租户数据列表
		dataRouter.POST("queryResInfoDetail", diskApi.GetResInfoDetail) //获取租户实例详情
	}
}
