package image

import (
	v1 "ygpt/server/api/v1"

	"github.com/gin-gonic/gin"
)

type ImageRouter struct {
}

func (s *ImageRouter) InitStoImageRouter(Router *gin.RouterGroup) {
	dataRouter := Router.Group("cloud/image/image")
	imageApi := v1.ApiGroupApp.CloudApiGroup.ImageApiGroup.ImageApi
	{

		dataRouter.POST("getImageList", imageApi.GetImageList)           // 获取镜像列表
		dataRouter.POST("getImageByName", imageApi.GetImageByName)       // 获取单个镜像
		dataRouter.DELETE("deleteImage", imageApi.DeleteImage)           // 删除镜像
		dataRouter.DELETE("deleteImageByIds", imageApi.DeleteImageByIds) // 删除镜像
		dataRouter.POST("uploadImage", imageApi.UploadImage)             // 镜像上传
		dataRouter.PUT("updateImageInfo", imageApi.UpdateImageInfo)      // 修改镜像信息
		dataRouter.GET("downloadImage", imageApi.DownloadImage)          // 镜像下载
		dataRouter.POST("betterUpImage", imageApi.BetterUpImage)
		dataRouter.POST("getRecycleImageList", imageApi.GetRecycleImageList) //获取回收站列表
		dataRouter.POST("clearRecycleImage", imageApi.ClearRecycleImage)     // 清除回收镜像
		dataRouter.POST("recoverRecycleImage", imageApi.RecoverRecycleImage) // 恢复回收镜像
		dataRouter.POST("updateImageTag", imageApi.UpdateImageTag)           // 镜像标签修改
	}
}
