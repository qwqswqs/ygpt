package backup

import (
	"github.com/gin-gonic/gin"
	v1 "ygpt/server/api/v1"
)

type BackRouter struct {
}

func (s *BackRouter) InitBackRouter(Router *gin.RouterGroup) {
	dataRouter := Router.Group("backup/back")
	backApi := v1.ApiGroupApp.YunGuanApiGroup.BackupApi.BackApi

	{
		dataRouter.POST("addBack", backApi.AddBack)             // 添加数据
		dataRouter.PUT("updateBack", backApi.UpdateBack)        // 修改数据
		dataRouter.DELETE("deleteBack", backApi.DeleteBack)     // 删除数据
		dataRouter.POST("queryBack", backApi.QueryBack)         // 获取数据列表
		dataRouter.POST("queryUserBack", backApi.QueryUserBack) // 获取租户数据列表
	}
}
