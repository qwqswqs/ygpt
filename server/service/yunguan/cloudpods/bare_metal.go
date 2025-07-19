package cloudpods

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
	"ygpt/server/global"
	"ygpt/server/model/yunguan/product"
	"ygpt/server/model/yunguan/renter"

	"yunion.io/x/onecloud/pkg/mcclient/modules/scheduler"

	"yunion.io/x/jsonutils"
	"yunion.io/x/onecloud/pkg/mcclient/modules/compute"
	"yunion.io/x/pkg/util/printutils"
)

type BareMetalService struct {
}

func (v *BareMetalService) HostDelete(id string) error {
	s, err := global.GetCloudClientSession()
	if err != nil {
		return err
	}

	params := jsonutils.NewDict()
	params.Set("protected", jsonutils.NewBool(false))
	params.Set("disable_delete", jsonutils.NewBool(false))
	_, err = compute.Servers.Update(s, id, params)
	if err != nil {
		return err
	}
	_, err = compute.Servers.Delete(s, id, nil)
	if err != nil {
		return err
	}
	return nil
}

func (v *BareMetalService) HostGet(id string) (jsonutils.JSONObject, error) {
	s, err := global.GetCloudClientSession()
	if err != nil {
		return nil, err
	}

	return compute.Baremetalagents.Get(s, id, nil)
}

func (v *BareMetalService) HostGetLoginInfo(id string) (string, string, error) {
	s, err := global.GetCloudClientSession()
	if err != nil {
		return "", "", err
	}

	loginInfo, err := compute.Servers.GetLoginInfo(s, id, nil)
	if err != nil {
		return "", "", err
	}

	userName, _ := loginInfo.GetString("username")
	password, _ := loginInfo.GetString("password")

	return userName, password, nil
}

func (v *BareMetalService) HostListByIds(getType string, id []string) ([]jsonutils.JSONObject, int, error) {
	s, err := global.GetCloudClientSession()
	if err != nil {
		return nil, 0, err
	}

	result := &printutils.ListResult{}
	params := jsonutils.NewDict()

	switch getType {
	case "bareMetal":
		params.Set("scope", jsonutils.NewString("system"))
		params.Set("show_fail_reason", jsonutils.NewBool(true))
		params.Set("details", jsonutils.NewBool(true))
		params.Set("hypervisor", jsonutils.NewString("baremetal"))
		params.Set("with_meta", jsonutils.NewBool(true))
		if id != nil && len(id) > 0 {
			params.Set("id", jsonutils.NewStringArray(id))
		}

		result, err = compute.Servers.List(s, params)
	case "physicalMachine":
		params.Set("scope", jsonutils.NewString("system"))
		params.Set("show_fail_reason", jsonutils.NewBool(true))
		params.Set("details", jsonutils.NewBool(true))
		params.Set("baremetal", jsonutils.NewBool(true))
		params.Set("host_type", jsonutils.NewString("baremetal"))
		params.Set("with_meta", jsonutils.NewBool(true))
		if id != nil && len(id) > 0 {
			params.Set("id", jsonutils.NewStringArray(id))
		}

		result, err = compute.Hosts.List(s, params)
	default:
		return nil, 0, errors.New("参数错误")
	}

	if err != nil {
		return nil, 0, err
	}

	return result.Data, result.Total, nil
}

func (v *BareMetalService) HostListByRenterID(getType string, renterId *int, pageIndex *int, pageSize *int) ([]jsonutils.JSONObject, int64, error) {
	s, err := global.GetCloudClientSession()
	if err != nil {
		return nil, 0, err
	}

	result := &printutils.ListResult{}
	params := jsonutils.NewDict()

	switch getType {
	case "bareMetal":
		params.Set("scope", jsonutils.NewString("system"))
		params.Set("show_fail_reason", jsonutils.NewBool(true))
		params.Set("details", jsonutils.NewBool(true))
		params.Set("hypervisor", jsonutils.NewString("baremetal"))
		params.Set("with_meta", jsonutils.NewBool(true))
		if renterId != nil {
			params.Set("filter", jsonutils.NewString("name.contains('"+strconv.Itoa(*renterId)+"')"))
		}

		if pageIndex != nil && pageSize != nil {
			params.Set("limit", jsonutils.NewInt(int64(*pageSize)))
			params.Set("offset", jsonutils.NewInt(int64(*pageSize*(*pageIndex-1))))
		}

		result, err = compute.Servers.List(s, params)
	case "physicalMachine":
		params.Set("scope", jsonutils.NewString("system"))
		params.Set("show_fail_reason", jsonutils.NewBool(true))
		params.Set("details", jsonutils.NewBool(true))
		params.Set("baremetal", jsonutils.NewBool(true))
		params.Set("host_type", jsonutils.NewString("baremetal"))
		params.Set("with_meta", jsonutils.NewBool(true))
		if renterId != nil {
			params.Set("filter", jsonutils.NewString("name.contains('"+strconv.Itoa(*renterId)+"')"))
		}

		if pageIndex != nil && pageSize != nil {
			params.Set("limit", jsonutils.NewInt(int64(*pageSize)))
			params.Set("offset", jsonutils.NewInt(int64(*pageSize*(*pageIndex-1))))
		}

		result, err = compute.Hosts.List(s, params)
	default:
		return nil, 0, errors.New("参数错误")
	}

	if err != nil {
		return nil, 0, err
	}

	return result.Data, int64(result.Total), nil
}

func (v *BareMetalService) HostList(getType string) ([]jsonutils.JSONObject, int, error) {
	s, err := global.GetCloudClientSession()
	if err != nil {
		return nil, 0, err
	}

	result := &printutils.ListResult{}
	params := jsonutils.NewDict()

	switch getType {
	case "bareMetal":
		params.Set("scope", jsonutils.NewString("system"))
		params.Set("show_fail_reason", jsonutils.NewBool(true))
		params.Set("details", jsonutils.NewBool(true))
		params.Set("hypervisor", jsonutils.NewString("baremetal"))
		params.Set("with_meta", jsonutils.NewBool(true))

		result, err = compute.Servers.List(s, params)
	case "physicalMachine":
		params.Set("scope", jsonutils.NewString("system"))
		params.Set("show_fail_reason", jsonutils.NewBool(true))
		params.Set("details", jsonutils.NewBool(true))
		params.Set("baremetal", jsonutils.NewBool(true))
		params.Set("host_type", jsonutils.NewString("baremetal"))
		params.Set("with_meta", jsonutils.NewBool(true))

		result, err = compute.Hosts.List(s, params)
	default:
		return nil, 0, errors.New("参数错误")
	}

	if err != nil {
		return nil, 0, err
	}

	return result.Data, result.Total, nil
}

type BaremetalCreateReq struct {
	ProjectId    string `json:"project_id"`    //
	Count        int    `json:"count"`         //
	VmemSize     int    `json:"vmem_size"`     //
	VcpuCount    int    `json:"vcpu_count"`    //
	GenerateName string `json:"generate_name"` //
	Hypervisor   string `json:"hypervisor"`    //
	AutoStart    bool   `json:"auto_start"`    //
	Vdi          string `json:"vdi"`           //
	Disks        []struct {
		Size       int64  `json:"size"`
		ImageId    string `json:"image_id,omitempty"`
		Mountpoint string `json:"mountpoint,omitempty"`
	} `json:"disks"`
	BaremetalDiskConfigs []BaremetalDiskConfig `json:"baremetal_disk_configs"`
	Nets                 []struct {
		Network string `json:"network"`
		Exit    bool   `json:"exit"`
	} `json:"nets"` //
	PreferHost      string   `json:"prefer_host"`   //
	PreferRegion    string   `json:"prefer_region"` //
	PreferZone      string   `json:"prefer_zone"`   //
	IsolatedDevices []device `json:"isolated_devices"`
}

type BaremetalDiskConfig struct {
	Conf    string `json:"conf"`
	Driver  string `json:"driver"`
	Count   int    `json:"count"`
	Range   []int  `json:"range"`
	Adapter int    `json:"adapter"`
	Type    string `json:"type"`
}

//type IsolatedDevice struct {
//	Model   string `json:"model"`
//	Vendor  string `json:"vendor"`
//	DevType string `json:"dev_type"`
//	PciId   string `json:"pci_id"`
//}

func (v *BareMetalService) HostCreate(productId int, imageName string, renterId int) (jsonutils.JSONObject, error) {
	var r BaremetalCreateReq

	s, err := global.GetCloudClientSession()
	if err != nil {
		return nil, err
	}

	var productComputing product.ProductComputing
	if err := global.GVA_DB.Model(&productComputing).Where("id = ?", productId).First(&productComputing).Error; err != nil {
		return nil, fmt.Errorf("未找到对应的云算力产品，productId为%d", productId)
	}

	r.PreferRegion = "default"
	r.PreferZone = GetZoneID()
	r.ProjectId = s.GetProjectId()
	r.Count = 1
	r.Hypervisor = "baremetal"
	r.AutoStart = true
	r.Vdi = "vnc"

	params := jsonutils.NewDict()
	params.Set("scope", jsonutils.NewString("system"))
	params.Set("show_fail_reason", jsonutils.NewBool(true))
	params.Set("details", jsonutils.NewBool(true))
	params.Set("baremetal", jsonutils.NewBool(true))
	params.Set("host_type", jsonutils.NewString("baremetal"))
	params.Set("with_meta", jsonutils.NewBool(true))

	//寻找物理机
	baremetal, err := compute.Hosts.GetById(s, productComputing.ResourceId, params)
	if err != nil {
		return nil, err
	}

	spec, err := baremetal.Get("spec")
	if err != nil {
		return nil, err
	}

	cpu, err := spec.Int("cpu")
	if err != nil {
		return nil, err
	}
	cpuStr := "cpu:" + strconv.FormatInt(cpu, 10) + "/"

	mem, err := spec.Int("mem")
	if err != nil {
		return nil, err
	}
	memStr := "mem:" + strconv.FormatInt(mem, 10) + "M/"

	r.VmemSize = int(mem)
	r.VcpuCount = int(cpu)
	nicCount, err := baremetal.Int("nic_count")
	if err != nil {
		return nil, err
	}

	nicStr := "nic:" + strconv.FormatInt(nicCount-1, 10)

	disk, err := spec.Get("disk")
	if err != nil {
		return nil, err
	}
	diskMap, err := disk.GetMap()
	if err != nil {
		return nil, err
	}

	diskStr := ""
	for driverName, driverObj := range diskMap {
		diskInfs, err := driverObj.GetArray("adapter0")
		if err != nil {
			return nil, err
		}
		if len(diskInfs) <= 0 {
			continue
		}

		diskInfo := diskInfs[0]
		count, err := diskInfo.Int("count")
		if err != nil {
			return nil, err
		}
		size, err := diskInfo.Int("size")
		if err != nil {
			return nil, err
		}
		diskType, err := diskInfo.GetString("type")
		if err != nil {
			return nil, err
		}

		diskStr = diskStr + "disk:" + driverName + "_adapter0_" + diskType + "_" + strconv.FormatInt(size/1024, 10) + "Gx" + strconv.FormatInt(count, 10) + "/"
	}

	sysInfo, err := baremetal.Get("sys_info")
	if err != nil {
		return nil, err
	}

	manufacture, err := sysInfo.GetString("manufacture")
	if err != nil {
		return nil, err
	}

	// 将空格替换为下划线
	manufacture = strings.ReplaceAll(manufacture, " ", "_")

	manufactureStr := "manufacture:" + manufacture + "/"

	model, err := sysInfo.GetString("model")
	if err != nil {
		return nil, err
	}
	// 将空格替换为下划线
	model = strings.ReplaceAll(model, " ", "_")

	// 构建 modelStr
	modelStr := "model:" + model + "/"

	skuStr := cpuStr + diskStr + manufactureStr + memStr + modelStr + nicStr

	podListParams := jsonutils.NewDict()
	podListParams.Set("show_emulated", jsonutils.NewBool(true))
	podListParams.Set("resource_type", jsonutils.NewString("shared"))
	//podListParams.Set("project_domain", jsonutils.NewString("default"))
	podListParams.Set("host_type", jsonutils.NewString("baremetal"))

	Capabilities, err := compute.Capabilities.List(s, podListParams)
	if err != nil {
		return nil, err
	}

	if len(Capabilities.Data) > 0 {
		capabilitie := Capabilities.Data[0]
		specs, err := capabilitie.Get("specs")
		if err != nil {
			return nil, err
		}

		hosts, err := specs.Get("hosts")
		if err != nil {
			return nil, err
		}

		if !hosts.Contains(skuStr) {
			return nil, errors.New("未找到物理机")
		}

		isolatedDevices, _ := hosts.GetArray("isolated_devices")

		jsonutils.Update(&r.IsolatedDevices, isolatedDevices)
	}

	//赋值磁盘参数
	diskMap, err = spec.GetMap("disk")
	if err != nil {
		return nil, err
	}

	driver, err := spec.GetString("driver")
	if err != nil {
		return nil, err
	}

	strogeType, err := baremetal.GetString("storage_type")
	if err != nil {
		return nil, err
	}

	//寻找镜像
	imageService := &ImageService{}
	sysImage, err := imageService.FindSystemImage(imageName)
	if err != nil {
		return nil, err
	}

	imageId, err := sysImage.GetString("id")
	if err != nil {
		return nil, err
	}

	for key, value := range diskMap {
		if key == driver {
			diskInfos, err := value.GetArray("adapter0")
			if err != nil {
				return nil, err
			}
			if len(diskInfos) > 0 {
				diskInfo := diskInfos[0]
				size, err := diskInfo.Int("size")
				if err != nil {
					return nil, err
				}

				startIndex, err := diskInfo.Int("start_index")
				if err != nil {
					return nil, err
				}

				r.Disks = []struct {
					Size       int64  `json:"size"`
					ImageId    string `json:"image_id,omitempty"`
					Mountpoint string `json:"mountpoint,omitempty"`
				}{
					{
						Size:    size - 30*1024,
						ImageId: imageId,
					},
					{
						Size:       -1,
						Mountpoint: "/(系统)",
					},
				}

				var baremetalDiskConfig BaremetalDiskConfig

				baremetalDiskConfigs := make([]BaremetalDiskConfig, 0)

				baremetalDiskConfig.Conf = "none"
				baremetalDiskConfig.Type = strogeType
				baremetalDiskConfig.Driver = driver
				baremetalDiskConfig.Count = 1
				baremetalDiskConfig.Adapter = 0
				baremetalDiskConfig.Range = append(baremetalDiskConfig.Range, int(startIndex))

				baremetalDiskConfigs = append(baremetalDiskConfigs, baremetalDiskConfig)

				r.BaremetalDiskConfigs = baremetalDiskConfigs

				break
			}
		}
	}

	if r.BaremetalDiskConfigs == nil || len(r.BaremetalDiskConfigs) == 0 {
		return nil, errors.New("磁盘参数错误")
	}

	//寻找网络
	networkService := &NetworkService{}
	networkId, err := networkService.GetGuestOneNetwork()
	if err != nil {
		return nil, err
	}
	r.Nets = []struct {
		Network string `json:"network"`
		Exit    bool   `json:"exit"`
	}{
		{
			Network: networkId,
			Exit:    false,
		},
	}

	r.GenerateName = global.GenerateUniqueName(renterId)

	_, err = scheduler.SchedManager.DoForecast(s, jsonutils.Marshal(r))
	if err != nil {
		return nil, err
	}

	result, err := compute.Servers.Create(s, jsonutils.Marshal(r))
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (v *BareMetalService) HostForecast(bareMetalRet jsonutils.JSONObject) error {
	var errMsg string
	if bareMetalRet == nil {
		return errors.New("参数为空")
	}

	s, err := global.GetCloudClientSession()
	if err != nil {
		return errors.New("获取CloudClientSession失败")
	}

	bareMetalId, err := bareMetalRet.GetString("id")
	if err != nil {
		errMsg = "获取裸金属id失败"
	}

	status, err := bareMetalRet.GetString("status")
	if err != nil {
		errMsg = "获取裸金属状态失败"
	}

	params := jsonutils.NewDict()
	params.Set("scope", jsonutils.NewString("system"))
	params.Set("show_fail_reason", jsonutils.NewBool(true))
	params.Set("details", jsonutils.NewBool(true))
	params.Set("hypervisor", jsonutils.NewString("baremetal"))
	params.Set("with_meta", jsonutils.NewBool(true))
	params.Set("id", jsonutils.NewString(bareMetalId))

	if status != "running" {
		size := 3000
		for i := 0; i < size; i++ {
			bareMetalRet, err = compute.Servers.Get(s, bareMetalId, params)
			if err != nil {
				errMsg = "获取裸金属信息失败"
			}
			if bareMetalRet == nil {
				continue
			}
			status, _ = bareMetalRet.GetString("status")
			if status == "running" {
				break
			}
			if i < size-1 {
				time.Sleep(2 * time.Second)
			}
		}
	}

	if status == "running" {
		ip, err := bareMetalRet.GetString("ips")
		if err != nil {
			errMsg = "获取裸金属ip失败"
		}

		username, password, err := v.HostGetLoginInfo(bareMetalId)
		if err != nil {
			errMsg = "获取裸金属登录信息失败"
		}

		var resInfo renter.RenterRes
		err = global.GVA_DB.Model(&renter.RenterRes{}).Where("resource_id = ?", bareMetalId).First(&resInfo).Error
		if err != nil {
			errMsg = "未找到租户资源"
		}

		resInfo.SshHost = ip
		resInfo.SshPort = 22
		resInfo.SshUser = username
		resInfo.SshPwd = password
		resInfo.Status = 2
		// 资源创建需要时间，资源确实创建完成才能确定开始和结束时间
		resInfo.StartTime = time.Now()
		resInfo.EndTime = time.Now().Add(time.Duration(resInfo.UseTime) * time.Hour)
		err = global.GVA_DB.Model(&renter.RenterRes{}).Where("resource_id = ?", bareMetalId).Updates(&resInfo).Error
		if err != nil {
			errMsg = "更新租户资源失败"
		}
	} else {
		errMsg = "裸金属未正常运行，请和运维人员联系"
	}

	if errMsg == "" {
		return nil
	}

	return errors.New(errMsg)
}

func (v *BareMetalService) HostForecastById(id string) error {
	var errMsg string
	if id == "" {
		return errors.New("参数为空")
	}

	s, err := global.GetCloudClientSession()
	if err != nil {
		return errors.New("获取CloudClientSession失败")
	}

	params := jsonutils.NewDict()
	params.Set("scope", jsonutils.NewString("system"))
	params.Set("show_fail_reason", jsonutils.NewBool(true))
	params.Set("details", jsonutils.NewBool(true))
	params.Set("hypervisor", jsonutils.NewString("baremetal"))
	params.Set("with_meta", jsonutils.NewBool(true))
	params.Set("id", jsonutils.NewString(id))

	status := ""
	size := 3000
	var bareMetalRet jsonutils.JSONObject
	for i := 0; i < size; i++ {
		bareMetalRet, err = compute.Servers.Get(s, id, params)
		if err != nil {
			errMsg = "获取裸金属信息失败"
		}
		if bareMetalRet == nil {
			continue
		}
		status, _ = bareMetalRet.GetString("status")
		if status == "running" {
			break
		}
		if i < size-1 {
			time.Sleep(2 * time.Second)
		}
	}

	if status == "running" {
		ip, err := bareMetalRet.GetString("ips")
		if err != nil {
			errMsg = "获取裸金属ip失败"
		}

		username, password, err := v.HostGetLoginInfo(id)
		if err != nil {
			errMsg = "获取裸金属登录信息失败"
		}

		var resInfo renter.RenterRes
		err = global.GVA_DB.Model(&renter.RenterRes{}).Where("resource_id = ?", id).First(&resInfo).Error
		if err != nil {
			errMsg = "未找到租户资源"
		}

		resInfo.SshHost = ip
		resInfo.SshPort = 22
		resInfo.SshUser = username
		resInfo.SshPwd = password
		resInfo.Status = 2
		// 资源创建需要时间，资源确实创建完成才能确定开始和结束时间
		resInfo.StartTime = time.Now()
		resInfo.EndTime = time.Now().Add(time.Duration(resInfo.UseTime) * time.Hour)
		err = global.GVA_DB.Model(&renter.RenterRes{}).Where("resource_id = ?", id).Updates(&resInfo).Error
		if err != nil {
			errMsg = "更新租户资源失败"
		}
	} else {
		errMsg = "裸金属未正常运行，请和运维人员联系"
	}

	if errMsg == "" {
		return nil
	}

	return errors.New(errMsg)
}

func GetZoneID() string {

	s, err := global.GetCloudClientSession()
	if err != nil {
		return ""
	}
	params := jsonutils.NewDict()
	result, err := compute.Zones.List(s, params)
	if err != nil {
		return ""
	}
	jsonString := jsonutils.Marshal(result).PrettyString()
	startIndex := strings.Index(jsonString, `"id": "`) + len(`"id": "`)
	if startIndex == -1+len(`"id": "`) {
		fmt.Println("id字段未找到")
		return ""
	}

	// 查找id字段的结束位置（即下一个双引号的位置）
	endIndex := strings.Index(jsonString[startIndex:], `"`)
	if endIndex == -1 {
		fmt.Println("id字段值未正确结束")
		return ""
	}
	endIndex += startIndex

	// 提取id字段的值
	id := jsonString[startIndex:endIndex]
	fmt.Println("ID:", id)
	return id
}
