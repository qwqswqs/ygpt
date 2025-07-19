package cloudpods

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strings"
	"ygpt/server/global"
	"yunion.io/x/jsonutils"
	"yunion.io/x/onecloud/pkg/mcclient/modulebase"
	"yunion.io/x/onecloud/pkg/mcclient/modules/image"
	"yunion.io/x/onecloud/pkg/mcclient/modules/k8s"
)

type ImageService struct {
}

type ContainerRegistry struct {
	Repositories []string `json:"repositories"`
}

type ContainerRegistryResponse struct {
	ContainerRegistry ContainerRegistry `json:"container_registry"`
}

type ImageTags struct {
	Image string   `json:"image"`
	Tags  []string `json:"tags"`
	Name  string   `json:"name"`
}

type RetImage struct {
	Name string   `json:"name"`
	Tags []string `json:"tags"`
}

type ImageResponse struct {
	ContainerRegistry ImageTags `json:"container_registry"`
}

//	func (i *ImageService) GetSystemImage(getType string) (list interface{}, total int64, err error) {
//		s, err := global.GetCloudClientSession()
//		if err != nil {
//			return
//		}
//
//		switch getType {
//		case "virtualMachine":
//			params := jsonutils.NewDict()
//			params.Set("scope", jsonutils.NewString("system"))
//			params.Set("show_fail_reason", jsonutils.NewBool(true))
//			params.Set("details", jsonutils.NewBool(true))
//			params.Set("is_guest_image", jsonutils.NewBool(false))
//
//			result, err := image.Images.List(s, params)
//		case "container":
//		case "bareMetal":
//		default:
//			return nil, 0, errors.New("参数错误")
//		}
//
//		if err != nil {
//			return
//		}
//		list = result.Data
//		total = int64(result.Total)
//		return
//	}
func (i *ImageService) GetSystemImage() (list interface{}, total int64, err error) {
	s, err := global.GetCloudClientSession()
	if err != nil {
		return
	}

	params := jsonutils.NewDict()
	params.Set("scope", jsonutils.NewString("system"))
	params.Set("show_fail_reason", jsonutils.NewBool(true))
	params.Set("details", jsonutils.NewBool(true))
	params.Set("is_guest_image", jsonutils.NewBool(false))

	result, err := image.Images.List(s, params)

	if err != nil {
		return
	}
	list = result.Data
	total = int64(result.Total)
	return
}

func (i *ImageService) FindSystemImage(name string) (jsonutils.JSONObject, error) {
	s, err := global.GetCloudClientSession()
	if err != nil {
		return nil, err
	}

	return image.Images.GetByName(s, name, nil)
}

func (i *ImageService) GetDockerImage() (interface{}, int64, error) {
	s, err := global.GetCloudClientSession()
	if err != nil {
		return nil, 0, err
	}

	params := jsonutils.NewDict()
	params.Set("details", jsonutils.NewBool(true))

	registries, err := k8s.ContainerRegistries.List(s, params)
	if err != nil {
		return nil, 0, err
	}

	var images []RetImage
	for _, registry := range registries.Data {
		id, err := registry.GetString("id")
		if err != nil {
			return nil, 0, err
		}
		path := fmt.Sprintf("%s/%s/images", k8s.ContainerRegistries.URLPath(), id)

		resp, err := modulebase.RawRequest(*k8s.ContainerRegistries.ResourceManager.ResourceManager, s, "GET", path, nil, nil)
		if err == nil && resp.StatusCode >= 200 && resp.StatusCode < 300 {
			bodyBytes, err := io.ReadAll(resp.Body)
			if err != nil {
				resp.Body.Close()
				return nil, 0, err
			}

			var containerRegistryResponse ContainerRegistryResponse
			err = json.Unmarshal(bodyBytes, &containerRegistryResponse)
			if err != nil {
				resp.Body.Close()
				return nil, 0, err
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
						Name: parts[1],
						Tags: tag.ContainerRegistry.Tags,
					})
				}
			}
		}
	}
	return images, int64(len(images)), nil
}

func (i *ImageService) FindDockerImage(name string) (string, error) {
	var (
		imageUrl string
	)

	s, err := global.GetCloudClientSession()
	if err != nil {
		return "", err
	}

	params := jsonutils.NewDict()
	params.Set("details", jsonutils.NewBool(true))

	registries, err := k8s.ContainerRegistries.List(s, params)
	if err != nil {
		return "", err
	}

	for _, registry := range registries.Data {
		id, err := registry.GetString("id")
		if err != nil {
			return "", err
		}
		path := fmt.Sprintf("%s/%s/images", k8s.ContainerRegistries.URLPath(), id)

		resp, err := modulebase.RawRequest(*k8s.ContainerRegistries.ResourceManager.ResourceManager, s, "GET", path, nil, nil)
		if err == nil && resp.StatusCode >= 200 && resp.StatusCode < 300 {
			bodyBytes, err := io.ReadAll(resp.Body)
			if err != nil {
				resp.Body.Close()
				return "", err
			}

			var containerRegistryResponse ContainerRegistryResponse
			err = json.Unmarshal(bodyBytes, &containerRegistryResponse)
			if err != nil {
				resp.Body.Close()
				return "", err
			}

			resp.Body.Close()

			for _, repository := range containerRegistryResponse.ContainerRegistry.Repositories {
				parts := strings.Split(repository, "/")
				if len(parts) < 2 {
					continue
				}

				if parts[1] != name {
					continue
				}

				path = fmt.Sprintf("%s/%s/image-tags?repository=%s", k8s.ContainerRegistries.URLPath(), id, parts[1])
				resp, err = modulebase.RawRequest(*k8s.ContainerRegistries.ResourceManager.ResourceManager, s, "GET", path, nil, nil)

				if err == nil && resp.StatusCode >= 200 && resp.StatusCode < 300 {
					bodyBytes, err = io.ReadAll(resp.Body)

					var imageTagRet ImageResponse
					err = json.Unmarshal(bodyBytes, &imageTagRet)
					if err != nil {
						resp.Body.Close()
						continue
					}

					if len(imageTagRet.ContainerRegistry.Tags) == 0 {
						resp.Body.Close()
						continue
					}

					imageUrl = imageTagRet.ContainerRegistry.Image + ":" + imageTagRet.ContainerRegistry.Tags[0]
					return imageUrl, nil
				}
			}
		}
	}

	return "", errors.New("image not found")
}
