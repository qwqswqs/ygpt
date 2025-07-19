package backup

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"time"
	"ygpt/server/global"
	"ygpt/server/model/common/response"
	"ygpt/server/model/yunguan/backup"
	"ygpt/server/utils"
)

type BackApi struct {
}

// 添加
func (p *BackApi) AddBack(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		panic(fmt.Sprintf("file参数不能为空"))
	}
	var bak backup.Backup
	bak.Name = file.Filename
	bak.Size = int(file.Size)
	bak.UploadTime = time.Now()
	bak.RenterID = int(utils.GetUserID(c))
	serviceReturn, err := BackService.AddBackService(bak)
	if err != nil {
		global.GVA_LOG.Error("添加失败!", zap.Error(err))
		response.FailWithDetailed(global.GVA_MODEL{ID: serviceReturn.ID}, "添加失败", c)
		return
	}
	response.OkWithMessage("添加成功", c)
}

// 修改
func (p *BackApi) UpdateBack(c *gin.Context) {
	var r backup.Backup
	err := c.ShouldBindJSON(&r)
	if err != nil {
		response.FailWithMessage(err, c)
		return
	}
	err = BackService.UpdateBackService(r)
	if err != nil {
		global.GVA_LOG.Error("更新失败!", zap.Error(err))
		response.FailWithMessage("更新失败", c)
		return
	}
	response.OkWithDetailed(gin.H{"updateAt": time.Now().Unix()}, "修改成功", c)
}

// 查询
func (p *BackApi) QueryBack(c *gin.Context) {
	var pageInfo backup.BackReq
	err := c.ShouldBindJSON(&pageInfo)
	if err != nil {
		response.FailWithMessage(err, c)
		return
	}
	err = utils.Verify(pageInfo, utils.RegisterVerify)
	if err != nil {
		response.FailWithMessage(err, c)
		return
	}
	list, total, err := BackService.GetBackListService(pageInfo)
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

// 查询租户网盘信息
func (p *BackApi) QueryUserBack(c *gin.Context) {
	var pageInfo backup.BackReq
	err := c.ShouldBindJSON(&pageInfo)
	if err != nil {
		response.FailWithMessage(err, c)
		return
	}
	userId := utils.GetUserID(c)
	pageInfo.RenterID = int(userId)
	list, total, err := BackService.GetBackListService(pageInfo)
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
func (p *BackApi) DeleteBack(c *gin.Context) {
	var r backup.Backup
	err := c.ShouldBindJSON(&r)
	if err != nil {
		response.FailWithMessage(err, c)
		return
	}
	err = BackService.DeleteBackService(r.ID)
	if err != nil {
		global.GVA_LOG.Error("删除失败!", zap.Error(err))
		response.FailWithMessage("删除失败", c)
		return
	}
	response.OkWithDetailed(gin.H{"deleteAt": time.Now().Unix()}, "删除成功", c)
}
