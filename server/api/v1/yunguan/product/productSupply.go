package product

import (
	"errors"
	"fmt"
	"math"
	"strconv"
	"time"
	"ygpt/server/global"
	"ygpt/server/model/common/request"
	"ygpt/server/model/common/response"
	"ygpt/server/model/yunguan/product"
	"ygpt/server/model/yunguan/system"
	wsService "ygpt/server/service/compute/websocket/service"
	"ygpt/server/utils"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"go.uber.org/zap"
)

type ProductSupplyApi struct {
}

// 添加
func (p *ProductSupplyApi) AddSupplyInfo(c *gin.Context) {
	var r product.ProductSupply
	err := c.ShouldBindJSON(&r)
	if err != nil {
		response.FailWithMessage(err, c)
		return
	}
	serviceReturn, err := productSupplyService.AddProductSupplyService(r)
	if err != nil {
		global.GVA_LOG.Error("添加失败!", zap.Error(err))
		response.FailWithDetailed(global.GVA_MODEL{ID: uint(serviceReturn.Id)}, "添加失败", c)
		return
	}
	response.OkWithDetailed(global.GVA_MODEL{ID: uint(serviceReturn.Id)}, "添加成功", c)
}

// 修改
func (p *ProductSupplyApi) UpdateSupplyInfo(c *gin.Context) {
	var r product.ProductSupply
	err := c.ShouldBindJSON(&r)
	if err != nil {
		response.FailWithMessage(err, c)
		return
	}
	err = productSupplyService.UpdateProductSupplyService(r)
	if err != nil {
		global.GVA_LOG.Error("更新失败!", zap.Error(err))
		response.FailWithMessage("更新失败", c)
		return
	}
	response.OkWithDetailed(gin.H{"updateAt": time.Now().Unix()}, "修改成功", c)
}

// 查询
func (p *ProductSupplyApi) QuerySupplyInfo(c *gin.Context) {
	var pageInfo product.GetProductSupplyReq
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
	list, total, err := productSupplyService.GetSupplyInfoListService(pageInfo)
	if err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage(err, c)
		return
	}
	response.OkWithDetailed(response.PageResult{
		List:     list,
		Total:    total,
		Page:     pageInfo.PageIndex,
		PageSize: pageInfo.PageSize,
	}, "获取成功", c)
}

// 删除
func (p *ProductSupplyApi) DeleteSupplyInfo(c *gin.Context) {
	// 获取前端传递的 ID
	id := c.Query("id") // 获取 URL 查询参数 "?id=123"
	//id为int
	Id, err := strconv.Atoi(id)
	if err != nil {
		response.FailWithMessage("ID 不能为空", c)
		return
	}
	// 供需信息在数据库中共用一张表，所以不去要特别区别两者的类型，因为ID不会重复，所以产品类型(3.供给 4.需求)在此处统一为3
	err = wsService.SuanLiServiceGroupApp.SyncProductStatusChangeToPlatform(Id, 3, 5)
	if err != nil {
		global.GVA_LOG.Error("删除失败!无法同步到算力平台", zap.Error(err))
		response.FailWithMessage("删除失败!无法同步到算力平台", c)
		return
	}
	err = productSupplyService.DeleteProductSupplyService(int64(Id))
	if err != nil {
		global.GVA_LOG.Error("删除失败!", zap.Error(err))
		response.FailWithMessage("删除失败", c)
		return
	}
	response.OkWithDetailed(gin.H{"deleteAt": time.Now().Unix()}, "删除成功", c)
}

// 同步单个数据
func (p *ProductSupplyApi) SyncSupplyDemandInfo(c *gin.Context) {
	var r product.ProductSupply
	err := c.ShouldBindJSON(&r)
	if err != nil {
		response.FailWithMessage(err, c)
		return
	}
	err = wsService.SuanLiServiceGroupApp.SupplyInfoRelease(int(r.Id))
	if err != nil {
		global.GVA_LOG.Error("同步失败!", zap.Error(err))
		response.FailWithMessage("同步失败", c)
		return
	}
	r.Status = -2
	err = productSupplyService.UpdateProductSupplyService(r)
	if err != nil {
		global.GVA_LOG.Error("同步失败!", zap.Error(err))
		response.FailWithMessage("同步失败", c)
		return
	}
	response.OkWithDetailed(gin.H{"updateAt": time.Now().Unix()}, "修改成功", c)
}

type ListSupplyReq struct {
	request.PageInfo
	//RegionId 城市id
	RegionId *int `json:"regionId" form:"regionId"`
	//Id 类别id
	Id *int `json:"id"  form:"id"`
	//SupplyDemand  1=云上算力产品 2=元素产品 3=供给信息 4=需求信息
	Type *int `json:"type" form:"type"`
	//SubType 1-1=裸金属，1-2=云主机，1-3=容器，1-4=存储，2-5=模型，2-6=数据集，3-7=算力，3-8=硬件，3-9=机柜，3-10=带宽，11其他
	SubType *int `json:"subType"  form:"subType"`
	//Status 审核状态（-1=驳回，0=未审核*，1=通过）
	Status *int `json:"status" form:"status"` //
	//ReNodeCode 节点编号
	ReNodeCode *string `json:"reNodeCode"   form:"reNodeCode"`
	//Source 来源 0=平台 1=用户 2=节点
	Source     *int   `json:"source"  form:"source"`
	TypeDesc   string `json:"typeDesc"`
	StatusDesc string `json:"statusDesc"`
}

type ListSupplyResp struct {
	RegionRoot system.ProductRegion  `json:"rootRegion"`
	Root       system.ProductSysType `json:"root"`
	Region     system.ProductRegion  `json:"region"`
	Product    product.ProductSupply `json:"product"`
	SysType    system.ProductSysType `json:"sysType"`
}

// GetSupplyList 获取product
func (p *ProductSupplyApi) GetSupplyList(c *gin.Context) {
	var req ListSupplyReq
	err := c.ShouldBindQuery(&req)
	if err != nil {
		response.FailWithMessage(err, c)
		return
	}
	fmt.Println(req)
	db := global.GVA_DB.Model(&product.ProductSupply{}).Where("is_deleted  = 0 ")
	condition := product.ProductSupply{}
	if req.RegionId != nil {
		db = db.Where("region_id = ?", req.RegionId)
		condition.RegionId = int64(*req.RegionId)
	}
	if req.Id != nil {
		db = db.Where("type_id = ?", req.Type)
		condition.TypeId = int64(*req.Type)
	}
	if req.Keyword != "" {
		db = db.Where("name LIKE ?", "%"+req.Keyword+"%")
	}
	if req.Type != nil {
		db = db.Where("type = ?", req.Type)
		condition.Type = *req.Type
	}
	if req.SubType != nil {
		db = db.Where("subtype = ?", req.SubType)
		condition.SubType = *req.SubType
	}
	if req.Status != nil {
		db = db.Where("status = ?", req.Status)
		condition.Status = int8(*req.Status)
	}
	if req.ReNodeCode != nil {
		db = db.Where("re_node_code = ?", req.ReNodeCode)
		condition.ReNodeCode = *req.ReNodeCode
	}
	if req.Source != nil {
		db = db.Where("source = ?", req.Source)
		condition.Source = *req.Source
	}
	if req.TypeDesc == "desc" {
		db = db.Order("type desc")
	} else {
		db = db.Order("type asc")
	}
	if req.StatusDesc == "desc" {
		db = db.Order("status desc")
	} else {
		db = db.Order("status asc")
	}
	var total int64
	db.Count(&total)
	var products []product.ProductSupply
	db.Offset((req.Page - 1) * req.PageSize).Limit(req.PageSize).Find(&products)
	var res []ListSupplyResp
	for _, product := range products {
		//查询区域
		var region system.ProductRegion
		global.GVA_DB.Model(&system.ProductRegion{}).Where("id = ?", product.RegionId).Find(&region)
		//查询SysType
		var SysType system.ProductSysType
		global.GVA_DB.Model(&system.ProductSysType{}).Where("id  = ?", product.TypeId).First(&SysType)
		// 查询rootRegion
		var rootRegion system.ProductRegion
		global.GVA_DB.Model(&system.ProductRegion{}).Where("id  = ?", region.ParentId).First(&rootRegion)
		// 查询rootType
		var rootType system.ProductSysType
		global.GVA_DB.Model(&system.ProductSysType{}).Where("id = ? ", SysType.ParentId).First(&rootType)
		res = append(res, ListSupplyResp{
			RegionRoot: rootRegion,
			Root:       rootType,
			Region:     region,
			Product:    product,
			SysType:    SysType,
		})
	}
	//响应
	response.OkWithDetailed(response.PageResult{
		List:     res,
		Total:    total,
		Page:     req.Page,
		PageSize: req.PageSize,
	}, "获取成功", c)
}

type CreateSupplyReq struct {
	//Name 标题
	Name string `json:"name"  form:"name"`
	//Source 0=平台 1=用户 2=节点
	Source int `json:"source" form:"source"`
	//Description 描述
	Description string `json:"description"  form:"description"`
	//Price 价格
	Price float64 `json:"price" form:"price"`
	//Inventory  库存
	Inventory int32 `json:"inventory" form:"inventory"`
	//HeadImage 图片地址
	HeadImage string `json:"headImage" form:"headImage"`
	//Type 1=云上算力产品 2=元素产品 3=供给信息 4=需求信息
	Type int `json:"type"   form:"type"`
	//SubType 1-1=裸金属，1-2=云主机，1-3=容器，1-4=存储，2-5=模型，2-6=数据集，3-7=算力，3-8=硬件，3-9=机柜，3-10=带宽，11其他
	SubType int `json:"subType"  form:"subType"`
	//StartTime 开始时间
	LeaseStartTime time.Time `json:"startTime"  form:"startTime"`
	//EndTime 结束时间
	LeaseEndTime time.Time `json:"endTime"  form:"endTime"`
	//IsDownload 0=可下载 1不可下载"
	IsDownload int8 `json:"isDownload"  form:"isDownload" `
	// 状态
	Status int `json:"status" form:"status"`

	//SpecJson 规格的JSON参数
	SpecJson string `json:"specJson" form:"specJson" `
	//RegionId 城市id
	RegionId int `json:"regionId"  form:"regionId"`

	//User 联系人
	User string `json:"user"  form:"user"`
	//UserPhone 联系电话
	UserPhone string `json:"userPhone"  form:"userPhone"`
	//UserEmail 联系电话
	UserEmail string `json:"userEmail"  form:"userEmail"`
	//ValidTime 有效时间
	ValidTime time.Time `json:"validTime" form:"validTime"`
	//TypeId 关联id
	TypeId int64 `json:"typeId" form:"typeId" `
	//uploadType
	UploadType *int `json:"uploadType"  form:"uploadType"`

	//Url  文件名称
	Url *string `json:"url"  form:"url"`
	// 模型参数
	Parameter *float64 `json:"parameter" form:"parameter"`
}

func (p *ProductSupplyApi) CreateSupplyInfo(c *gin.Context) {
	var req CreateSupplyReq
	err := c.ShouldBindJSON(&req)
	if err != nil {
		global.GVA_LOG.Error("参数错误!", zap.Error(err))
		response.FailWithMessage("参数错误!", c)
		return
	}
	req.Source = 6
	//正常供需
	var productSupply product.ProductSupply
	err = copier.Copy(&productSupply, &req)
	productSupply.CreateTime = time.Now()
	productSupply.Status = 0
	productSupply.CreateBy = utils.GetUserName(c)
	if req.LeaseStartTime.IsZero() || req.LeaseEndTime.IsZero() {
		// 默认信息开始时间是此刻
		productSupply.LeaseStartTime = time.Now()
		// 默认信息有效期是一年
		productSupply.LeaseEndTime = time.Now().AddDate(0, 1, 0)
	}
	if err = global.GVA_DB.Debug().Model(&product.ProductSupply{}).Create(&productSupply).Error; err != nil {
		global.GVA_LOG.Error("发布供需错误!", zap.Error(err))
		response.FailWithMessage("发布供需错误!", c)
		return
	}
	response.OkWithMessage("发布供需成功", c)
	return
}

func (p *ProductSupplyApi) DeleteSupply(c *gin.Context) {
	id := c.Query("id")
	productId, err := strconv.Atoi(id)
	if err != nil {
		global.GVA_LOG.Error("删除失败!", zap.Error(err))
		response.FailWithMessage("删除失败:"+err.Error(), c)
		return
	}
	deleteMap := map[string]interface{}{
		"is_deleted":  1,
		"update_time": time.Now(),
	}
	if err = global.GVA_DB.Model(&product.ProductSupply{}).Where("id = ?", productId).Updates(deleteMap).Error; err != nil {
		global.GVA_LOG.Error("删除失败!", zap.Error(err))
		response.FailWithMessage("删除失败:"+err.Error(), c)
		return
	}

	response.OkWithMessage("删除成功", c)
}

func StrslicetoInt(strings []string) []int {
	intSlice := make([]int, 0)
	for _, v := range strings {
		atoi, err := strconv.Atoi(v)
		if err != nil {
			atoi = 0
		}
		intSlice = append(intSlice, atoi)
	}
	return intSlice
}

func (p *ProductSupplyApi) DeleteSupplyByIds(c *gin.Context) {
	ids := c.QueryArray("ids[]")
	idsInt := StrslicetoInt(ids)
	var productV2 product.ProductSupply
	productV2.UpdateTime = time.Now()
	productV2.IsDeleted = 1
	if err := global.GVA_DB.Model(&product.ProductSupply{}).Where("id  in ? ", idsInt).Updates(&productV2).Error; err != nil {
		global.GVA_LOG.Error("批量删除失败!", zap.Error(err))
		response.FailWithMessage("批量删除失败!", c)
		return
	}
	response.OkWithMessage("批量删除成功", c)
}

type SysTypeSearch struct {
	request.PageInfo
	//Type 1表示元素 2表示算力
	Type *int `json:"type"  form:"type"`
}

type SupplyConfigResp struct {
	Root   system.ProductSysType   `json:"root" form:"root"`
	Second []system.ProductSysType `json:"second" form:"second"`
	Third  []system.ProductSysType `json:"third" form:"third"`
}

func GetSysTypeInfo(info SysTypeSearch) (list []SupplyConfigResp, total int64, err error) {
	//limit := info.PageSize
	//offset := info.PageSize * (info.Page - 1)
	//元素
	ele := []int{5, 6, 7, 8}
	//算力
	algo := []int{1, 2, 3, 4}
	var list2 []system.ProductSysType
	// 创建db
	db := global.GVA_DB.Model(&system.ProductSysType{})
	if info.Type != nil && *info.Type > 0 {
		if *info.Type == 1 {
			db.Where("is_deleted = 0 and parent_id  =0 and category  IN  ?", ele)
			db.Count(&total)
			if err = db.Debug().Find(&list2).Error; err != nil {
				return list, total, err
			}
			fmt.Println(list2)
			for _, v := range list2 {
				var eRet SupplyConfigResp
				eRet.Root = v
				var types []system.ProductSysType
				global.GVA_DB.Debug().Model(&system.ProductSysType{}).
					Where("is_deleted = 0 and parent_id = ?", v.Id).Find(&types)
				eRet.Second = types
				list = append(list, eRet)
			}
			fmt.Println(len(list))
			return list, total, nil
		}

		if *info.Type == 2 {
			root := make([]system.ProductSysType, 0)
			//res := make([]response.SysTypeListResponse, 0)

			db.Where("is_deleted = 0 and parent_id  =0 and category  IN  ?", algo)
			db.Count(&total)
			if err = db.Debug().Find(&root).Error; err != nil {
				return list, total, err
			}
			for _, v := range root {
				var e SupplyConfigResp
				e.Root = v
				list = append(list, e)
			}
			for i := range list {
				second := make([]system.ProductSysType, 0)
				global.GVA_DB.Debug().Model(&system.ProductSysType{}).Where("is_deleted = 0 and parent_id = ?", list[i].Root.Id).Find(&second)
				list[i].Second = second
				for _, v := range second {
					var third []system.ProductSysType
					global.GVA_DB.Debug().Model(&system.ProductSysType{}).Where("is_deleted = 0 and parent_id =?", v.Id).Find(&third)
					list[i].Third = append(list[i].Third, third...)
				}
			}

		}
		return list, total, nil

	}
	return list, total, nil
}

func (p *ProductSupplyApi) GetSupplyConfig(c *gin.Context) {
	var pageInfo SysTypeSearch
	err := c.ShouldBindQuery(&pageInfo)
	if err != nil {
		response.FailWithMessage(err, c)
		return
	}
	list, total, err := GetSysTypeInfo(pageInfo)
	if err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage("获取失败:"+err.Error(), c)
		return
	}
	response.OkWithDetailed(response.PageResult{
		List:     list,
		Total:    total,
		Page:     pageInfo.Page,
		PageSize: pageInfo.PageSize,
	}, "获取成功", c)
}

func (p *ProductSupplyApi) GetSupplyCityList(c *gin.Context) {
	var req request.PageInfo
	err := c.ShouldBindQuery(&req)
	if err != nil {
		global.GVA_LOG.Error("参数错误!", zap.Error(err))
		response.FailWithMessage("参数错误:"+err.Error(), c)
		return
	}
	type Ret struct {
		Province system.ProductRegion   `json:"province"`
		Cities   []system.ProductRegion `json:"cities"`
	}
	var ret []Ret
	var count int64
	var provinces []system.ProductRegion
	//limit := req.PageSize
	//offset := (req.Page - 1) * req.PageSize

	db := global.GVA_DB.Model(&system.ProductRegion{}).Where("is_deleted = 0  and parent_id  = 1")
	//筛选父id为0
	if req.Keyword != "" {
		db.Where("name  LIKE  ?", req.Keyword)
	}
	err = db.Count(&count).Error
	//if limit > 0 {
	//	db.Limit(limit).Offset(offset).Order("id  asc")
	//}

	err = db.Debug().Find(&provinces).Error
	for _, v := range provinces {
		var eRet Ret
		eRet.Province = v
		var cities []system.ProductRegion
		global.GVA_DB.Debug().Model(&system.ProductRegion{}).Where("is_deleted  = 0  and  parent_id  = ?", v.Id).Find(&cities)
		eRet.Cities = cities
		ret = append(ret, eRet)
	}
	response.OkWithDetailed(response.PageResult{
		List:     ret,
		Total:    count,
		Page:     req.Page,
		PageSize: req.PageSize,
	}, "查询成功", c)
}

func (p *ProductSupplyApi) UpdateSupply(c *gin.Context) {
	var req product.ProductSupply
	err := c.ShouldBindJSON(&req)
	if err != nil {
		global.GVA_LOG.Error("参数错误", zap.Error(err))
		response.FailWithMessage(err, c)
		return
	}
	fmt.Printf("%+v\n", req)
	req.UpdateBy = utils.GetUserName(c)
	req.UpdateTime = time.Now()
	if req.LeaseStartTime.IsZero() || req.LeaseEndTime.IsZero() {
		// 默认信息开始时间是此刻
		req.LeaseStartTime = time.Now()
		// 默认信息有效期是一年
		req.LeaseEndTime = time.Now().AddDate(0, 1, 0)
	}
	err = global.GVA_DB.Debug().Model(&product.ProductSupply{}).Where("id = ?", req.Id).Updates(&req).Error
	if err != nil {
		global.GVA_LOG.Error("更新失败", zap.Error(err))
		response.FailWithMessage("更新失败", c)
		return
	}
	response.OkWithMessage("更新成功", c)
}

// 更新系统类别
func (*ProductSupplyApi) AlterSysType(c *gin.Context) {
	var req struct {
		Id       int    `json:"id"`
		Name     string `json:"name"`
		Category *int8  `json:"category"`
	}
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.FailWithMessage(err, c)
		return
	}
	//更新
	err = global.GVA_DB.Model(&system.ProductSysType{}).Where("id = ?", req.Id).Updates(&system.ProductSysType{
		Name:     req.Name,
		Category: req.Category,
	}).Error
	if err != nil {
		response.FailWithMessage(err, c)
		return
	}
	response.OkWithDetailed(&gin.H{
		"updateTime": time.Now().Unix(),
	}, "更新成功", c)
}

var LabelMap = map[int]string{
	10: "计算类型",
	11: "品牌",
	12: "供电方式",
	13: "运营商",
	14: "卡类型",
	15: "型号",
	16: "路由类型",
	17: "供需规格",
}

var descMap = map[bool]string{
	false: " asc ",
	true:  "desc ",
}

func SlicePage(page, pageSize, nums int) (sliceStart, sliceEnd int) {
	if page < 0 {
		page = 1
	}

	if pageSize < 0 {
		pageSize = 20
	}

	if pageSize > nums {
		return 0, nums
	}

	// 总页数
	pageCount := int(math.Ceil(float64(nums) / float64(pageSize)))
	if page > pageCount {
		return 0, 0
	}
	sliceStart = (page - 1) * pageSize
	sliceEnd = sliceStart + pageSize

	if sliceEnd > nums {
		sliceEnd = nums
	}
	return sliceStart, sliceEnd
}

// distinct 去除同名的数据
func distinct(list []struct {
	Id   int    `json:"id"`
	Type string `json:"type"`
	Name string `json:"name"`
	Unit string `json:"unit"`
}) []struct {
	Id   int    `json:"id"`
	Type string `json:"type"`
	Name string `json:"name"`
	Unit string `json:"unit"`
} {
	var result []struct {
		Id   int    `json:"id"`
		Type string `json:"type"`
		Name string `json:"name"`
		Unit string `json:"unit"`
	}
	var distinctMap map[string]bool = make(map[string]bool)
	for _, v := range list {
		if _, ok := distinctMap[v.Name+v.Type]; !ok {
			distinctMap[v.Name+v.Type] = true
			result = append(result, v)
		}
	}
	return result
}

// 获取类别名称列表
func (*ProductSupplyApi) GetTypeNameList(c *gin.Context) {
	var req struct {
		Type     *int   `json:"type" form:"type"`         // 类别
		Page     int    `json:"pageNum" form:"pageNum"`   // 页码
		PageSize int    `json:"pageSize" form:"pageSize"` // 每页大小
		Keyword  string `json:"keyword" form:"keyword"`   // 关键字
		OrderBy  string `json:"orderBy" form:"orderBy"`   //排序依据
		IsDesc   *bool  `json:"isDesc" form:"isDesc"`     //升降序
		TypeDesc string `json:"typeDesc"`
	}
	err := c.ShouldBindQuery(&req)
	if err != nil {
		response.FailWithMessage(err, c)
		return
	}
	db := global.GVA_DB.Model(&system.ProductSysType{}).Where("is_deleted = 0")
	if req.Type == nil {
		// 没有选择Type查询出来全部
		var TypeList []struct {
			Id   int    `json:"id"`
			Type string `json:"type"`
			Name string `json:"name"`
			Unit string `json:"unit"`
		}
		var parents []system.ProductSysType
		global.GVA_DB.Model(&system.ProductSysType{}).Where("category in (1,2,3,4)").Order("category " + req.TypeDesc).Find(&parents)
		var MapTypeName = map[int8]string{
			3: "计算类型",
			4: "品牌",
			1: "供电方式",
			2: "运营商",
		}
		for _, item := range parents {
			if item.Category != nil {
				var list []system.ProductSysType
				db := global.GVA_DB.Model(&system.ProductSysType{}).Where("parent_id = ?", item.Id)
				if req.Keyword != "" {
					db = db.Where("name LIKE ?", "%"+req.Keyword+"%")
				}
				db.Find(&list)
				for _, listItem := range list {
					TypeList = append(TypeList, struct {
						Id   int    "json:\"id\""
						Type string "json:\"type\""
						Name string "json:\"name\""
						Unit string "json:\"unit\""
					}{
						Id:   *listItem.Id,
						Name: listItem.Name,
						Type: MapTypeName[*item.Category],
						Unit: listItem.Description,
					})
				}
			}
		}
		// 查询下下面的
		for _, item := range parents {
			if item.Category != nil && *item.Category != 1 {
				var ids []int
				var listType []system.ProductSysType
				global.GVA_DB.Model(&system.ProductSysType{}).Select("id").Where("parent_id = ?", item.Id).Find(&ids)
				db := global.GVA_DB.Model(&system.ProductSysType{}).Where("parent_id in ?", ids)
				if req.Keyword != "" {
					db = db.Where("name LIKE ?", "%"+req.Keyword+"%")
				}
				db.Find(&listType)
				var TypeNameMap = map[int8]string{
					3: "卡类型",
					4: "型号",
					2: "路由类型",
				}
				for _, listItem := range listType {
					TypeList = append(TypeList, struct {
						Id   int    "json:\"id\""
						Type string "json:\"type\""
						Name string "json:\"name\""
						Unit string "json:\"unit\""
					}{
						Id:   *listItem.Id,
						Type: TypeNameMap[*item.Category],
						Name: listItem.Name,
						Unit: listItem.Description,
					})
				}
			}
		}
		// 查询规格信息
		TypeList = distinct(TypeList)
		var list []system.ProductSysType
		db := global.GVA_DB.Model(&system.ProductSysType{}).Where("category = 17")
		if req.Keyword != "" {
			db = db.Where("name LIKE ?", "%"+req.Keyword+"%")
		}
		db.Find(&list)
		for _, listItem := range list {
			TypeList = append(TypeList, struct {
				Id   int    "json:\"id\""
				Type string "json:\"type\""
				Name string "json:\"name\""
				Unit string "json:\"unit\""
			}{
				Id:   *listItem.Id,
				Type: "供需规格",
				Name: listItem.Name,
				Unit: listItem.Description,
			})
		}
		begin, end := SlicePage(req.Page, req.PageSize, len(TypeList))
		//响应分页
		response.OkWithData(&gin.H{
			"list":  TypeList[begin:end],
			"total": len(TypeList),
		}, c)
		return
	} else {
		switch *req.Type {
		case 10, 11, 12, 13:
			var parent system.ProductSysType
			var CategoryMap = map[int]int{
				10: 3,
				11: 4,
				12: 1,
				13: 2,
			}
			global.GVA_DB.Model(&system.ProductSysType{}).Where("category = ? ", CategoryMap[*req.Type]).First(&parent)
			db = db.Where("parent_id = ?", parent.Id)
			if req.Keyword != "" {
				db = db.Where("name LIKE ?", "%"+req.Keyword+"%")
			}
		// 如果下面的类型还需要往下匹配
		case 14, 15, 16:
			var parent system.ProductSysType
			var CategoryMap = map[int]int{
				14: 3,
				15: 4,
				16: 2,
			}
			global.GVA_DB.Model(&system.ProductSysType{}).Where("category = ? ", CategoryMap[*req.Type]).First(&parent)
			var ids []int
			global.GVA_DB.Model(&system.ProductSysType{}).Where("parent_id = ?", parent.Id).Select("id").Find(&ids)
			db = db.Where("parent_id in ?", ids)
			if req.Keyword != "" {
				db = db.Where("name LIKE ?", "%"+req.Keyword+"%")
			}
		case 17:
			//这个是对规格选项的获取，搜索category为17的信息并且这里要获取规格
			db = db.Where("category = ? ", 17)
			if req.Keyword != "" {
				db = db.Where("name LIKE ?", "%"+req.Keyword+"%")
			}
		default:
			response.FailWithMessage("错误的类别", c)
			return
		}
	}
	if req.OrderBy != "" && req.IsDesc != nil {
		db = db.Order(fmt.Sprintf("%s %s", "name", descMap[*req.IsDesc]))
	}
	if req.Keyword != "" {
		db = db.Where("name LIKE ?", "%"+req.Keyword+"%")
	}
	var list []system.ProductSysType
	db.Find(&list)
	var resList []struct {
		Id   int    `json:"id"`
		Type string `json:"type"`
		Name string `json:"name"`
		Unit string `json:"unit"`
	}
	for _, item := range list {
		resList = append(resList, struct {
			Id   int    "json:\"id\""
			Type string "json:\"type\""
			Name string "json:\"name\""
			Unit string "json:\"unit\""
		}{
			Id:   *item.Id,
			Type: LabelMap[*req.Type],
			Name: item.Name,
			Unit: item.Description, //这里只能使用description作为Unit了
		})
	}
	//resList = distinct(resList)
	if req.Page == 0 && req.PageSize == 0 {
		response.OkWithData(&gin.H{
			"list":  resList,
			"total": len(resList),
		}, c)
		return
	}
	begin, end := SlicePage(req.Page, req.PageSize, len(resList))
	response.OkWithData(&gin.H{
		"list":  resList[begin:end],
		"total": len(resList),
	}, c)
}

func (*ProductSupplyApi) GetElemSysTypeList(c *gin.Context) {
	var pageInfo struct {
		Type *int `json:"type"  form:"type"`
	}
	err := c.ShouldBindQuery(&pageInfo)
	if err != nil {
		response.FailWithMessage(err, c)
		return
	}
	//limit := info.PageSize
	//offset := info.PageSize * (info.Page - 1)
	//元素
	ele := []int{5, 6, 7, 8}
	//算力
	//algo := []int{1, 2, 3, 4}
	var list2 []system.ProductSysType
	var list []system.SysTypeListResponse
	// 创建db
	db := global.GVA_DB.Model(&system.ProductSysType{})
	if pageInfo.Type != nil && *pageInfo.Type > 0 {
		if *pageInfo.Type == 1 {
			db.Where("is_deleted = 0 and parent_id  =0 and category  IN  ?", ele)
			if err = db.Debug().Find(&list2).Error; err != nil {
				global.GVA_LOG.Error("获取失败!", zap.Error(err))
				response.FailWithMessage("获取失败:"+err.Error(), c)
				return
			}
			fmt.Println(list2)
			for _, v := range list2 {
				var eRet system.SysTypeListResponse
				eRet.Root = v
				var types []system.ProductSysType
				global.GVA_DB.Debug().Model(&system.ProductSysType{}).
					Where("is_deleted = 0 and parent_id = ?", v.Id).Find(&types)
				eRet.Second = types
				list = append(list, eRet)
			}
			fmt.Println(len(list))
		}
	}
	if err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage("获取失败:"+err.Error(), c)
		return
	}
	response.OkWithDetailed(response.PageResult{
		List: list,
	}, "获取成功", c)
}

func findParentId(Type int) ([]int, error) {
	switch Type {
	case 10, 11, 12, 13:
		var CategoryMap = map[int]int{
			11: 4,
			10: 3,
			12: 1,
			13: 2,
		}
		var sysType system.ProductSysType
		err := global.GVA_DB.Model(&system.ProductSysType{}).Where("category = ?", CategoryMap[Type]).First(&sysType).Error
		if err != nil {
			return nil, err
		} else {
			return []int{*sysType.Id}, nil
		}
	case 14, 15, 16:
		var CategoryMap = map[int]int{
			14: 3,
			15: 4,
			16: 2,
		}
		var sysType system.ProductSysType
		err := global.GVA_DB.Model(&system.ProductSysType{}).Where("category = ?", CategoryMap[Type]).First(&sysType).Error
		if err != nil {
			return nil, err
		}
		//再获取下
		var ids []int
		err = global.GVA_DB.Model(&system.ProductSysType{}).Where("parent_id = ?", sysType.Id).Select("id").Find(&ids).Error
		if err != nil {
			return nil, err
		}
		return ids, nil
	default:
		return nil, errors.New("错误的类别")
	}
}

// 更新类别
func (*ProductSupplyApi) UpdateCategory(c *gin.Context) {
	var req struct {
		Type int    `json:"type"`
		Name string `json:"name"`
		Id   int    `json:"id"`
		Unit string `json:"unit"`
	}
	var err = c.ShouldBindJSON(&req)
	if err != nil {
		response.FailWithMessage(err, c)
		return
	}
	// 对规格配置进行特殊处理
	if req.Type == 17 {
		dataList := []system.ProductSysType{}
		global.GVA_DB.Model(&system.ProductSysType{}).Where("category = 17 and name = ? and description = ?", req.Name, req.Unit).Find(&dataList)
		if len(dataList) > 0 {
			response.FailWithMessage("修改失败，该配置已存在", c)
			return
		}
		err = global.GVA_DB.Model(&system.ProductSysType{}).Where("id = ? and category = 17", req.Id).Updates(&system.ProductSysType{Name: req.Name, Description: req.Unit}).Error
		if err != nil {
			global.GVA_LOG.Error("更新失败", zap.Error(err))
			response.FailWithMessage("更新失败", c)
			return
		}
		response.OkWithMessage("更改成功", c)
		return
	}
	//找到对应Type的ParentId
	parentId, err := findParentId(req.Type)
	if err != nil {
		response.FailWithMessage("错误的类型", c)
		return
	}
	// 获取原来的类别
	var typeInfo system.ProductSysType
	err = global.GVA_DB.Model(&system.ProductSysType{}).Where("id = ?", req.Id).Find(&typeInfo).Error
	if err != nil {
		global.GVA_LOG.Error("查询失败", zap.Error(err))
		response.FailWithMessage("更新失败", c)
		return
	}
	dataList := []system.ProductSysType{}
	global.GVA_DB.Model(&system.ProductSysType{}).Where("parent_id = ? and name = ? and description = ?", typeInfo.ParentId, req.Name, req.Unit).Find(&dataList)
	if len(dataList) > 0 {
		response.FailWithMessage("修改失败，该配置已存在", c)
		return
	}
	// 更改数据，更新对应的Parent 下面所有同样名称的信息
	global.GVA_DB.Model(&system.ProductSysType{}).Where("parent_id in ? and name = ? ", parentId, typeInfo.Name).Updates(&system.ProductSysType{
		Name:        req.Name,
		Description: req.Unit,
	})
	response.OkWithMessage("更改成功", c)
}

func (*ProductSupplyApi) AddCategory(c *gin.Context) {
	var req struct {
		Type int    `json:"type"`
		Name string `json:"name"`
		Unit string `json:"unit"`
	}
	var err = c.ShouldBindJSON(&req)
	if err != nil {
		response.FailWithMessage(err, c)
		return
	}
	// 这里要对规格配置类型进行特殊的处理
	if req.Type == 17 {
		err = global.GVA_DB.Model(&system.ProductSysType{}).Create(&system.ProductSysType{
			Name:        req.Name,
			Description: req.Unit,
			Category:    utils.Pointer[int8](17),
		}).Error
		if err != nil {
			global.GVA_LOG.Error("添加失败", zap.Error(err))
			response.FailWithMessage("添加失败", c)
		} else {
			response.OkWithMessage("添加成功", c)
		}
		return
	}
	//找到对应Type的ParentId
	parentIds, err := findParentId(req.Type)
	if err != nil {
		response.FailWithMessage("错误的类型", c)
		return
	}
	dataList := []system.ProductSysType{}
	global.GVA_DB.Model(&system.ProductSysType{}).Where("parent_id = ? and name = ? and description = ?", parentIds[0], req.Name, req.Unit).Find(&dataList)
	if len(dataList) > 0 {
		response.FailWithMessage("创建失败，该配置已存在", c)
		return
	}
	// 添加数据，这里要添加到所有的父类别上面
	for _, parentId := range parentIds {
		global.GVA_DB.Create(&system.ProductSysType{
			ParentId:    &parentId,
			Name:        req.Name,
			Description: req.Unit,
		})
	}
	response.OkWithMessage("创建成功", c)
}
func (*ProductSupplyApi) DeleteCategory(c *gin.Context) {
	id := c.Query("id")
	// err := sysTypeService.DeleteSysType(id)
	// 这个要把所有的类型删除掉
	var typeInfo system.ProductSysType
	err := global.GVA_DB.Where("id = ?", id).First(&typeInfo).Error
	if err != nil {
		global.GVA_LOG.Error("删除失败: 找不到要删除的类", zap.Error(err))
		response.FailWithMessage("删除失败: 找不到要删除的类型", c)
		return
	}
	global.GVA_DB.Where("name = ? and category = ?", typeInfo.Name, typeInfo.Category).Delete(&system.ProductSysType{})
	response.OkWithMessage("删除成功", c)
}
func (*ProductSupplyApi) DeleteCategoryByIds(c *gin.Context) {
	ids := c.QueryArray("ids[]")
	err := global.GVA_DB.Delete(&[]system.ProductSysType{}, "id in ?", ids).Error
	if err != nil {
		global.GVA_LOG.Error("批量删除失败!", zap.Error(err))
		response.FailWithMessage("批量删除失败:"+err.Error(), c)
		return
	}
	response.OkWithMessage("批量删除成功", c)
}
