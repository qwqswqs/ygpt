package cloudpods

import (
	"ygpt/server/global"
	"ygpt/server/model/common/response"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type ResInfoApi struct {
}

func (r *ResInfoApi) GetResInfo(c *gin.Context) {
	resInfo, err := ResInfoService.GetResInfo()
	if err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage(err, c)
	}
	response.OkWithDetailed(resInfo, "获取成功", c)
}
