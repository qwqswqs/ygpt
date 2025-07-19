package system

import "ygpt/server/service"

type SystemApi struct {
	SystemLogApi
	SystemToolApi
	SystemTypeApi
	SystemConfigApi
	SystemCacheApi
	SystemTicketApi
	SystemAlertApi
	SystemSoftwareApi
	SystemOperateApi
	PlatformConnectApi
}

var (
	systemLogService      = service.ServiceGroupApp.YunGuanServiceGroup.SystemLogService
	systemOperateService  = service.ServiceGroupApp.YunGuanServiceGroup.SystemOperateService
	systemToolService     = service.ServiceGroupApp.YunGuanServiceGroup.SystemToolService
	systemTypeService     = service.ServiceGroupApp.YunGuanServiceGroup.SystemTypeService
	systemTicketService   = service.ServiceGroupApp.YunGuanServiceGroup.SystemTicketService
	systemConfigService   = service.ServiceGroupApp.YunGuanServiceGroup.SystemConfigService
	systemCacheService    = service.ServiceGroupApp.YunGuanServiceGroup.SystemCacheService
	systemAlertService    = service.ServiceGroupApp.YunGuanServiceGroup.SystemAlertService
	systemSoftwareService = service.ServiceGroupApp.YunGuanServiceGroup.SystemSoftwareService
)
