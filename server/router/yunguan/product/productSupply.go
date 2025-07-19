package product

import (
	"github.com/gin-gonic/gin"
	v1 "ygpt/server/api/v1"
)

type ProductSupplyRouter struct {
}

func (s *ProductSupplyRouter) InitProductSupplyRouter(Router *gin.RouterGroup) {
	productRouter := Router.Group("/")

	dataRouter := Router.Group("product/supply")

	supplyInfoApi := v1.ApiGroupApp.YunGuanApiGroup.ProductSupplyApi

	{
		productRouter.GET("sysType/getSysTypeList", supplyInfoApi.GetSupplyConfig)         //获取系统类型
		productRouter.GET("infoRegion/getInfoRegionList", supplyInfoApi.GetSupplyCityList) //获取区域信息
	}
	{
		dataRouter.POST("addSupplyInfo", supplyInfoApi.AddSupplyInfo)            // 添加数据
		dataRouter.PUT("updateSupplyInfo", supplyInfoApi.UpdateSupplyInfo)       // 修改数据
		dataRouter.DELETE("deleteSupplyInfo", supplyInfoApi.DeleteSupplyInfo)    // 删除数据
		dataRouter.POST("querySupplyInfo", supplyInfoApi.QuerySupplyInfo)        // 获取数据列表
		dataRouter.POST("syncSupplyInfo", supplyInfoApi.SyncSupplyDemandInfo)    // 同步单个数据
		dataRouter.POST("createProduct", supplyInfoApi.CreateSupplyInfo)         // 创建供需
		dataRouter.DELETE("deleteProduct", supplyInfoApi.DeleteSupplyInfo)       //删除供需
		dataRouter.DELETE("deleteProductByIds", supplyInfoApi.DeleteSupplyByIds) //删除供需
		dataRouter.GET("getProductList2", supplyInfoApi.GetSupplyList)           //获取供需
		dataRouter.PUT("updateProduct", supplyInfoApi.UpdateSupplyInfo)          //供需更新
	}
}
