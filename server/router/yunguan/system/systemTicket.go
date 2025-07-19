package system

import (
	v1 "ygpt/server/api/v1"

	"github.com/gin-gonic/gin"
)

type SystemTicketRouter struct{}

func (s *SystemTicketRouter) InitSystemTicketRouter(Router *gin.RouterGroup) {
	dataRouter := Router.Group("sys/ticket")
	modelApi := v1.ApiGroupApp.YunGuanApiGroup.SystemTicketApi
	{
		dataRouter.POST("addSysTicket", modelApi.AddSystemTicket)         // 添加数据
		dataRouter.PUT("updateSysTicket", modelApi.UpdateSystemTicket)    // 修改数据
		dataRouter.PUT("assignSysTicket", modelApi.AssignSystemTicket)    // 分配数据
		dataRouter.DELETE("deleteSysTicket", modelApi.DeleteSystemTicket) // 删除数据
		dataRouter.POST("getSysTicketList", modelApi.GetSystemTicketList) // 获取数据列表
		dataRouter.POST("uploadContract", modelApi.UploadContract)        // 上传合同
		dataRouter.POST("uploadInvoice", modelApi.UploadInvoice)          // 发票上传
		dataRouter.POST("postQuotation", modelApi.PostQuotation)          //工单报价
		dataRouter.POST("assignResource", modelApi.AssignResource)        //故障工单，资源分配
		dataRouter.POST("getNoAssignRes", modelApi.GetNoAssignRes)        //获取未分配的资源
		dataRouter.POST("getAllRes", modelApi.GetAllRes)                  //获取所有资源
	}
}
