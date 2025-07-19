package system

import (
	"fmt"
	"time"
	"ygpt/server/global"
	"ygpt/server/model/common/response"
	"ygpt/server/model/yunguan/system"
	"ygpt/server/utils"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type SystemAlertApi struct {
}

// 添加系统报警
func (d *SystemAlertApi) AddSystemAlert(c *gin.Context) {
	var r system.SystemAlert
	err := c.ShouldBindJSON(&r)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	r.AlertTime = time.Now()
	serviceReturn, err := systemAlertService.AddSystemAlertService(r)
	if err != nil {
		global.GVA_LOG.Error("添加失败!", zap.Error(err))
		response.FailWithDetailed(global.GVA_MODEL{ID: serviceReturn.ID}, "添加失败", c)
		return
	}
	response.OkWithMessage("添加成功", c)
}

// 获取系统报警列表
func (b *SystemAlertApi) GetSystemAlertList(c *gin.Context) {
	var dataReq struct {
		PageSize     int        `json:"pageSize"`
		PageIndex    int        `json:"pageIndex"`
		ResourceType int        `json:"resourceType"`
		AlertType    int        `json:"alertType"`
		AlertTime    *time.Time `json:"alertTime"`
		StartTime    *time.Time `json:"startTime"`
		EndTime      *time.Time `json:"endTime"`
		IsRenterData bool       `json:"isRenterData"`
	}
	err := c.ShouldBindJSON(&dataReq)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	var dataList []system.SystemAlert
	var total int64
	db := global.GVA_DB.Model(&system.SystemAlert{})

	if dataReq.IsRenterData == true {
		id := utils.GetUserID(c)
		fmt.Printf("id:%v", id)
		db.Where("is_renter_data= ?", dataReq.IsRenterData)
		db.Where("renter_id=?", utils.GetUserID(c))
	}
	if dataReq.ResourceType != 0 {
		db.Where("resource_type=?", dataReq.ResourceType)
	}
	if dataReq.AlertType != 0 {
		db.Where("alert_type=?", dataReq.AlertType)
	}
	if dataReq.StartTime != nil {
		db.Where("alert_time > ?", dataReq.StartTime)
	}
	if dataReq.EndTime != nil {
		db.Where("alert_time < ?", dataReq.EndTime)
	}
	err = db.Count(&total).Error
	if err != nil {
		return
	}
	err = db.Order("alert_time DESC").Limit(dataReq.PageSize).Offset(dataReq.PageSize * (dataReq.PageIndex - 1)).Find(&dataList).Error
	if err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage("获取失败", c)
		return
	}
	response.OkWithDetailed(response.PageResult{
		List:     dataList,
		Total:    total,
		Page:     dataReq.PageIndex,
		PageSize: dataReq.PageSize,
	}, "获取成功", c)
}

// 修改查看状态
func (b *SystemAlertApi) UpdateSystemAlert(c *gin.Context) {
	var model struct {
		ID int `json:"id"`
	}
	err := c.ShouldBindJSON(&model)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = global.GVA_DB.Model(&system.SystemAlert{}).Where("id", model.ID).Updates(&system.SystemAlert{
		IsCheck: true,
	}).Error
	if err != nil {
		global.GVA_LOG.Error("修改失败!", zap.Error(err))
		response.FailWithMessage("修改失败", c)
		return
	}
	response.OkWithDetailed(gin.H{"updateAt": time.Now().Unix()}, "修改成功", c)
}

// 删除系统报警
func (b *SystemAlertApi) DeleteSystemAlert(c *gin.Context) {
	var model system.SystemAlert
	err := c.ShouldBindJSON(&model)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = systemAlertService.DeleteSystemAlertService(model.ID)
	if err != nil {
		global.GVA_LOG.Error("删除失败!", zap.Error(err))
		response.FailWithMessage("删除失败", c)
		return
	}
	response.OkWithDetailed(gin.H{"deleteAt": time.Now().Unix()}, "删除成功", c)
}
