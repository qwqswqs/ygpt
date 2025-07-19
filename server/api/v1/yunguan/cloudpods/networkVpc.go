package cloudpods

import (
	"github.com/gin-gonic/gin"
	"ygpt/server/global"
	"ygpt/server/model/common/response"
	"yunion.io/x/jsonutils"
	"yunion.io/x/onecloud/pkg/mcclient/modules/compute"
)

type VpcApi struct {
}
type getVpcReq struct {
	PageIndex int64 `json:"pageIndex"`
	PageSize  int64 `json:"pageSize"`
}

// 获取VPC列表
func (d *VpcApi) GetVpcList(c *gin.Context) {
	var r getVpcReq
	err := c.ShouldBindJSON(&r)
	if err != nil {
		response.FailWithMessage(err, c)
		return
	}
	s, err := global.GetCloudClientSession()
	if err != nil {
		response.FailWithMessage(err, c)
	}
	params := jsonutils.NewDict()
	params.Set("limit", jsonutils.NewInt(r.PageSize))
	params.Set("offset", jsonutils.NewInt(r.PageSize*(r.PageIndex-1)))
	result, _ := compute.Vpcs.List(s, params)

	response.OkWithDetailed(jsonutils.Marshal(result).PrettyString(), "获取成功", c)
}

//
//// 修改
//func (u *VpcApi) UpdateVpc(c *gin.Context) {
//	var r network.NetworkVpc
//	err := c.ShouldBindJSON(&r)
//	if err != nil {
//		response.FailWithMessage(err, c)
//		return
//	}
//	fmt.Println(r)
//	err = VpcService.UpdateVpcService(r)
//	if err != nil {
//		global.GVA_LOG.Error("更新失败!", zap.Error(err))
//		response.FailWithMessage("更新失败", c)
//		return
//	}
//	response.OkWithDetailed(gin.H{"updateAt": time.Now().Unix()}, "修改成功", c)
//}
//
//// 查询
//func (u *VpcApi) QueryVpc(c *gin.Context) {
//	var pageInfo network.NetworkVpcReq
//	err := c.ShouldBindJSON(&pageInfo)
//	if err != nil {
//		response.FailWithMessage(err, c)
//		return
//	}
//	err = utils.Verify(pageInfo, utils.RegisterVerify)
//	if err != nil {
//		response.FailWithMessage(err, c)
//		return
//	}
//	list, total, err := VpcService.GetVpcListService(pageInfo)
//	if err != nil {
//		global.GVA_LOG.Error("获取失败!", zap.Error(err))
//		response.FailWithMessage(err, c)
//		return
//	}
//	response.OkWithDetailed(response.PageResult{
//		List:     list,
//		Total:    total,
//		Page:     pageInfo.Page,
//		PageSize: pageInfo.PageSize,
//	}, "获取成功", c)
//}
//
//// 删除
//func (u *VpcApi) DeleteVpc(c *gin.Context) {
//	var r network.NetworkVpc
//	err := c.ShouldBindJSON(&r)
//	if err != nil {
//		response.FailWithMessage(err, c)
//		return
//	}
//	err = VpcService.DeleteVpcService(r.ID)
//	if err != nil {
//		global.GVA_LOG.Error("删除失败!", zap.Error(err))
//		response.FailWithMessage("删除失败", c)
//		return
//	}
//	response.OkWithDetailed(gin.H{"deleteAt": time.Now().Unix()}, "删除成功", c)
//}
