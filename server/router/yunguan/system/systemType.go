package system

import (
	"github.com/gin-gonic/gin"
	v1 "ygpt/server/api/v1"
)

type SystemTypeRouter struct{}

func (s *SystemTypeRouter) InitSystemTypeRouter(Router *gin.RouterGroup) {
	dataRouter := Router.Group("sys/type")
	typeRouter := Router.Group("sysType")
	modelApi := v1.ApiGroupApp.YunGuanApiGroup.SystemTypeApi

	supplyInfoApi := v1.ApiGroupApp.YunGuanApiGroup.ProductSupplyApi
	{
		dataRouter.POST("addSysType", modelApi.AddSystemType)         // 添加数据
		dataRouter.PUT("updateSysType", modelApi.UpdateSystemType)    // 添加数据
		dataRouter.DELETE("deleteSysType", modelApi.DeleteSystemType) // 删除数据
		dataRouter.POST("getSysTypeList", modelApi.GetSystemTypeList) // 获取数据列表

		typeRouter.GET("getCategoryList", supplyInfoApi.GetTypeNameList)
		typeRouter.GET("getElemSysTypeList", supplyInfoApi.GetElemSysTypeList)
		typeRouter.PUT("updateCategory", supplyInfoApi.UpdateCategory)
		typeRouter.POST("addCategory", supplyInfoApi.AddCategory)
		typeRouter.DELETE("deleteCategory", supplyInfoApi.DeleteCategory)           // 删除数据
		typeRouter.DELETE("deleteCategoryByIds", supplyInfoApi.DeleteCategoryByIds) // 批量删除sysType表

	}
}
