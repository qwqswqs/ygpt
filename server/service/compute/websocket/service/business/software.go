package business

import (
	"encoding/json"
	"go.uber.org/zap"
	"time"
	"ygpt/server/global"
	"ygpt/server/model/yunguan/system"
	"ygpt/server/service/compute/websocket/protocol"
)

type SoftwareService struct{}

func (s *SoftwareService) QuerySoftwareVersion() error {
	queryMsg := protocol.SoftVersionRequest{}
	queryMsg.ID = global.GVA_WS.Counter.GetID()
	queryMsg.Command = protocol.SoftVersionCmd
	queryMsg.NodeCode = global.GVA_CONFIG.System.NodeCode
	return global.GVA_WS.SendingWithRetry(queryMsg)
}

func (s *SoftwareService) HandleQuerySoftwareVersionResult(message *protocol.SoftVersionResponse) error {
	softWares := message.SoftWaresJson
	for _, v := range softWares {
		var software system.SystemSoftware
		// 细节Json解析到 map 中
		var result map[string]interface{}
		err := json.Unmarshal([]byte(v), &result)
		if err != nil {
			global.GVA_LOG.Error("Error decoding JSON:", zap.Error(err))
			continue // 继续处理下一个软件
		}
		PlatSoftware, ok := result["id"].(float64)
		if !ok {
			global.GVA_LOG.Error("Invalid ID format for software")
			continue // 继续处理下一个软件
		}
		// 使用 FirstOrCreate 查找或创建记录
		err = global.GVA_DB.Where("re_software_id = ?", int(PlatSoftware)).FirstOrCreate(&software).Error
		if err != nil {
			global.GVA_LOG.Error("Failed to find or create software record:", zap.Error(err))
			continue // 继续处理下一个软件
		}
		// 更新软件信息
		software.Name = result["name"].(string)
		//software.Type = result["category"].(string)
		software.Version = result["version"].(string)
		software.Url = result["url"].(string)
		software.Size = result["size"].(string)
		lastUpdateTimeStr, ok := result["lastUpdate"].(string)
		lastUpdateTime, err := time.Parse(time.RFC3339, lastUpdateTimeStr)
		if ok || err == nil {
			software.LastUpdate = lastUpdateTime
		} else {
			global.GVA_LOG.Error("Invalid last update time format")
		}
		software.ReSoftwareID = int(PlatSoftware)
		// 保存更新后的软件信息
		err = global.GVA_DB.Save(&software).Error
		if err != nil {
			global.GVA_LOG.Error("Failed to save software record:", zap.Error(err))
		}
	}
	return nil
}
func (s *SoftwareService) SoftDownload(SoftWareID int) error {
	softDownloadMsg := protocol.SoftDownloadRequest{}
	softDownloadMsg.ID = global.GVA_WS.Counter.GetID()
	softDownloadMsg.Command = protocol.SoftDownloadCmd
	//softDownloadMsg.AuthCode = global.GVA_CONFIG.System.AuthCode
	//softDownloadMsg.SoftWareID = SoftWareID
	return global.GVA_WS.SendingWithRetry(softDownloadMsg)

}

func (s *SoftwareService) HandleSoftDownloadResult(message *protocol.SoftDownloadResponse) error {
	// TODO:写入数据库
	return nil
}
