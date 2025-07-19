package baseRes

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"net"
	"net/url"
	"ygpt/server/global"
	"ygpt/server/model/common/response"
	data "ygpt/server/model/yunguan/k8s"
	"yunion.io/x/jsonutils"
	"yunion.io/x/onecloud/pkg/mcclient/modules/k8s"
)

type BaseK8SApi struct {
}
type getBaseK8SReq struct {
	PageIndex int64  `json:"pageIndex"`
	PageSize  int64  `json:"pageSize"`
	Name      string `json:"name"`
	Status    string `json:"status"`
}

// 获取集群列表
func (d *BaseK8SApi) GetBaseK8SList(c *gin.Context) {
	var r getBaseK8SReq
	err := c.ShouldBindJSON(&r)
	if err != nil {
		response.FailWithMessage(err, c)
		return
	}
	s, err := global.GetCloudClientSession()
	if err != nil {
		response.FailWithMessage(err, c)
		return
	}
	params := jsonutils.NewDict()
	params.Set("limit", jsonutils.NewInt(r.PageSize))
	params.Set("show_fail_reason", jsonutils.NewBool(true))
	params.Set("details", jsonutils.NewBool(true))
	params.Set("scope", jsonutils.NewString("system"))
	params.Set("offset", jsonutils.NewInt(r.PageSize*(r.PageIndex-1)))

	filters := make([]string, 0)

	if r.Name != "" {
		filters = append(filters, "name.contains('"+r.Name+"')")
	}

	if r.Status != "" {
		filters = append(filters, "status.equals("+r.Status+")")
	}

	params.Set("filter", jsonutils.NewStringArray(filters))

	result, err := k8s.KubeClusters.List(s, params)
	if err != nil {
		response.FailWithMessage(err, c)
		return
	}
	renterIds := make([]data.K8sModel, 0)
	list := jsonutils.Marshal(result).Interface()
	rawMap, _ := list.(map[string]interface{})
	if arr, ok := rawMap["data"].([]interface{}); ok {
		for _, item := range arr {
			if obj, i := item.(map[string]interface{}); i {
				if id, j := obj["id"].(string); j {
					var model data.K8sModel
					global.GVA_DB.Where("k8s_id = ?", id).First(&model)
					renterIds = append(renterIds, model)
				}
			}
		}
	}
	for i, item := range rawMap["data"].([]interface{}) {
		if obj, ok := item.(map[string]interface{}); ok {
			// 将 renterId 添加到当前对象中
			obj["serverIp"] = renterIds[i].ApiIP
			obj["port"] = renterIds[i].Endpoint
		} else {
			fmt.Printf("Item at index %d is not a map[string]interface{}\n", i)
		}
	}

	response.OkWithDetailed(rawMap, "获取成功", c)
}

type deleteReq struct {
	ID string `json:"id"`
}
type updateReq struct {
	ID       string `json:"id"`
	Endpoint string `json:"endpoint"`
}
type deleteByIds struct {
	Ids []string `json:"ids"`
}

// 删除集群
func (d *BaseK8SApi) DeleteBaseCluster(c *gin.Context) {
	var r deleteReq
	err := c.ShouldBindJSON(&r)
	if err != nil {
		response.FailWithMessage(err, c)
		return
	}
	s, err := global.GetCloudClientSession()
	if err != nil {
		response.FailWithMessage(err, c)
		return
	}
	params := jsonutils.NewDict()
	params.Set("id", jsonutils.NewString(r.ID))
	result, err := k8s.KubeClusters.Delete(s, r.ID, params)
	if err != nil {
		response.FailWithMessage(err, c)
		return
	}
	err = global.GVA_DB.Model(&data.K8sModel{}).Where("k8s_id", r.ID).Delete(&data.K8sModel{}).Error
	if err != nil {
		response.FailWithMessage(err, c)
		return
	}
	response.OkWithDetailed(jsonutils.Marshal(result).PrettyString(), "删除成功", c)
}

// 编辑集群端口
func (d *BaseK8SApi) UpdateBaseCluster(c *gin.Context) {
	var r updateReq
	err := c.ShouldBindJSON(&r)
	if err != nil {
		response.FailWithMessage(err, c)
		return
	}
	err = global.GVA_DB.Model(&data.K8sModel{}).Where("k8s_id", r.ID).Updates(&data.K8sModel{
		Endpoint: r.Endpoint,
	}).Error
	if err != nil {
		response.FailWithMessage(err, c)
		return
	}
	response.OkWithMessage("修改成功", c)
}

// 批量删除集群
func (d *BaseK8SApi) DeleteBaseClusterByIds(c *gin.Context) {
	var r deleteByIds
	err := c.ShouldBindJSON(&r)
	if err != nil {
		response.FailWithMessage(err, c)
		return
	}
	s, err := global.GetCloudClientSession()
	if err != nil {
		response.FailWithMessage(err, c)
		return
	}
	for _, id := range r.Ids {

		params := jsonutils.NewDict()
		params.Set("id", jsonutils.NewString(id))
		_, err := k8s.KubeClusters.Delete(s, id, params)
		if err != nil {
			response.FailWithMessage(err, c)
			return
		}
		err = global.GVA_DB.Model(&data.K8sModel{}).Where("k8s_id", id).Delete(&data.K8sModel{}).Error
		if err != nil {
			response.FailWithMessage(err, c)
			return
		}
	}

	response.OkWithMessage("删除成功", c)
}

// 导入集群
type importClusterReq struct {
	Name       string `json:"name"`
	KubeConfig string `json:"kubeConfig"`
}
type ErrorResponse struct {
	Error struct {
		Class   string `json:"class"`
		Code    int    `json:"code"`
		Details string `json:"details"`
		Request struct {
			Body    string            `json:"body"`
			Headers map[string]string `json:"headers"`
			Method  string            `json:"method"`
			URL     string            `json:"url"`
		} `json:"request"`
	} `json:"error"`
}

func (d *BaseK8SApi) ImportBaseCluster(c *gin.Context) {
	var r importClusterReq
	err := c.ShouldBindJSON(&r)
	if err != nil {
		response.FailWithMessage(err, c)
		return
	}
	s, err := global.GetCloudClientSession()
	if err != nil {
		response.FailWithMessage(err, c)
		return
	}
	params2 := jsonutils.NewDict()
	params2.Set("kubeconfig", jsonutils.NewString(r.KubeConfig))
	params := jsonutils.NewDict()
	params.Set("mode", jsonutils.NewString("import"))
	params.Set("name", jsonutils.NewString(r.Name))
	params.Set("project_domain_id", jsonutils.NewString("default"))
	params.Set("provider", jsonutils.NewString("external"))
	params.Set("resource_type", jsonutils.NewString("unknown"))
	params.Set("import_data", params2)
	result, err := k8s.KubeClusters.Create(s, params)
	if err != nil {
		var errorResponse ErrorResponse

		// 解析 JSON 字符串
		json.Unmarshal([]byte(err.Error()), &errorResponse)
		message := ""
		if errorResponse.Error.Code == 409 {
			message = "kubecluster资源重名"
		} else {
			message = "导入失败，请检查配置"
		}
		response.FailWithMessage(message, c)
		return
	}

	list := jsonutils.Marshal(result).Interface()
	rawMap, _ := list.(map[string]interface{})
	urlData := rawMap["api_server"].(string)
	parsedURL, err := url.Parse(urlData)
	hostData := parsedURL.Host
	ip, port, err := net.SplitHostPort(hostData)
	modelData := data.K8sModel{
		ApiIP:    ip,
		K8sId:    rawMap["id"].(string),
		Endpoint: port,
		Name:     rawMap["name"].(string),
	}
	err = global.GVA_DB.Create(&modelData).Error
	if err != nil {
		response.FailWithMessage(err, c)
	}
	response.OkWithDetailed(jsonutils.Marshal(result).Interface(), "创建成功", c)
}
