package utils

import (
	"ygpt/server/global"
	"yunion.io/x/jsonutils"
	"yunion.io/x/onecloud/pkg/mcclient/modules/k8s"
)

func GetNamespaceId(clusterID string) string {
	s, err := global.GetCloudClientSession()
	if err != nil {
		return ""
	}
	// 只检索ygpt命名空间下的deployment
	params := jsonutils.NewDict()
	params.Set("cluster", jsonutils.NewString(clusterID))

	namespaceResult, err := k8s.Namespaces.GetByName(s, "ygpt", params)
	if err != nil {
		params.Set("name", jsonutils.NewString("ygpt"))
		namespaceResult, err = k8s.Namespaces.Create(s, params)
		if err != nil {
			return ""
		}
	}

	namespaceId, err := namespaceResult.GetString("id")
	if err != nil {
		return ""
	}

	return namespaceId
}
