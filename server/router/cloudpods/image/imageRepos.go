package image

import (
	"github.com/gin-gonic/gin"
	v1 "ygpt/server/api/v1"
)

type ImageReposRouter struct {
}

func (s *ImageReposRouter) InitImageReposRouter(Router *gin.RouterGroup) {
	dataRouter := Router.Group("cloud/image/imageRepos")
	imageApi := v1.ApiGroupApp.CloudApiGroup.ImageApiGroup.ImageReposApi
	{

		dataRouter.POST("getImageReposList", imageApi.GetImageReposList)           // 获取镜像仓库列表
		dataRouter.POST("addImageRepos", imageApi.AddImageRepos)                   // 新增镜像仓库
		dataRouter.POST("deleteImageRepos", imageApi.DeleteImageRepos)             // 删除镜像仓库
		dataRouter.POST("deleteImageReposByIds", imageApi.DeleteImageReposByIds)   // 删除镜像仓库
		dataRouter.POST("getImageReposImageList", imageApi.GetImageReposImageList) // 获取镜像仓库的镜像列表
	}
}
