package task

import (
	"encoding/json"
	"ygpt/server/global"
	"ygpt/server/service/compute/websocket/service"
	"ygpt/server/service/yunguan/cloudpods"
)

func PushResInfo() {
	// 向算力平台推送资源信息
	resInfoService := cloudpods.ResInfoService{}
	info, err := resInfoService.GetResInfo()
	if err != nil {
		global.GVA_LOG.Error("向算力平台推送资源信息失败:" + err.Error())
		return
	}

	infoByte, err := json.Marshal(info)
	if err != nil {
		global.GVA_LOG.Error("向算力平台推送资源信息失败:" + err.Error())
		return
	}

	infoMap := make(map[string]interface{})
	err = json.Unmarshal(infoByte, &infoMap)
	if err != nil {
		global.GVA_LOG.Error("向算力平台推送资源信息失败:" + err.Error())
		return
	}

	infoMap["nodeType"] = global.GVA_CONFIG.System.NodeType
	infoByte, err = json.Marshal(infoMap)
	err = service.SuanLiServiceGroupApp.NodeStaticsUpload(string(infoByte))
	if err != nil {
		global.GVA_LOG.Error("向算力平台推送资源信息失败:" + err.Error())
		return
	}
	global.GVA_LOG.Info("向算力平台推送资源信息成功")
}
