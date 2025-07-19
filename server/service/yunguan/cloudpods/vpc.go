package cloudpods

import (
	"ygpt/server/global"
	"yunion.io/x/jsonutils"
	"yunion.io/x/onecloud/pkg/mcclient/modules/compute"
)

type VpcService struct {
}

func (d *VpcService) GetVpcList() ([]jsonutils.JSONObject, int, error) {
	s, err := global.GetCloudClientSession()
	if err != nil {
		return nil, 0, err
	}
	params := jsonutils.NewDict()
	params.Set("scope", jsonutils.NewString("system"))
	params.Set("show_fail_reason", jsonutils.NewBool(true))
	params.Set("details", jsonutils.NewBool(true))
	result, err := compute.Vpcs.List(s, params)
	if err != nil {
		return nil, 0, err
	}
	return result.Data, result.Total, nil
}
