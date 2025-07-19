package image

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"
	"ygpt/server/global"
	"ygpt/server/model/common/response"

	"yunion.io/x/onecloud/pkg/mcclient/modulebase"

	"github.com/gin-gonic/gin"
	"yunion.io/x/jsonutils"
	"yunion.io/x/onecloud/pkg/mcclient/modules/k8s"
)

type ImageReposApi struct {
}

type getImageReposReq struct {
	PageIndex int64  `json:"pageIndex"`
	PageSize  int64  `json:"pageSize"`
	Name      string `json:"name"`
}

// 查询镜像仓库列表
func (p *ImageReposApi) GetImageReposList(c *gin.Context) {
	var r getImageReposReq
	err := c.ShouldBindJSON(&r)
	if err != nil {
		response.FailWithMessage(err, c)
		return
	}
	s, err := global.GetCloudClientSession()
	if err != nil {
		response.FailWithMessage(err, c)
	}
	filters := make([]string, 0)
	params := jsonutils.NewDict()
	params.Set("details", jsonutils.NewBool(true))
	params.Set("limit", jsonutils.NewInt(r.PageSize))
	params.Set("offset", jsonutils.NewInt(r.PageSize*(r.PageIndex-1)))

	if r.Name != "" {
		filters = append(filters, "name.contains('"+r.Name+"')")
	}
	params.Set("filter", jsonutils.NewStringArray(filters))
	registries, err := k8s.ContainerRegistries.List(s, params)

	response.OkWithDetailed(jsonutils.Marshal(registries).PrettyString(), "获取成功", c)
}

// 新增镜像仓库
type addImageReposReq struct {
	Name     string `json:"name"`
	Url      string `json:"url"`
	Username string `json:"username"`
	Password string `json:"password"`
}

func (p *ImageReposApi) AddImageRepos(c *gin.Context) {
	var r addImageReposReq
	err := c.ShouldBindJSON(&r)
	if err != nil {
		response.FailWithMessage(err, c)
		return
	}
	s, err := global.GetCloudClientSession()
	if err != nil {
		response.FailWithMessage("创建失败", c)
	}
	param := jsonutils.NewDict()
	param1 := jsonutils.NewDict()
	param2 := jsonutils.NewDict()
	param2.Set("username", jsonutils.NewString(r.Username))
	param2.Set("password", jsonutils.NewString(r.Password))
	param1.Set("common", param2)
	param.Set("name", jsonutils.NewString(r.Name))
	param.Set("type", jsonutils.NewString("common"))
	param.Set("url", jsonutils.NewString(r.Url))
	param.Set("config", param1)
	result, err := k8s.ContainerRegistries.Create(s, param)
	if err != nil {
		response.FailWithMessage("创建失败，该镜像仓库不可用或账号密码错误", c)
		return
	}
	response.OkWithDetailed(jsonutils.Marshal(result).PrettyString(), "创建成功", c)
}

// 删除镜像仓库
type deleteImageReposReq struct {
	ID string `json:"id"`
}

// 删除镜像仓库
type deleteImageReposByIdsReq struct {
	Ids []string `json:"ids"`
}

func (p *ImageReposApi) DeleteImageRepos(c *gin.Context) {
	var r deleteImageReposReq
	err := c.ShouldBindJSON(&r)
	if err != nil {
		response.FailWithMessage(err, c)
		return
	}
	s, err := global.GetCloudClientSession()
	if err != nil {
		response.FailWithMessage("删除失败", c)
	}
	param := jsonutils.NewDict()
	result, err := k8s.ContainerRegistries.Delete(s, r.ID, param)

	response.OkWithDetailed(jsonutils.Marshal(result).PrettyString(), "删除成功", c)
}

func (p *ImageReposApi) DeleteImageReposByIds(c *gin.Context) {
	var r deleteImageReposByIdsReq
	err := c.ShouldBindJSON(&r)
	if err != nil {
		response.FailWithMessage(err, c)
		return
	}
	s, err := global.GetCloudClientSession()
	if err != nil {
		response.FailWithMessage("删除失败", c)
	}
	for _, id := range r.Ids {

		param := jsonutils.NewDict()
		_, err := k8s.ContainerRegistries.Delete(s, id, param)
		if err != nil {
			response.FailWithMessage("删除失败", c)
			return
		}
	}
	response.OkWithMessage("删除成功", c)
}

type ContainerRegistry struct {
	Repositories []string `json:"repositories"`
}
type ImageTags struct {
	Image string   `json:"image"`
	Tags  []string `json:"tags"`
	Name  string   `json:"name"`
}
type ImageResponse struct {
	ContainerRegistry ImageTags `json:"container_registry"`
}
type ContainerRegistryResponse struct {
	ContainerRegistry ContainerRegistry `json:"container_registry"`
}
type RetImage struct {
	Name string   `json:"name"`
	Tags []string `json:"tags"`
}

// 获取镜像仓库的镜像列表
func (p *ImageReposApi) GetImageReposImageList(c *gin.Context) {
	var r deleteImageReposReq
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
	var images []RetImage
	id := r.ID
	path := fmt.Sprintf("%s/%s/images", k8s.ContainerRegistries.URLPath(), id)

	resp, err := modulebase.RawRequest(*k8s.ContainerRegistries.ResourceManager.ResourceManager, s, "GET", path, nil, nil)
	if err == nil && resp.StatusCode >= 200 && resp.StatusCode < 300 {
		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			resp.Body.Close()
			response.FailWithMessage(err, c)
			return
		}

		var containerRegistryResponse ContainerRegistryResponse
		err = json.Unmarshal(bodyBytes, &containerRegistryResponse)
		if err != nil {
			resp.Body.Close()
			response.FailWithMessage(err, c)
			return
		}

		resp.Body.Close()

		for _, repository := range containerRegistryResponse.ContainerRegistry.Repositories {
			parts := strings.Split(repository, "/")
			if len(parts) < 2 {
				continue
			}

			path = fmt.Sprintf("%s/%s/image-tags?repository=%s", k8s.ContainerRegistries.URLPath(), id, parts[1])
			resp, err = modulebase.RawRequest(*k8s.ContainerRegistries.ResourceManager.ResourceManager, s, "GET", path, nil, nil)

			if err == nil && resp.StatusCode >= 200 && resp.StatusCode < 300 {
				bodyBytes, err = io.ReadAll(resp.Body)

				var tag ImageResponse
				err = json.Unmarshal(bodyBytes, &tag)
				if err != nil {
					resp.Body.Close()
					continue
				}
				images = append(images, RetImage{
					Name: parts[0] + "/" + parts[1],
					Tags: tag.ContainerRegistry.Tags,
				})
			}
		}
	}

	response.OkWithDetailed(response.PageResult{
		List:  images,
		Total: int64(len(images)),
	}, "获取成功", c)
}
