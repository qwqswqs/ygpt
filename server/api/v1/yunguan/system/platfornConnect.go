package system

import (
	"github.com/gin-gonic/gin"
	"time"
	"ygpt/server/config"
	"ygpt/server/global"
	"ygpt/server/model/common/response"
	"ygpt/server/service/compute/websocket/service"
)

type PlatformConnectApi struct {
}

func (d *PlatformConnectApi) GetIsConnected(c *gin.Context) {
	isConnected := global.GVA_WS.WsCli.GetIsConnected()
	response.OkWithDetailed(gin.H{"isConnected": isConnected}, "获取成功", c)
}

func (d *PlatformConnectApi) UpdateConnectConfigAndConnect(c *gin.Context) {
	var sysConfig config.System
	err := c.ShouldBindJSON(&sysConfig)

	// 断开连接
	global.GVA_WS.WsCli.Disconnect()
	// 等待自动重连,带超时
	dl := 0
	for {
		if global.GVA_WS.WsCli.GetIsConnected() {
			break
		}
		if dl > 10 {
			response.FailWithMessage("连接失败", c)
			return
		}
		time.Sleep(time.Second)
		dl++
	}
	// 系统会自动登录，手动发送登录报文测试是否能成功登陆
	err = service.SuanLiServiceGroupApp.Login()
	if err != nil {
		response.FailWithMessage("登录失败"+err.Error(), c)
	} else {
		response.OkWithMessage("登录成功", c)
	}
}
