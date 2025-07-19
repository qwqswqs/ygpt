package renter

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"time"
	"ygpt/server/global"
	"ygpt/server/model/common/request"
	"ygpt/server/model/common/response"
	"ygpt/server/model/yunguan/renter"
	"ygpt/server/utils"
)

type RenterRecordApi struct {
}

// 添加
func (p *RenterRecordApi) AddRenterRecord(c *gin.Context) {
	var r renter.RenterRecord
	err := c.ShouldBindJSON(&r)
	if err != nil {
		response.FailWithMessage(err, c)
		return
	}
	serviceReturn, err := renterRecordService.AddRenterRecordService(r)
	if err != nil {
		global.GVA_LOG.Error("添加失败!", zap.Error(err))
		response.FailWithDetailed(global.GVA_MODEL{ID: serviceReturn.ID}, "添加失败", c)
		return
	}
	response.OkWithMessage("添加成功", c)
}

// 修改
func (p *RenterRecordApi) UpdateRenterRecord(c *gin.Context) {
	var r renter.RenterRecord
	err := c.ShouldBindJSON(&r)
	if err != nil {
		response.FailWithMessage(err, c)
		return
	}
	err = renterRecordService.UpdateRenterRecordService(r)
	if err != nil {
		global.GVA_LOG.Error("更新失败!", zap.Error(err))
		response.FailWithMessage("更新失败", c)
		return
	}
	response.OkWithDetailed(gin.H{"updateAt": time.Now().Unix()}, "修改成功", c)
}

// 查询
func (p *RenterRecordApi) QueryRenterRecord(c *gin.Context) {
	var pageInfo request.PageInfo
	err := c.ShouldBindJSON(&pageInfo)
	if err != nil {
		response.FailWithMessage(err, c)
		return
	}

	err = utils.Verify(pageInfo, utils.PageInfoVerify)
	if err != nil {
		response.FailWithMessage(err, c)
		return
	}
	list, total, err := renterRecordService.GetRenterRecordListService(pageInfo)
	if err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage(err, c)
		return
	}
	response.OkWithDetailed(response.PageResult{
		List:     list,
		Total:    total,
		Page:     pageInfo.Page,
		PageSize: pageInfo.PageSize,
	}, "获取成功", c)
}

// 删除
func (p *RenterRecordApi) DeleteRenterRecord(c *gin.Context) {
	var r renter.RenterRecord
	err := c.ShouldBindJSON(&r)
	if err != nil {
		response.FailWithMessage(err, c)
		return
	}
	err = renterRecordService.DeleteRenterRecordService(r.ID)
	if err != nil {
		global.GVA_LOG.Error("删除失败!", zap.Error(err))
		response.FailWithMessage("删除失败", c)
		return
	}
	response.OkWithDetailed(gin.H{"deleteAt": time.Now().Unix()}, "删除成功", c)
}
