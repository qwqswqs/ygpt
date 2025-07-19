package system

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"strconv"
	"time"
	"ygpt/server/global"
	"ygpt/server/model/common/request"
	"ygpt/server/model/common/response"
	"ygpt/server/model/yunguan/system"
	"ygpt/server/utils"
)

type SystemCacheApi struct {
}

// 添加系统运行工具
func (d *SystemCacheApi) AddSystemCache(c *gin.Context) {
	var r system.SystemCache
	err := c.ShouldBindJSON(&r)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = utils.Verify(r, utils.RegisterVerify)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	data := system.SystemCache{
		Address:    r.Address,
		Size:       r.Size,
		Type:       r.Type,
		UploadTime: time.Now(),
		UserType:   r.UserType,
		UserID:     r.UserID,
	}
	serviceReturn, err := systemCacheService.AddSystemCacheService(data)
	if err != nil {
		global.GVA_LOG.Error("添加失败!", zap.Error(err))
		response.FailWithDetailed(global.GVA_MODEL{ID: serviceReturn.ID}, "添加失败", c)
		return
	}
	response.OkWithMessage("添加成功", c)
}

// 修改系统运行工具
func (d *SystemCacheApi) UpdateSystemCache(c *gin.Context) {
	var model system.SystemCache
	err := c.ShouldBindJSON(&model)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = systemCacheService.UpdateSystemCacheService(model)
	if err != nil {
		global.GVA_LOG.Error("更新失败!", zap.Error(err))
		response.FailWithMessage("更新失败", c)
		return
	}
	response.OkWithDetailed(gin.H{"updateAt": time.Now().Unix()}, "修改成功", c)
}

// 获取系统运行工具列表
func (b *SystemCacheApi) GetSystemCacheList(c *gin.Context) {
	PageIndex, err := strconv.Atoi(c.Query("pageIndex"))
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	PageSize, err := strconv.Atoi(c.Query("pageSize"))
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	var pageInfo request.PageInfo
	pageInfo.Page = PageIndex
	pageInfo.PageSize = PageSize
	list, total, err := systemCacheService.GetSystemCacheListService(pageInfo)
	if err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage("获取失败", c)
		return
	}
	response.OkWithDetailed(response.PageResult{
		List:     list,
		Total:    total,
		Page:     pageInfo.Page,
		PageSize: pageInfo.PageSize,
	}, "获取成功", c)
}

// 删除系统运行工具
func (b *SystemCacheApi) DeleteSystemCache(c *gin.Context) {
	var model system.SystemCache
	err := c.ShouldBindJSON(&model)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = systemCacheService.DeleteSystemCacheService(model.ID)
	if err != nil {
		global.GVA_LOG.Error("删除失败!", zap.Error(err))
		response.FailWithMessage("删除失败", c)
		return
	}
	response.OkWithDetailed(gin.H{"deleteAt": time.Now().Unix()}, "删除成功", c)
}
