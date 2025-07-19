package system

type SystemServiceGroup struct {
	SystemLogService
	SystemToolService
	SystemConfigService
	SystemCacheService
	SystemTicketService
	SystemAlertService
	SystemSoftwareService
	SystemTypeService
	SystemOperateService
}
