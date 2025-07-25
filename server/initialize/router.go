package initialize

import (
	"net/http"
	"os"

	"ygpt/server/docs"
	"ygpt/server/global"
	"ygpt/server/middleware"
	"ygpt/server/router"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type justFilesFilesystem struct {
	fs http.FileSystem
}

func (fs justFilesFilesystem) Open(name string) (http.File, error) {
	f, err := fs.fs.Open(name)
	if err != nil {
		return nil, err
	}

	stat, err := f.Stat()
	if stat.IsDir() {
		return nil, os.ErrPermission
	}

	return f, nil
}

// 初始化总路由

func Routers() *gin.Engine {
	Router := gin.New()
	Router.Use(gin.Recovery())
	if gin.Mode() == gin.DebugMode {
		Router.Use(gin.Logger())
	}

	systemRouter := router.RouterGroupApp.System
	exampleRouter := router.RouterGroupApp.Example
	yunRouter := router.RouterGroupApp.YunGuan
	comRouter := router.RouterGroupApp.Compute
	cloudRouter := router.RouterGroupApp.Cloud
	// 如果想要不使用nginx代理前端网页，可以修改 web/.env.production 下的
	// VUE_APP_BASE_API = /
	// VUE_APP_BASE_PATH = http://localhost
	// 然后执行打包命令 npm run build。在打开下面3行注释
	// Router.Static("/favicon.ico", "./dist/favicon.ico")
	// Router.Static("/assets", "./dist/assets")   // dist里面的静态资源
	// Router.StaticFile("/", "./dist/index.html") // 前端网页入口页面

	Router.StaticFS(global.GVA_CONFIG.Local.StorePath, justFilesFilesystem{http.Dir(global.GVA_CONFIG.Local.StorePath)}) // Router.Use(middleware.LoadTls())  // 如果需要使用https 请打开此中间件 然后前往 core/server.go 将启动模式 更变为 Router.RunTLS("端口","你的cre/pem文件","你的key文件")
	// 跨域，如需跨域可以打开下面的注释
	// Router.Use(middleware.Cors()) // 直接放行全部跨域请求
	// Router.Use(middleware.CorsByRules()) // 按照配置的规则放行跨域请求
	// global.GVA_LOG.Info("use middleware cors")
	docs.SwaggerInfo.BasePath = global.GVA_CONFIG.System.RouterPrefix
	Router.GET(global.GVA_CONFIG.System.RouterPrefix+"/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	global.GVA_LOG.Info("register swagger handler")
	// 方便统一添加路由组前缀 多服务器上线使用

	PublicGroup := Router.Group(global.GVA_CONFIG.System.RouterPrefix)
	PrivateGroup := Router.Group(global.GVA_CONFIG.System.RouterPrefix)

	PrivateGroup.Use(middleware.JWTAuth()).Use(middleware.CasbinHandler())

	{
		// 健康监测
		PublicGroup.GET("/health", func(c *gin.Context) {
			c.JSON(http.StatusOK, "ok")
		})
	}
	{
		systemRouter.InitBaseRouter(PublicGroup) // 注册基础功能路由 不做鉴权
		systemRouter.InitInitRouter(PublicGroup) // 自动初始化相关

		comRouter.InitUdpMonitorRouter(PublicGroup) // udp操作实例的api

		//cloudpods相关路由注册
		cloudRouter.InitServersRouter(PublicGroup)
		cloudRouter.InitScalingConfigRouter(PublicGroup)
		cloudRouter.InitScalingGroupRouter(PublicGroup)
		cloudRouter.InitScalingPolicyRouter(PublicGroup)
		cloudRouter.InitStoImageRouter(PublicGroup)
		cloudRouter.InitNetworkIpRouter(PublicGroup)
		cloudRouter.InitNetworkVpcRouter(PublicGroup)
		cloudRouter.InitNetworkEipRouter(PublicGroup)
		cloudRouter.InitNetworkSecgroupRouter(PublicGroup)
		cloudRouter.InitVmEvsRouter(PublicGroup)
		cloudRouter.InitVmNasRouter(PublicGroup)
		cloudRouter.InitContainerEvsRouter(PublicGroup)
		cloudRouter.InitContainerNasRouter(PublicGroup)
		cloudRouter.InitImageReposRouter(PublicGroup)
		cloudRouter.InitStoSnapshotRouter(PublicGroup)
		cloudRouter.InitBaseHostRouter(PublicGroup)
		cloudRouter.InitBaseK8SRouter(PublicGroup)
		cloudRouter.InitBaseDeviceRouter(PublicGroup)
		cloudRouter.InitK8sPodRouter(PublicGroup)
		cloudRouter.InitK8sDeploymentRouter(PublicGroup)
		cloudRouter.InitCloudCommonLogRouter(PublicGroup)

		//云管相关路由注册
		yunRouter.InitLicenseConfigRouter(PublicGroup)
		yunRouter.InitAlertConfigRouter(PublicGroup)

		yunRouter.InitSystemLogRouter(PublicGroup)
		yunRouter.InitSystemOperateRouter(PublicGroup)
		yunRouter.InitSystemToolRouter(PublicGroup)
		yunRouter.InitSystemTypeRouter(PublicGroup)
		yunRouter.InitSystemTicketRouter(PublicGroup)
		yunRouter.InitSystemCacheRouter(PublicGroup)
		yunRouter.InitSystemConfigRouter(PublicGroup)
		yunRouter.InitSystemAlertRouter(PublicGroup)
		yunRouter.InitSystemSoftwareRouter(PublicGroup)

		yunRouter.InitResDiskRouter(PublicGroup)
		yunRouter.InitResDeviceRouter(PublicGroup)
		yunRouter.InitResInfoRouter(PublicGroup)

		yunRouter.InitRenterRouter(PublicGroup)
		yunRouter.InitRenterResRouter(PublicGroup)
		yunRouter.InitRenterTaskRouter(PublicGroup)

		yunRouter.InitBackupImageRouter(PublicGroup)
		yunRouter.InitBackRouter(PublicGroup)
		yunRouter.InitBackupSnapshotRouter(PublicGroup)

		yunRouter.InitNetworkIpRouter(PublicGroup)
		yunRouter.InitNetworkNatRouter(PublicGroup)
		yunRouter.InitNetworkSubRouter(PublicGroup)
		yunRouter.InitNetworkVpcRouter(PublicGroup)

		yunRouter.InitProductConfigRouter(PublicGroup)
		yunRouter.InitProductElementRouter(PublicGroup)
		yunRouter.InitProductElementDownloadRouter(PublicGroup)
		yunRouter.InitProductSupplyRouter(PublicGroup)
		yunRouter.InitProductComputingRouter(PublicGroup)
		yunRouter.InitProductCategoryRouter(PublicGroup)

		yunRouter.InitSystemImageRouter(PublicGroup)
		yunRouter.InitDockerImageRouter(PublicGroup)
		yunRouter.InitContainerRouter(PublicGroup)
		yunRouter.InitVirtualMachineRouter(PublicGroup)
		yunRouter.InitBareMetalRouter(PublicGroup)
		yunRouter.InitJumpRouter(PublicGroup)
		yunRouter.InitCloudpodsResInfoRouter(PublicGroup)
	}

	{
		systemRouter.InitApiRouter(PublicGroup, PublicGroup)       // 注册功能api路由
		systemRouter.InitJwtRouter(PublicGroup)                    // jwt相关路由
		systemRouter.InitUserRouter(PublicGroup)                   // 注册用户路由
		systemRouter.InitMenuRouter(PublicGroup)                   // 注册menu路由
		systemRouter.InitSystemRouter(PublicGroup)                 // system相关路由
		systemRouter.InitCasbinRouter(PublicGroup)                 // 权限相关路由
		systemRouter.InitAutoCodeRouter(PublicGroup, PublicGroup)  // 创建自动化代码
		systemRouter.InitAuthorityRouter(PublicGroup)              // 注册角色路由
		systemRouter.InitSysDictionaryRouter(PublicGroup)          // 字典管理
		systemRouter.InitAutoCodeHistoryRouter(PublicGroup)        // 自动化代码历史
		systemRouter.InitSysOperationRecordRouter(PublicGroup)     // 操作记录
		systemRouter.InitSysDictionaryDetailRouter(PublicGroup)    // 字典详情管理
		systemRouter.InitAuthorityBtnRouterRouter(PublicGroup)     // 按钮权限管理
		systemRouter.InitSysExportTemplateRouter(PublicGroup)      // 导出模板
		exampleRouter.InitCustomerRouter(PublicGroup)              // 客户路由
		exampleRouter.InitFileUploadAndDownloadRouter(PublicGroup) // 文件上传下载功能路由

	}
	//插件路由安装
	InstallPlugin(PrivateGroup, PublicGroup, Router)

	// 注册业务路由
	initBizRouter(PrivateGroup, PublicGroup)

	global.GVA_ROUTERS = Router.Routes()

	global.GVA_LOG.Info("router register success")
	return Router
}
