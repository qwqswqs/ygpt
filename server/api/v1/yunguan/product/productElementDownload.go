package product

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"ygpt/server/global"
	"ygpt/server/model/common/response"
)

type ProductElemDownloadApi struct {
}

type getDownloadElemReq struct {
	PrivateIp string `json:"privateIp"`
}

func (p *ProductElemDownloadApi) ListProductElemDownload(c *gin.Context) {
	var r getDownloadElemReq
	err := c.ShouldBindJSON(&r)
	if err != nil {
		response.FailWithMessage("参数错误", c)
		return
	}
	serviceReturn, err := ProductElemDownloadService.ListElementDownload(r.PrivateIp)
	if err != nil {
		global.GVA_LOG.Error("查询失败!", zap.Error(err))
		response.FailWithMessage(err.Error()+"查询失败", c)
		return
	}
	//返回一个数组
	response.OkWithDetailed(serviceReturn, "添加成功", c)
}
