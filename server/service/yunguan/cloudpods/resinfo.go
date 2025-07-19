package cloudpods

import (
	"ygpt/server/global"

	"yunion.io/x/onecloud/pkg/mcclient"

	"yunion.io/x/onecloud/pkg/mcclient/modules/k8s"

	"yunion.io/x/jsonutils"
	"yunion.io/x/onecloud/pkg/mcclient/modules/compute"
)

type ResInfoService struct {
}

type res struct {
	Total   int64 `json:"total"`
	Used    int64 `json:"used"`
	NotUsed int64 `json:"not_used"`
	UsedPct int64 `json:"used_percent"`
}

type vmRes struct {
	Total      int64 `json:"total"`            //云主机总数
	Ready      int64 `json:"ready"`            //云主机开机数
	NoReady    int64 `json:"no_ready"`         //云主机关机数
	NoReadyPct int64 `json:"no_ready_percent"` //云主机关机率
}

type Node struct {
	Total    int64 `json:"total"`         //k8s集群节点数量
	Ready    int64 `json:"ready"`         //k8s集群节点健康数量
	NoReady  int64 `json:"no_ready"`      //k8s集群节点不健康数量
	ReadyPct int64 `json:"ready_percent"` //k8s集群节点健康率
}

type k8sRes struct {
	ClusterNum int64 `json:"cluster_num"` //k8s集群数量
	Node       Node  `json:"node"`        //k8s集群节点信息
	Cpu        res   `json:"cpu"`
	Memory     res   `json:"memory"`
}

type ResInfo struct {
	HostNum      int64  `json:"host_num"`      //宿主机数量
	BareMetalNum int64  `json:"baremetal_num"` //物理机数量
	Cpu          res    `json:"cpu"`           //单位：核
	Memory       res    `json:"memory"`        //单位：GB
	Storage      res    `json:"storage"`       //单位：GB
	Gpu          res    `json:"gpu"`           //单位：卡
	IP           res    `json:"ip"`            //单位：个
	VPC          res    `json:"vpc"`           //单位：个
	VM           vmRes  `json:"vm"`            //单位：台
	K8s          k8sRes `json:"k8s"`           //k8s的资源信息
}

const (
	mbToGB = 1024
)

func (s *ResInfoService) GetResInfo() (interface{}, error) {
	c, err := global.GetCloudClientSession()
	if err != nil {
		return "", err
	}

	params := jsonutils.NewDict()
	params.Set("scope", jsonutils.NewString("system"))
	params.Set("refresh", jsonutils.NewBool(true))
	usageObj, err := compute.Usages.GetGeneralUsage(c, params)
	if err != nil {
		return "", err
	}

	cpuTotal, err := usageObj.Int("hosts.cpu")
	if err != nil {
		return "", err
	}
	cpuUsed, err := usageObj.Int("all.servers.cpu")
	if err != nil {
		return "", err
	}
	cpuInfo := calculateResInfo(cpuTotal, cpuUsed)

	memTotal, err := usageObj.Int("hosts.memory")
	if err != nil {
		return "", err
	}
	memUsed, err := usageObj.Int("all.servers.memory")
	if err != nil {
		return "", err
	}
	memInfo := calculateResInfo(memTotal/mbToGB, memUsed/mbToGB)

	storages, err := usageObj.Int("storages")
	if err != nil {
		return "", err
	}
	storageUsed, err := usageObj.Int("all.servers.disk")
	if err != nil {
		return "", err
	}
	storageInfo := calculateResInfo(storages/mbToGB, storageUsed/mbToGB)

	isolatedDevices, err := usageObj.Int("isolated_devices")
	if err != nil {
		return "", err
	}
	isolatedDevicesUsed, err := usageObj.Int("all.servers.isolated_devices")
	if err != nil {
		return "", err
	}
	gpuInfo := calculateResInfo(isolatedDevices, isolatedDevicesUsed)

	ipNum, err := usageObj.Int("all.ports")
	if err != nil {
		return "", err
	}
	ipUsed, err := usageObj.Int("all.nics")
	if err != nil {
		return "", err
	}
	ipInfo := calculateResInfo(ipNum, ipUsed)

	vpcNum, err := usageObj.Int("all.vpcs")
	if err != nil {
		return "", err
	}
	vpcUsed, err := getVPCUsedCount(c)
	if err != nil {
		return "", err
	}
	vpcInfo := calculateResInfo(vpcNum, vpcUsed)

	vmNum, err := usageObj.Int("all.servers")
	if err != nil {
		return "", err
	}
	vmNoReady, err := usageObj.Int("all.ready_servers") //cloudpods云主机关机数用的这个字段
	if err != nil {
		return "", err
	}
	vmInfo := calculateVMResInfo(vmNum, vmNoReady)

	hostNum, err := usageObj.Int("hosts")
	if err != nil {
		return "", err
	}
	bareMetalNum, err := usageObj.Int("baremetals")
	if err != nil {
		return "", err
	}

	k8sResInfo, err := getK8sResInfo(c)
	if err != nil {
		return "", err
	}

	resInfo := &ResInfo{
		HostNum:      hostNum,
		BareMetalNum: bareMetalNum,
		Cpu:          cpuInfo,
		Memory:       memInfo,
		Storage:      storageInfo,
		Gpu:          gpuInfo,
		IP:           ipInfo,
		VPC:          vpcInfo,
		VM:           vmInfo,
		K8s:          *k8sResInfo,
	}

	return resInfo, nil
}

func getVPCUsedCount(c *mcclient.ClientSession) (int64, error) {
	params := jsonutils.NewDict()
	params.Set("scope", jsonutils.NewString("system"))
	params.Set("show_fail_reason", jsonutils.NewBool(true))
	params.Set("details", jsonutils.NewBool(true))

	result, err := compute.Vpcs.List(c, params)
	if err != nil {
		return 0, err
	}
	vpcUsed := int64(0)

	for _, vpc := range result.Data {
		networkCount, err := vpc.Int("network_count")
		if err != nil {
			return 0, err
		}
		if networkCount > 0 {
			vpcUsed++
		}
	}
	return vpcUsed, nil
}

func getK8sResInfo(c *mcclient.ClientSession) (*k8sRes, error) {
	params := jsonutils.NewDict()
	params.Set("scope", jsonutils.NewString("system"))
	params.Set("refresh", jsonutils.NewBool(true))

	k8sUsage, err := k8s.Usages.GetUsage(c, params)
	if err != nil {
		return nil, err
	}

	allObj, err := k8sUsage.Get("all")
	if err != nil {
		return nil, err
	}

	cluster, err := allObj.Get("cluster")
	if err != nil {
		return nil, err
	}

	clusterNum, err := cluster.Int("count")
	if err != nil {
		return nil, err
	}

	node, err := cluster.Get("node")
	if err != nil {
		return nil, err
	}

	nodeTotal, err := node.Int("count")
	if err != nil {
		return nil, err
	}
	nodeReady, err := node.Int("ready_count")
	if err != nil {
		return nil, err
	}
	nodeNoReady, err := node.Int("not_ready_count")
	if err != nil {
		return nil, err
	}

	nodeReadyPct := int64(0)
	if nodeTotal != 0 {
		nodeReadyPct = int64(int((float64(nodeReady) / float64(nodeTotal)) * 100))
	}

	cpuInfo, err := node.Get("cpu")
	if err != nil {
		return nil, err
	}
	memoryInfo, err := node.Get("memory")
	if err != nil {
		return nil, err
	}

	cpuTotal, err := cpuInfo.Int("capacity")
	if err != nil {
		return nil, err
	}
	cpuRequest, err := cpuInfo.Int("request")
	if err != nil {
		return nil, err
	}

	memoryTotal, err := memoryInfo.Int("capacity")
	if err != nil {
		return nil, err
	}
	memoryRequest, err := memoryInfo.Int("request")
	if err != nil {
		return nil, err
	}

	cpu := calculateResInfo(cpuTotal/1000, cpuRequest/1000)
	memory := calculateResInfo(memoryTotal/(1000*1000*1000*1000), memoryRequest/(1000*1000*1000*1000))

	return &k8sRes{
		ClusterNum: clusterNum,
		Node: Node{
			Total:    nodeTotal,
			Ready:    nodeReady,
			NoReady:  nodeNoReady,
			ReadyPct: nodeReadyPct,
		},
		Cpu:    cpu,
		Memory: memory,
	}, nil
}

func calculateResInfo(total, used int64) res {
	notUsed := total - used
	usedPct := int64(0)
	if total != 0 {
		usedPct = int64((float64(used) / float64(total)) * 100)
	}
	return res{
		Total:   total,
		Used:    used,
		NotUsed: notUsed,
		UsedPct: usedPct,
	}
}

func calculateVMResInfo(total, noReady int64) vmRes {
	ready := total - noReady
	noReadyPct := int64(0)
	if total != 0 {
		noReadyPct = int64((float64(noReady) / float64(total)) * 100)
	}
	return vmRes{
		Total:      total,
		Ready:      ready,
		NoReady:    noReady,
		NoReadyPct: noReadyPct,
	}
}
