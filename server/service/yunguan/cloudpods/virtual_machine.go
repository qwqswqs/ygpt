package cloudpods

import (
	"errors"
	"strconv"
	"time"
	"ygpt/server/global"
	"ygpt/server/model/yunguan/renter"

	"yunion.io/x/jsonutils"
	"yunion.io/x/onecloud/pkg/mcclient/modules/compute"
	"yunion.io/x/pkg/util/regutils"
)

type VirtualMachineService struct {
}

type diskData struct {
	Backend  string `json:"backend"`
	DiskType string `json:"disk_type"`
	ImageId  string `json:"image_id"`
	Index    int    `json:"index"`
	//Medium   string `json:"medium"`
	Size int `json:"size"`
}

type netData struct {
	Network string `json:"network"`
}

type device struct {
	Model    string `json:"model"`
	Vendor   string `json:"vendor"`
	Dev_type string `json:"dev_type"`
}

type hostCreateReq struct {
	AutoStart       bool       `json:"auto_start"`
	DeployTelegraf  bool       `json:"deploy_telegraf"`
	Description     string     `json:"description"`
	Disks           []diskData `json:"disks"`
	GenerateName    string     `json:"generate_name"`
	Hypervisor      string     `json:"hypervisor"`
	Machine         string     `json:"machine"`
	Nets            []netData  `json:"nets"`
	OsArch          string     `json:"os_arch"`
	PreferRegion    string     `json:"prefer_region"`
	PreferZone      string     `json:"prefer_zone"`
	ProjectId       string     `json:"project_id"`
	Sku             string     `json:"sku"`
	VcpuCount       int        `json:"vcpu_count"`
	Vdi             string     `json:"vdi"`
	Vga             string     `json:"vga"`
	VmemSize        int        `json:"vmem_size"`
	Count           int        `json:"__count__"`
	Cdrom           string     `json:"cdrom"`
	IsolatedDevices []device   `json:"isolated_devices"`
}

func (v *VirtualMachineService) HostCreate(imageName string, cpu, memory int, storageType string, sysDisk int, dataDisk int, gpu int, renterId int) (jsonutils.JSONObject, error) {
	s, err := global.GetCloudClientSession()
	if err != nil {
		return nil, err
	}

	var r hostCreateReq

	r.AutoStart = true
	r.DeployTelegraf = true
	r.Description = ""
	r.Count = 1
	r.GenerateName = global.GenerateUniqueName(renterId)
	r.Hypervisor = "kvm"
	r.PreferRegion = "default"
	r.VcpuCount = cpu
	r.VmemSize = memory * 1024

	r.Sku = "ecs.g1.c" + strconv.Itoa(cpu) + "m" + strconv.Itoa(memory)

	_, err = compute.ServerSkus.GetByName(s, r.Sku, nil)
	if err != nil { // sku不存在，创建sku
		skuParams := jsonutils.NewDict()
		skuParams.Set("cpu_core_count", jsonutils.NewInt(int64(cpu)))
		skuParams.Set("memory_size_mb", jsonutils.NewInt(int64(memory*1024)))
		skuInfo, err := compute.ServerSkus.Create(s, skuParams)
		if err != nil {
			return nil, err
		}

		r.Sku, err = skuInfo.GetString("name")
		if err != nil {
			return nil, err
		}
	}
	//r.PreferZone = "2e74be3a-42e0-42a9-827e-6dcd4bb8a0e1"
	//r.ProjectId = "1ac77f18bd49409089cedd0572964728"

	imageService := &ImageService{}
	sysImage, err := imageService.FindSystemImage(imageName)
	if err != nil {
		return nil, err
	}

	imageId, err := sysImage.GetString("id")
	if err != nil {
		return nil, err
	}
	diskFormat, err := sysImage.GetString("disk_format")
	if err != nil {
		return nil, err
	}

	r.OsArch = "x86"
	osArch, err := sysImage.GetString("os_arch")
	if err == nil {
		r.OsArch = osArch
	}

	if storageType != "nfs" && storageType != "local" {
		return nil, errors.New("暂不支持该存储类型")
	}

	// 初始化系统磁盘
	sysDiskData := diskData{
		Backend:  storageType,
		DiskType: "sys",
		ImageId:  imageId,
		Index:    0,
		//Medium:   "rotate",
		Size: sysDisk * 1024,
	}

	// 根据 dataDisk 的值决定是否添加数据磁盘
	if dataDisk > 0 {
		dataDiskData := diskData{
			Backend:  storageType,
			DiskType: "data",
			Index:    1,
			//Medium:   "rotate",
			Size: dataDisk * 1024,
		}
		if diskFormat == "qcow2" {
			r.Disks = []diskData{sysDiskData, dataDiskData}
		} else if diskFormat == "iso" {
			r.Disks = []diskData{sysDiskData, dataDiskData}
			r.Cdrom = imageId
		} else {
			r.Disks = []diskData{sysDiskData, dataDiskData}
		}
	} else {
		if diskFormat == "qcow2" {
			r.Disks = []diskData{sysDiskData}
		} else if diskFormat == "iso" {
			r.Disks = []diskData{sysDiskData}
			r.Cdrom = imageId
		} else {
			r.Disks = []diskData{sysDiskData}
		}
	}

	networkService := &NetworkService{}
	networkId, err := networkService.GetGuestOneNetwork()
	if err != nil {
		return nil, err
	}
	r.Nets = []netData{
		{
			Network: networkId,
		},
	}

	//todo 目前写死只分配英伟达的gpu
	isolatedDevices := make([]device, 0)
	for i := 0; i < gpu; i++ {
		var gpuDevice device
		gpuDevice.Vendor = "NVIDIA"
		gpuDevice.Model = "GeForce RTX 4090"
		gpuDevice.Dev_type = "GPU-HPC"
		isolatedDevices = append(isolatedDevices, gpuDevice)
	}

	r.IsolatedDevices = isolatedDevices

	params := jsonutils.Marshal(r)
	vmRet, err := compute.Servers.Create(s, params)
	if err != nil {
		return nil, err
	}

	return vmRet, nil

	//_, err = scheduler.SchedManager.DoForecast(s, jsonutils.Marshal(params))
	//if err != nil {
	//	if err != nil {
	//		response.FailWithMessage("创建失败", c)
	//	}
	//}

	//params = jsonutils.Marshal(r)
	//result, err = compute.Servers.Create(s, params)
	//
	//if err != nil {
	//	response.FailWithMessage("创建失败", c)
	//}
}

func (v *VirtualMachineService) HostForecast(vmRet jsonutils.JSONObject) error {
	var errMsg string
	if vmRet == nil {
		return errors.New("参数为空")
	}

	s, err := global.GetCloudClientSession()
	if err != nil {
		return errors.New("获取CloudClientSession失败")
	}

	vmId, err := vmRet.GetString("id")
	if err != nil {
		errMsg = "获取云主机id失败"
	}

	status, err := vmRet.GetString("status")
	if err != nil {
		errMsg = "获取云主机状态失败"
	}
	if status != "running" {
		size := 3000
		for i := 0; i < size; i++ {
			vmRet, err = compute.Servers.Get(s, vmId, nil)
			if err != nil {
				errMsg = "获取云主机信息失败"
			}
			if vmRet == nil {
				continue
			}
			status, _ = vmRet.GetString("status")
			if status == "running" {
				break
			}
			if i < size-1 {
				time.Sleep(2 * time.Second)
			}
		}
	}

	if status == "running" {
		ip, err := vmRet.GetString("ips")
		if err != nil {
			errMsg = "获取云主机ip失败"
		}

		username, password, err := v.HostGetLoginInfo(vmId)
		if err != nil {
			errMsg = "获取云主机登录信息失败"
		}

		var resInfo renter.RenterRes
		err = global.GVA_DB.Model(&renter.RenterRes{}).Where("resource_id = ?", vmId).First(&resInfo).Error
		if err != nil {
			errMsg = "未找到租户资源"
		}

		resInfo.SshHost = ip
		resInfo.SshPort = 22
		resInfo.SshUser = username
		resInfo.SshPwd = password
		resInfo.Status = 2
		resInfo.StartTime = time.Now()
		resInfo.EndTime = time.Now().Add(time.Duration(resInfo.UseTime) * time.Hour)
		err = global.GVA_DB.Model(&renter.RenterRes{}).Where("resource_id = ?", vmId).Updates(&resInfo).Error
		if err != nil {
			errMsg = "更新租户资源失败"
		}
	} else {
		errMsg = "云主机未正常运行，请和运维人员联系"
	}

	if errMsg == "" {
		return nil
	}

	return errors.New(errMsg)
}

func (v *VirtualMachineService) HostForecastById(id string) error {
	var errMsg string
	if id == "" {
		return errors.New("参数为空")
	}

	s, err := global.GetCloudClientSession()
	if err != nil {
		return errors.New("获取CloudClientSession失败")
	}

	status := ""
	size := 3000
	var vmRet jsonutils.JSONObject
	for i := 0; i < size; i++ {
		vmRet, err = compute.Servers.Get(s, id, nil)
		if err != nil {
			errMsg = "获取云主机信息失败"
		}
		if vmRet == nil {
			continue
		}
		status, _ = vmRet.GetString("status")
		if status == "running" {
			break
		}
		if i < size-1 {
			time.Sleep(2 * time.Second)
		}
	}

	if status == "running" {
		ip, err := vmRet.GetString("ips")
		if err != nil {
			errMsg = "获取云主机ip失败"
		}

		username, password, err := v.HostGetLoginInfo(id)
		if err != nil {
			errMsg = "获取云主机登录信息失败"
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
		resInfo.StartTime = time.Now()
		resInfo.EndTime = time.Now().Add(time.Duration(resInfo.UseTime) * time.Hour)
		err = global.GVA_DB.Model(&renter.RenterRes{}).Where("resource_id = ?", id).Updates(&resInfo).Error
		if err != nil {
			errMsg = "更新租户资源失败"
		}
	} else {
		errMsg = "云主机未正常运行，请和运维人员联系"
	}

	if errMsg == "" {
		return nil
	}

	return errors.New(errMsg)
}

func (v *VirtualMachineService) HostBandwidthLimit(id string, bandwidth int) error {
	s, err := global.GetCloudClientSession()
	if err != nil {
		return errors.New("获取CloudClientSession失败")
	}

	vmRet, err := compute.Servers.Get(s, id, nil)
	if err != nil {
		return errors.New("获取云主机信息失败")
	}

	mac, err := vmRet.GetString("macs")
	if err != nil {
		return errors.New("获取云主机mac地址失败")
	}

	params := jsonutils.NewDict()
	if regutils.MatchMacAddr(mac) {
		params.Add(jsonutils.NewString(mac), "mac")
	}

	params.Add(jsonutils.NewInt(int64(bandwidth)), "bandwidth")
	_, err = compute.Servers.PerformAction(s, id, "change-bandwidth", params)
	if err != nil {
		return err
	}
	return nil
}
func (v *VirtualMachineService) HostDelete(id string) error {
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

func (v *VirtualMachineService) HostGet(id string) (jsonutils.JSONObject, error) {
	s, err := global.GetCloudClientSession()
	if err != nil {
		return nil, err
	}

	return compute.Servers.Get(s, id, nil)
}

func (v *VirtualMachineService) HostListByIds(id []string) ([]jsonutils.JSONObject, int, error) {
	s, err := global.GetCloudClientSession()
	if err != nil {
		return nil, 0, err
	}

	params := jsonutils.NewDict()
	params.Set("scope", jsonutils.NewString("system"))
	params.Set("show_fail_reason", jsonutils.NewBool(true))
	params.Set("details", jsonutils.NewBool(true))
	params.Set("with_meta", jsonutils.NewBool(true))
	params.Set("filter", jsonutils.NewString("hypervisor.notin(baremetal,container)"))
	params.Set("id", jsonutils.NewStringArray(id))

	result, err := compute.Servers.List(s, params)
	if err != nil {
		return nil, 0, err
	}
	return result.Data, result.Total, nil
}

func (v *VirtualMachineService) HostList() ([]jsonutils.JSONObject, int, error) {
	s, err := global.GetCloudClientSession()
	if err != nil {
		return nil, 0, err
	}

	params := jsonutils.NewDict()
	params.Set("scope", jsonutils.NewString("system"))
	params.Set("show_fail_reason", jsonutils.NewBool(true))
	params.Set("details", jsonutils.NewBool(true))
	params.Set("with_meta", jsonutils.NewBool(true))
	params.Set("filter", jsonutils.NewString("hypervisor.notin(baremetal,container)"))

	result, err := compute.Servers.List(s, params)
	if err != nil {
		return nil, 0, err
	}
	return result.Data, result.Total, nil
}
func (v *VirtualMachineService) HostGetLoginInfo(id string) (string, string, error) {
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

func (v *VirtualMachineService) HostListByRenterID(renterId *int, pageIndex *int, pageSize *int) ([]jsonutils.JSONObject, int64, error) {
	s, err := global.GetCloudClientSession()
	if err != nil {
		return nil, 0, err
	}

	params := jsonutils.NewDict()
	params.Set("scope", jsonutils.NewString("system"))
	params.Set("show_fail_reason", jsonutils.NewBool(true))
	params.Set("details", jsonutils.NewBool(true))
	params.Set("with_meta", jsonutils.NewBool(true))
	if pageIndex != nil && pageSize != nil {
		params.Set("limit", jsonutils.NewInt(int64(*pageSize)))
		params.Set("offset", jsonutils.NewInt(int64(*pageSize*(*pageIndex-1))))
	}

	filters := make([]string, 0)
	if renterId != nil {
		filters = append(filters, "name.contains('"+strconv.Itoa(*renterId)+"')")
	}
	filters = append(filters, "hypervisor.notin(baremetal,container)")

	params.Set("filter", jsonutils.NewStringArray(filters))
	result, err := compute.Servers.List(s, params)
	if err != nil {
		return nil, 0, err
	}
	return result.Data, int64(result.Total), nil
}
