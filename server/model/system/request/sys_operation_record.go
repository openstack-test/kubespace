package request

import (
	"kubespace/server/model/common/request"
	"kubespace/server/model/system"
)

type SysOperationRecordSearch struct {
	system.SysOperationRecord
	request.PageInfo
}
