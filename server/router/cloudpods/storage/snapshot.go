package storage

import (
	"github.com/gin-gonic/gin"
	v1 "ygpt/server/api/v1"
)

type StoSnapshotRouter struct {
}

func (s *StoSnapshotRouter) InitStoSnapshotRouter(Router *gin.RouterGroup) {
	dataRouter := Router.Group("cloud/storage/snapshot")
	vpcApi := v1.ApiGroupApp.CloudApiGroup.StorageApiGroup.StoSnapshotApi

	{
		dataRouter.POST("getSnapshotList", vpcApi.GetStoSnapshotList)      // 获取数据列表
		dataRouter.POST("snapshotReset", vpcApi.SnapshotReset)             // 回滚主机
		dataRouter.POST("deleteSnapshot", vpcApi.DeleteSnapshot)           // 删除主机快照
		dataRouter.POST("deleteSnapshotByIds", vpcApi.DeleteSnapshotByIds) // 删除主机快照
	}
}
