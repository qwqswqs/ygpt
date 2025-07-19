package system

type SystemRouter struct {
	SystemLogRouter
	SystemToolRouter
	SystemTypeRouter
	SystemTicketRouter
	SystemConfigRouter
	SystemCacheRouter
	SystemAlertRouter
	SystemSoftwareRouter
	systemOperateRouter
	SystemPlatformConnectRouter
}
