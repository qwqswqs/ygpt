package backup

import "ygpt/server/service"

type BackupApi struct {
	SnapshotApi
	BackApi
	ImageApi
}

var (
	BackService     = service.ServiceGroupApp.YunGuanServiceGroup.BackupServiceGroup.BackService
	ImageService    = service.ServiceGroupApp.YunGuanServiceGroup.BackupServiceGroup.ImageService
	SnapshotService = service.ServiceGroupApp.YunGuanServiceGroup.BackupServiceGroup.SnapshotService
)
