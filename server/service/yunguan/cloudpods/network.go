package cloudpods

import (
	"errors"
	"ygpt/server/global"
	"yunion.io/x/jsonutils"
	"yunion.io/x/onecloud/pkg/mcclient/modules/compute"
)

type NetworkService struct {
}

func (n *NetworkService) GetGuestOneNetwork() (string, error) {
	s, err := global.GetCloudClientSession()
	if err != nil {
		return "", err
	}
	params := jsonutils.NewDict()
	params.Set("scope", jsonutils.NewString("system"))
	params.Set("show_fail_reason", jsonutils.NewBool(true))
	params.Set("details", jsonutils.NewBool(true))

	networks, err := compute.Networks.List(s, params)
	if err != nil {
		return "", err
	}
	for _, network := range networks.Data {
		serverType, err := network.GetString("server_type")
		if err != nil {
			return "", err
		}
		if serverType != "guest" {
			continue
		}

		ipTotal, err := network.Int("ports")
		if err != nil {
			continue
		}

		ipUsed, err := network.Int("ports_used")
		if ipUsed >= ipTotal {
			continue
		}

		networkId, err := network.GetString("id")
		if err != nil {
			return "", err
		}
		return networkId, nil
	}

	return "", errors.New("one guest network not found")
}
