package renter

import (
	"time"
	"ygpt/server/global"
	"ygpt/server/model/common/response"
	"ygpt/server/model/system"
	"ygpt/server/model/yunguan/renter"
	"ygpt/server/service/compute/websocket/protocol"
	"ygpt/server/service/compute/websocket/service"
	"ygpt/server/utils"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type RenterResApi struct {
}

// 添加
func (p *RenterResApi) AddRenterRes(c *gin.Context) {
	var r renter.RenterRes
	err := c.ShouldBindJSON(&r)
	if err != nil {
		response.FailWithMessage(err, c)
		return
	}
	r.StartTime = time.Now()
	serviceReturn, err := renterResService.AddRenterResService(r)
	if err != nil {
		global.GVA_LOG.Error("添加失败!", zap.Error(err))
		response.FailWithDetailed(global.GVA_MODEL{ID: serviceReturn.ID}, "添加失败", c)
		return
	}
	response.OkWithDetailed(global.GVA_MODEL{ID: serviceReturn.ID}, "添加成功", c)
}

// 修改
func (p *RenterResApi) UpdateRenterRes(c *gin.Context) {
	var r renter.RenterRes
	err := c.ShouldBindJSON(&r)
	if err != nil {
		response.FailWithMessage(err, c)
		return
	}
	err = renterResService.UpdateRenterResService(r)
	if err != nil {
		global.GVA_LOG.Error("更新失败!", zap.Error(err))
		response.FailWithMessage("更新失败", c)
		return
	}
	response.OkWithDetailed(gin.H{"updateAt": time.Now().Unix()}, "修改成功", c)
}

// 查询
func (p *RenterResApi) QueryRenterRes(c *gin.Context) {
	var pageInfo renter.GetRenterResReq
	err := c.ShouldBindJSON(&pageInfo)
	if err != nil {
		response.FailWithMessage(err, c)
		return
	}
	if pageInfo.RenterID == 0 {
		pageInfo.RenterID = int(utils.GetUserID(c))
	}

	err = utils.Verify(pageInfo, utils.PageInfoVerify)
	if err != nil {
		response.FailWithMessage(err, c)
		return
	}
	list, total, err := renterResService.GetRenterResListService(pageInfo)
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

func (p *RenterResApi) QueryRenterResInfo(c *gin.Context) {
	var info renter.GetRenterResReqByResourceID
	err := c.ShouldBindJSON(&info)
	if err != nil {
		response.FailWithMessage(err, c)
		return
	}

	serviceReturn, err := renterResService.QueryRenterResInfo(info)
	if err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage("获取失败:"+err.Error(), c)
		return
	}

	response.OkWithDetailed(serviceReturn, "获取成功", c)
}

func (p *RenterResApi) QueryRenterResCountByTicketId(c *gin.Context) {
	var info renter.GetRenterResCountByTicketId
	err := c.ShouldBindJSON(&info)
	if err != nil {
		response.FailWithMessage(err, c)
		return
	}
	if info.TicketId == 0 {
		response.FailWithMessage("ticketId不能为空", c)
		return
	}
	count, err := renterResService.QueryRenterResCountByTicketId(info.TicketId)
	if err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage("获取失败:"+err.Error(), c)
		return
	}
	response.OkWithDetailed(gin.H{"count": count}, "获取成功", c)
}

// 获取所有租户资源列表
func (p *RenterResApi) QueryRenterResList(c *gin.Context) {
	var info struct {
		PageIndex int    `json:"pageIndex"`
		PageSize  int    `json:"pageSize"`
		RenterID  *int64 `json:"renterID"`
		IsRenter  *int64 `json:"isRenter"`
		Type      int64  `json:"type"`
	}
	var dataList []renter.RenterRes
	err := c.ShouldBindJSON(&info)
	if err != nil {
		response.FailWithMessage(err, c)
		return
	}
	db := global.GVA_DB.Model(&renter.RenterRes{})
	if info.RenterID != nil {
		db.Where("renter_id = ?", *info.RenterID)
	}
	if info.IsRenter != nil {
		db.Where("renter_id = ?", utils.GetUserID(c))
	}
	if info.Type != 0 {
		db.Where("type = ?", info.Type)
	}
	var total int64
	err = db.Count(&total).Error
	if err != nil {
		return
	}
	err = db.Order("created_at DESC").Limit(info.PageSize).Offset((info.PageIndex - 1) * info.PageSize).Find(&dataList).Error
	if err != nil {
		response.FailWithMessage("获取失败", c)
	}
	response.OkWithDetailed(gin.H{"list": dataList, "total": total}, "获取成功", c)
}

// 获取所有租户资源列表
func (p *RenterResApi) QueryRenterFinance(c *gin.Context) {
	var req protocol.ReqInfo
	err := c.ShouldBindQuery(&req)
	if err != nil {
		global.GVA_LOG.Error("参数错误", zap.Error(err))
		response.FailWithMessage("参数错误", c)
		return
	}
	//用keyword区别前端的行为
	if req.PageInfo.Keyword == "node" {
		//云管租户端需要查询所有人的财务信息，还有条件查询
		req.PageInfo.Keyword = ""
		req.NodeCode = &global.GVA_CONFIG.System.NodeCode
	} else {
		//云管租户端需要查询个人的财务信息
		userId := int(utils.GetUserID(c))
		var user system.SysUser
		if err := global.GVA_DB.Where("id = ?", userId).First(&user).Error; err != nil {
			global.GVA_LOG.Error("查询失败", zap.Error(err))
			response.FailWithMessage("查询失败", c)
			return
		}
		req.BuyerPhone = &user.Phone
	}
	//发送请求,解析应答
	ret, err := service.SuanLiServiceGroupApp.TenantService.RequestTenantInfoList(req)
	if err != nil {
		global.GVA_LOG.Error("请求失败", zap.Error(err))
		response.FailWithMessage("请求失败，请检查网络连接", c)
		return
	}
	response.OkWithDetailed(response.PageResult{
		List:     ret.Payload.List,
		Total:    int64(ret.Payload.Total),
		Page:     req.PageInfo.Page,
		PageSize: req.PageInfo.PageSize,
	}, "查询成功", c)
}

// 删除
func (p *RenterResApi) DeleteRenterRes(c *gin.Context) {
	var r renter.RenterRes
	err := c.ShouldBindJSON(&r)
	if err != nil {
		response.FailWithMessage(err, c)
		return
	}
	err = renterResService.DeleteRenterResService(r.ID)
	if err != nil {
		global.GVA_LOG.Error("删除失败!", zap.Error(err))
		response.FailWithMessage("删除失败", c)
		return
	}
	response.OkWithDetailed(gin.H{"deleteAt": time.Now().Unix()}, "删除成功", c)
}

// 分配资源信息
func (p *RenterResApi) BindRenterRes(c *gin.Context) {
	var r renter.RenterRes
	err := c.ShouldBindJSON(&r)
	if err != nil {
		response.FailWithMessage(err, c)
		return
	}
	err = renterResService.BindResInfoService(r)
	if err != nil {
		global.GVA_LOG.Error("分配失败!", zap.Error(err))
		response.FailWithMessage("分配失败", c)
		return
	}
	response.OkWithDetailed(gin.H{"updateAt": time.Now().Unix()}, "分配成功", c)

	// 向算力圈发送请求   传入租户资源表id
	//wsService.SuanLiServiceGroupApp.ResourceService.AllocateResponse(int(r.ID))
}

// 释放资源信息
func (p *RenterResApi) ReleaseRenterRes(c *gin.Context) {
	var r renter.RenterRes
	err := c.ShouldBindJSON(&r)
	if err != nil {
		response.FailWithMessage(err, c)
		return
	}
	err = renterResService.ReleaseResInfoService(r)
	if err != nil {
		global.GVA_LOG.Error("释放失败!", zap.Error(err))
		response.FailWithMessage("释放失败", c)
		return
	}
	response.OkWithDetailed(gin.H{"updateAt": time.Now().Unix()}, "释放成功", c)

	// 向算力圈发送请求
}

// 获取与订单相关的租户资源列表
func (p *RenterResApi) GetRenterResByTicketId(c *gin.Context) {
	var r renter.RenterRes
	err := c.ShouldBindJSON(&r)
	if err != nil {
		response.FailWithMessage(err, c)
		return
	}
	data, err := renterResService.GetResInfoByTicketService(r)
	if err != nil {
		global.GVA_LOG.Error("释放失败!", zap.Error(err))
		response.FailWithMessage("释放失败", c)
		return
	}
	response.OkWithDetailed(data, "获取成功", c)
}
