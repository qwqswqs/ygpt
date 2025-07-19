package request

import (
	"ygpt/server/model/common/request"
	"ygpt/server/model/system"
)

type SysOperationRecordSearch struct {
	system.SysOperationRecord
	request.PageInfo
}
