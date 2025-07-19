package cloudpods

import (
	"github.com/gin-gonic/gin"
	"ygpt/server/model/common/response"
	"yunion.io/x/jsonutils"
)

type ImageApi struct {
}

type getSystemImageReq struct {
	Type string `json:"type"`
}

func (i *ImageApi) GetSystemImage(c *gin.Context) {
	//var r getSystemImageReq
	//if err := c.ShouldBindJSON(&r); err != nil {
	//	response.FailWithMessage("参数失败", c)
	//	return
	//}

	list, total, err := ImageService.GetSystemImage()
	if err != nil {
		response.FailWithMessage("获取失败:"+err.Error(), c)
		return
	}

	response.OkWithDetailed(response.PageResult{
		List:  jsonutils.Marshal(list).PrettyString(),
		Total: total,
	}, "获取成功", c)
}

func (i *ImageApi) GetDockerImage(c *gin.Context) {
	list, total, err := ImageService.GetDockerImage()
	if err != nil {
		response.FailWithMessage("获取失败:"+err.Error(), c)
		return
	}

	response.OkWithDetailed(response.PageResult{
		List:  list,
		Total: total,
	}, "获取成功", c)
}
