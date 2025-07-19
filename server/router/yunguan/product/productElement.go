package product

import (
	"github.com/gin-gonic/gin"
	v1 "ygpt/server/api/v1"
)

type ProductElementRouter struct {
}

func (s *ProductElementRouter) InitProductElementRouter(Router *gin.RouterGroup) {
	dataRouter := Router.Group("product/elem")
	infoApi := v1.ApiGroupApp.YunGuanApiGroup.ProductElemApi

	{
		dataRouter.POST("upload", infoApi.UploadFile)                       // 上传文件
		dataRouter.PUT("updateElemInfo", infoApi.UpdateInfo)                // 修改数据
		dataRouter.PUT("updateElemFileInfo", infoApi.UpdateFileInfo)        // 修改文件数据
		dataRouter.POST("queryElemInfo", infoApi.QueryInfo)                 // 获取数据列表
		dataRouter.POST("queryUserElemInfo", infoApi.QueryUserInfo)         // 获取租户数据列表
		dataRouter.POST("queryShareElemInfo", infoApi.QueryShareInfo)       // 获取共享数据列表
		dataRouter.POST("queryAllShareElemInfo", infoApi.QueryAllShareInfo) // 获取共享数据列表
		dataRouter.POST("addElemInfo", infoApi.AddInfo)                     // 添加数据
		dataRouter.DELETE("deleteElemInfo", infoApi.DeleteInfo)             // 删除数据
		dataRouter.POST("syncElemInfo", infoApi.SyncInfo)                   // 同步数据
		dataRouter.GET("getAllModels", infoApi.GetAllAvailableModel)        // 获取所有模型
		dataRouter.GET("getAllDataSets", infoApi.GetAllAvailableDataSet)    //获取所有数据集

		dataRouter.POST("createProduct", infoApi.CreateProduct)        //获取所有数据集
		dataRouter.PUT("updateProduct", infoApi.UpdateProduct)         //获取所有数据集
		dataRouter.GET("getEleListInfo", infoApi.GetEleListInfo)       //获取所有数据集
		dataRouter.GET("getUserSupplyList", infoApi.GetUserSupplyList) //获取所有数据集
		dataRouter.DELETE("deleteProduct", infoApi.DeleteProduct)      // 删除product表
		dataRouter.DELETE("deleteProductByIds", infoApi.DeleteProductByIds)
	}
}
