package network

import (
	"github.com/gin-gonic/gin"
	v1 "ygpt/server/api/v1"
)

type NetworkSecgroupRouter struct {
}

func (s *NetworkVpcRouter) InitNetworkSecgroupRouter(Router *gin.RouterGroup) {
	dataRouter := Router.Group("cloud/network/secgroup")
	secgroupApi := v1.ApiGroupApp.CloudApiGroup.NetApiGroup.NetworkSecgroupApi

	{
		dataRouter.POST("getSecgroupList", secgroupApi.GetSecgroupList)                 // 获取数据列表
		dataRouter.POST("addSecgroup", secgroupApi.AddSecgroup)                         // 新增
		dataRouter.POST("deleteSecgroup", secgroupApi.DeleteSecgroup)                   // 删除
		dataRouter.POST("deleteSecgroupByIds", secgroupApi.DeleteSecgroupByIds)         // 删除
		dataRouter.POST("getSecgroupRuleList", secgroupApi.GetSecgroupRuleList)         // 获取数据列表
		dataRouter.POST("addSecgroupRule", secgroupApi.AddSecgroupRule)                 // 新增
		dataRouter.POST("updateSecgroupRule", secgroupApi.UpdateSecgroupRule)           // 新增
		dataRouter.POST("deleteSecgroupRule", secgroupApi.DeleteSecgroupRule)           // 删除
		dataRouter.POST("deleteSecgroupRuleByIds", secgroupApi.DeleteSecgroupRuleByIds) // 删除

		dataRouter.POST("addSeverSecgroupRule", secgroupApi.AddSeverSecgroupRule)       // 新增
		dataRouter.POST("revokeSeverSecgroupRule", secgroupApi.RevokeSeverSecgroupRule) // 删除
	}
}
