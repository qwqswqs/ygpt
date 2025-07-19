package main

import (
	"encoding/json"
	"fmt"
	"testing"
	"ygpt/server/global"
	"ygpt/server/service"
)

var systemConfigService = service.ServiceGroupApp.YunGuanServiceGroup.SystemConfigService

func TestGet(t *testing.T) {

	config, _ := systemConfigService.GetSystemConfigByKeyService("image")
	var imageConfig global.MinioData
	fmt.Println(config.Value)
	json.Unmarshal([]byte(config.Value), imageConfig)
	fmt.Println(imageConfig)
}
