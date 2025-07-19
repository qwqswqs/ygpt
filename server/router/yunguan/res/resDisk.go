package res

import (
	"github.com/gin-gonic/gin"
	v1 "ygpt/server/api/v1"
)

type ResDiskRouter struct {
}

func (s *ResDiskRouter) InitResDiskRouter(Router *gin.RouterGroup) {
	dataRouter := Router.Group("res/disk")
	diskApi := v1.ApiGroupApp.YunGuanApiGroup.ResApi.DiskApi

	{
		dataRouter.POST("addDisk", diskApi.AddDisk)         // 添加数据
		dataRouter.PUT("updateDisk", diskApi.UpdateDisk)    // 修改数据
		dataRouter.DELETE("deleteDisk", diskApi.DeleteDisk) // 删除数据
		dataRouter.POST("queryDisk", diskApi.QueryDisk)     // 获取数据列表
	}
}
