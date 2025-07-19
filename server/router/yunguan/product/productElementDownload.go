package product

import (
	"github.com/gin-gonic/gin"
	v1 "ygpt/server/api/v1"
)

type ProductElementDownloadRouter struct {
}

func (s *ProductElementRouter) InitProductElementDownloadRouter(Router *gin.RouterGroup) {
	dataRouter := Router.Group("product/elemDownload")
	infoApi := v1.ApiGroupApp.YunGuanApiGroup.ProductElemDownloadApi
	{
		dataRouter.POST("listDownloadedElement", infoApi.ListProductElemDownload) //获取任务列表
	}
}
