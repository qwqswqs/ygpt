package cloudpods

import (
	"github.com/gin-gonic/gin"
	"net/url"
	"ygpt/server/global"
	"ygpt/server/model/common/response"
)

type JumpApi struct {
}

func (j *JumpApi) Jump(c *gin.Context) {
	jumpUrl := global.GVA_CONFIG.Cloudpods.JumpUrl

	if jumpUrl == "" {
		response.FailWithMessage("获取失败:跳转url为空", c)
		return
	}

	parsedURL, err := url.Parse(jumpUrl)
	if err != nil || parsedURL.Scheme == "" || parsedURL.Host == "" {
		response.FailWithMessage("获取失败:跳转url不合法", c)
		return
	}

	response.OkWithDetailed(parsedURL.String(), "获取成功", c)
}
