package request

import (
	"ygpt/server/model/common/request"
	"ygpt/server/model/system"
)

type SysDictionaryDetailSearch struct {
	system.SysDictionaryDetail
	request.PageInfo
}
