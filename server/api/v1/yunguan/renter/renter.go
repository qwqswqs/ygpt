package renter

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"time"
	"ygpt/server/global"
	"ygpt/server/model/common/request"
	"ygpt/server/model/common/response"
	"ygpt/server/model/system"
	systemReq "ygpt/server/model/system/request"
	"ygpt/server/model/yunguan/renter"
	"ygpt/server/service"
	"ygpt/server/utils"
)

type RenterApi struct {
}

// 添加
func (p *RenterApi) AddRenter(c *gin.Context) {
	var r systemReq.AddTenantReq
	r.NickName = r.Username
	err := c.ShouldBindJSON(&r)
	if err != nil {
		response.FailWithMessage(err, c)
		return
	}
	var authorities []system.SysAuthority
	for _, v := range r.AuthorityIds {
		authorities = append(authorities, system.SysAuthority{
			AuthorityId: v,
		})
	}
	user := &system.SysUser{
		Username:    r.Username,
		NickName:    r.Username,
		Password:    r.Password,
		CompanyName: r.CompanyName,
		HeaderImg:   r.HeaderImg,
		AuthorityId: r.AuthorityId,
		Authorities: authorities,
		Enable:      r.Enable,
		Phone:       r.Phone,
		Email:       r.Email,
		Source:      r.Source,
	}
	userReturn, err := userService.Register(*user)
	if err != nil {
		global.GVA_LOG.Error("注册失败!", zap.Error(err))
		response.FailWithMessage("注册失败", c)
		return
	}
	//response.OkWithDetailed(gin.H{"id": userReturn.ID, "username": userReturn.Username}, "注册成功", c)
	renter := renter.Renter{
		RenterID:    int(userReturn.ID),
		Username:    userReturn.Username,
		Email:       userReturn.Email,
		CompanyName: userReturn.CompanyName,
		Source:      userReturn.Source,
		Status:      userReturn.Enable,
	}

	serviceReturn, err := renterService.AddRenterService(renter)
	if err != nil {
		global.GVA_LOG.Error("添加失败!", zap.Error(err))
		response.FailWithDetailed(global.GVA_MODEL{ID: serviceReturn.ID}, "添加失败", c)
		return
	}
	response.OkWithMessage("添加成功", c)
}

// 修改
func (p *RenterApi) UpdateRenter(c *gin.Context) {
	var r renter.Renter
	err := c.ShouldBindJSON(&r)
	if err != nil {
		response.FailWithMessage(err, c)
		return
	}
	err = renterService.UpdateRenterService(r)
	if err != nil {
		global.GVA_LOG.Error("更新失败!", zap.Error(err))
		response.FailWithMessage("更新失败", c)
		return
	}
	response.OkWithDetailed(gin.H{"updateAt": time.Now().Unix()}, "修改成功", c)
}

// 查询
func (p *RenterApi) QueryRenter(c *gin.Context) {
	var pageInfo renter.GetRenterListReq
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
	list, total, err := renterService.GetRenterListService(pageInfo)
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

// 查询
func (p *RenterApi) QueryRenterInfo(c *gin.Context) {
	var data renter.Renter
	err := global.GVA_DB.Model(&renter.Renter{}).Where("renter_id = ?", utils.GetUserID(c)).Find(&data).Error
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithDetailed(gin.H{"data": data}, "获取成功", c)
}

// 查询
func (p *RenterApi) QueryResRenter(c *gin.Context) {
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
	list, total, err := renterService.GetRenterResListService(pageInfo)
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
func (p *RenterApi) DeleteRenter(c *gin.Context) {
	var r renter.Renter
	err := c.ShouldBindJSON(&r)
	if err != nil {
		response.FailWithMessage(err, c)
		return
	}
	userService.DeleteUser(r.RenterID)
	err = renterService.DeleteRenterService(r.ID)
	if err != nil {
		global.GVA_LOG.Error("删除失败!", zap.Error(err))
		response.FailWithMessage("删除失败", c)
		return
	}
	response.OkWithDetailed(gin.H{"deleteAt": time.Now().Unix()}, "删除成功", c)
}

var userService = service.ServiceGroupApp.SystemServiceGroup.UserService

func (b *RenterApi) GetUserList(c *gin.Context) {
	var pageInfo systemReq.GetUserList
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
	list, total, err := userService.GetOpsUserInfoList(pageInfo)
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
func (b *RenterApi) GetAllRenterList(c *gin.Context) {
	list, total, err := userService.GetAllRenterInfoList()
	if err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage(err, c)
		return
	}
	response.OkWithDetailed(response.PageResult{
		List:  list,
		Total: total,
	}, "获取成功", c)
}
func (b *RenterApi) GetAllUserList(c *gin.Context) {
	list, total, err := userService.GetAllUserInfoList()
	if err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage(err, c)
		return
	}
	response.OkWithDetailed(response.PageResult{
		List:  list,
		Total: total,
	}, "获取成功", c)
}
