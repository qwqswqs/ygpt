package main

import (
	"go.uber.org/zap"
	"ygpt/server/core"
	"ygpt/server/global"
	"ygpt/server/initialize"
	"ygpt/server/model/yunguan/backup"
	"ygpt/server/model/yunguan/config"
	"ygpt/server/model/yunguan/k8s"
	"ygpt/server/model/yunguan/network"
	"ygpt/server/model/yunguan/product"
	"ygpt/server/model/yunguan/renter"
	"ygpt/server/model/yunguan/res"
	"ygpt/server/model/yunguan/system"
	"ygpt/server/service/compute/udp/server"
	"ygpt/server/service/compute/websocket"
	"ygpt/server/service/compute/websocket/service"
	"ygpt/server/task"
	//"ygpt/server/service/compute/websocket/service"
)

//go:generate go env -w GO111MODULE=on
//go:generate go env -w GOPROXY=https://goproxy.cn,direct
//go:generate go mod tidy
//go:generate go mod download

// @title                       Gin-Vue-Admin Swagger API接口文档
// @version                     v2.7.5
// @description                 使用gin+vue进行极速开发的全栈开发基础平台
// @securityDefinitions.apikey  ApiKeyAuth
// @in                          header
// @name                        x-token
// @BasePath                    /
func main() {
	global.GVA_VP, global.GVA_CONFIG_PATH = core.Viper() // 初始化Viper
	initialize.OtherInit()
	global.GVA_LOG = core.Zap() // 初始化zap日志库
	zap.ReplaceGlobals(global.GVA_LOG)
	global.GVA_DB = initialize.Gorm() // gorm连接数据库
	initialize.DBList()
	initialize.Timer()
	if global.GVA_DB != nil {
		initialize.RegisterTables() // 初始化表
		// 程序结束前关闭数据库链接
		db, _ := global.GVA_DB.DB()
		defer db.Close()
	}
	system.Migrate()
	res.Migrate()
	renter.Migrate()
	backup.Migrate()
	network.Migrate()
	product.Migrate()
	config.Migrate()
	k8s.Migrate()
	go func() {
		global.GVA_WS = websocket.InitWebsocketClient(global.GVA_CONFIG.System.SlEndpoint) //初始化websocket客户端
		if global.GVA_WS != nil {
			//尝试登录平台
			err := service.SuanLiServiceGroupApp.Login()
			if err != nil {
				global.GVA_LOG.Error("登录平台失败", zap.Error(err))
			} else {
				task.PushResInfo() // 定时任务方法定在task文件包中
			}
		} else {
			global.GVA_LOG.Error("websocket客户端初始化失败")
		}
		websocket.KeepLogin(global.GVA_WS)
	}()

	go func() {
		err := server.WsServer.NewWebsocketServer()
		if err != nil {
			global.GVA_LOG.Error("websocket服务端初始化失败", zap.Error(err))
		}
		global.GVA_LOG.Error("websocket服务端异常结束")
	}()
	core.RunWindowsServer()
}
