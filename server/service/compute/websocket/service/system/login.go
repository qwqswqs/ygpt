package system

import (
	"fmt"
	"go.uber.org/zap"
	"gopkg.in/yaml.v3"
	"os"
	"time"
	"ygpt/server/global"
	"ygpt/server/service/compute/websocket/protocol"
)

type LoginService struct{}

func (l *LoginService) Login() error {
	login := protocol.Login{}
	login.Command = protocol.LoginCmd
	login.ID = global.GVA_WS.Counter.GetID()
	login.Payload.NodeCode = global.GVA_CONFIG.System.NodeCode
	login.Payload.AuthCode = global.GVA_CONFIG.System.AuthCode
	login.Payload.CloudType = 2
	//var sysConfig system.SystemConfig
	//if err := global.GVA_DB.Where("key", "instanceEnv").First(&sysConfig).Error; err != nil {
	//	global.GVA_LOG.Error("获取实例环境失败!", zap.Any("err", err))
	//} else {
	//	login.Payload.InstanceEnvJson = sysConfig.Value
	//}

	//先创建通道，再发送消息，等待接收结果
	resultCh := make(chan any, 1)
	global.AddTransferChannel(int64(login.ID), resultCh)
	defer func() {
		global.DeleteTransferChannel(int64(login.ID))
	}()

	if err := global.GVA_WS.SendingWithRetry(login); err != nil {
		global.GVA_LOG.Error("发送登录消息失败!", zap.Any("err", err))
		return err
	}

	deadLine := time.After(time.Second * 6)
	var result *protocol.LoginResponse
	select {
	case <-deadLine:
		global.GVA_LOG.Error("login timeout")
		return fmt.Errorf("login timeout")
	case r := <-resultCh:
		err, ok := r.(error)
		if ok {
			global.GVA_LOG.Error(err.Error())
			return err
		}
		result, ok = r.(*protocol.LoginResponse)
		if !ok {
			global.GVA_LOG.Error("error type of LoginResponse receive from task channel")
			return fmt.Errorf("error type of LoginResponse receive from task channel")
		}
		if result.RetCode != 0 {
			global.GVA_LOG.Error(fmt.Sprintf("Login result: %+v", result))
			return fmt.Errorf("login result: %+v", result.RetMsg)
		}
		return nil
	}
}

func (l *LoginService) HandleLoginResult(message protocol.LoginResponse) error {
	if message.RetCode == 0 {
		err := global.GVA_VP.ReadInConfig()
		if err != nil {
			global.GVA_LOG.Error("读取配置文件失败!", zap.Any("err", err))
		}
		global.GVA_CONFIG.SLMinio.EndPoint = message.Payload.MinioHost + ":" + message.Payload.MinioPort
		global.GVA_CONFIG.SLMinio.AccessKey = message.Payload.AccessKey
		global.GVA_CONFIG.SLMinio.SecretKey = message.Payload.SecretKey
		global.GVA_CONFIG.SLMinio.UseSSL = message.Payload.UseSSL
		global.GVA_CONFIG.Minio.BucketName = message.Payload.UploadBucket
		//
		//global.GVA_VP.Set("sl-minio.EndPoint", message.Payload.MinioHost+":"+message.Payload.MinioPort)
		//global.GVA_VP.Set("sl-minio.AccessKey", message.Payload.AccessKey)
		//global.GVA_VP.Set("sl-minio.SecretKey", message.Payload.SecretKey)
		//global.GVA_VP.Set("sl-minio.UseSSL", message.Payload.UseSSL)
		//global.GVA_VP.Set("Minio.BucketName", message.Payload.UploadBucket)

		//// 将更新保存到文件
		//err = global.GVA_VP.WriteConfig()
		//if err != nil {
		//	fmt.Printf("Error writing config file, %s", err)
		//}

		// 将结构体转换为YAML格式
		data, err := yaml.Marshal(global.GVA_CONFIG)
		if err != nil {
			return fmt.Errorf("无法将配置转换为YAML格式: %v", err)
		}

		// 将YAML格式写入配置文件
		err = os.WriteFile(global.GVA_CONFIG_PATH, data, 0644)
		if err != nil {
			return fmt.Errorf("无法写入配置文件: %v", err)
		}
	} else {
		global.GVA_LOG.Error(fmt.Sprintf("Login result: %+v", message))
	}
	return nil
}

func (l *LoginService) Logout() error {
	logout := protocol.Logout{}
	logout.Command = protocol.LogoutCmd
	logout.ID = global.GVA_WS.Counter.GetID()
	//logout.Payload.AuthCode = global.GVA_CONFIG.System.AuthCode
	return global.GVA_WS.SendingWithRetry(logout)
}

func (l *LoginService) HandleLogoutResult(message *protocol.LogoutResponse) error {
	global.GVA_LOG.Info(fmt.Sprintf("Logout result: %+v", message))
	return nil
}

func (l *LoginService) Keep() error {
	keepMsg := protocol.Keep{}
	keepMsg.ID = global.GVA_WS.Counter.GetID()
	keepMsg.Command = protocol.KeepCmd
	return global.GVA_WS.SendingWithRetry(keepMsg)
}

func (l *LoginService) HandleKeepResult(message protocol.KeepResponse) error {
	return nil
}
