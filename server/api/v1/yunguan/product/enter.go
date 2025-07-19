package product

import "ygpt/server/service"

type ProductApi struct {
	ProductElemApi
	ProductElemDownloadApi
	ProductSupplyApi
	ProductComputingApi
	ProductConfigApi
	ProductCategoryApi
}

var (
	productElemService          = service.ServiceGroupApp.YunGuanServiceGroup.ProductElemService
	ProductElemDownloadService  = service.ServiceGroupApp.YunGuanServiceGroup.ProductElemDownloadService
	productConfigService        = service.ServiceGroupApp.YunGuanServiceGroup.ProductConfigService
	productCategoryService      = service.ServiceGroupApp.YunGuanServiceGroup.ProductCategoryService
	productSupplyService        = service.ServiceGroupApp.YunGuanServiceGroup.ProductSupplyService
	productComputingService     = service.ServiceGroupApp.YunGuanServiceGroup.ProductComputingService
	productElementRecordService = service.ServiceGroupApp.YunGuanServiceGroup.ProductElementRecordService
)
