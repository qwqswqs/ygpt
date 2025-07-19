package backup

import (
	"github.com/gin-gonic/gin"
	v1 "ygpt/server/api/v1"
)

type SnapShotRouter struct {
}

func (s *SnapShotRouter) InitBackupSnapshotRouter(Router *gin.RouterGroup) {
	dataRouter := Router.Group("backup/snapshot")
	snapshotApi := v1.ApiGroupApp.YunGuanApiGroup.BackupApi.SnapshotApi

	{
		dataRouter.POST("addSnapshot", snapshotApi.AddSnapshot)             // 添加数据
		dataRouter.PUT("updateSnapshot", snapshotApi.UpdateSnapshot)        // 修改数据
		dataRouter.DELETE("deleteSnapshot", snapshotApi.DeleteSnapshot)     // 删除数据
		dataRouter.POST("querySnapshot", snapshotApi.QuerySnapshot)         // 获取数据列表
		dataRouter.POST("queryUserSnapshot", snapshotApi.QueryUserSnapshot) // 获取数据列表
	}
}
