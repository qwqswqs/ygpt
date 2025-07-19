package backup

import (
	"github.com/gin-gonic/gin"
	v1 "ygpt/server/api/v1"
)

type ImageRouter struct {
}

func (s *ImageRouter) InitBackupImageRouter(Router *gin.RouterGroup) {
	dataRouter := Router.Group("backup/image")
	imageApi := v1.ApiGroupApp.YunGuanApiGroup.BackupApi.ImageApi
	{
		dataRouter.POST("addBackImage", imageApi.AddBackImage)           // 添加数据
		dataRouter.POST("addBackImageFile", imageApi.AddBackImageFile)   // 添加数据
		dataRouter.PUT("updateBackImage", imageApi.UpdateBackImage)      // 修改数据
		dataRouter.DELETE("deleteBackImage", imageApi.DeleteBackImage)   // 删除数据
		dataRouter.POST("queryBackImage", imageApi.QueryBackImage)       // 获取数据列表
		dataRouter.POST("queryBackImageAll", imageApi.QueryBackImageAll) // 获取数据列表
		dataRouter.GET("downloadImage", imageApi.DownloadBackImage)      // 下载镜像
	}
}
