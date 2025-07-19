package service

import (
	"ygpt/server/service/compute/websocket/service/business"
	"ygpt/server/service/compute/websocket/service/maintenance"
	"ygpt/server/service/compute/websocket/service/system"
)

var SuanLiServiceGroupApp = new(SuanLiServiceGroup)

type SuanLiServiceGroup struct {
	business.SoftwareService
	business.TenantService
	business.ResourceService
	business.SpecUploadService
	business.ElementService
	business.TicketService
	business.SupplyDemandService
	maintenance.MonitorService
	system.LoginService
	system.ModifyService
}
