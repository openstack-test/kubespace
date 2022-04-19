package request

import (
	"kubespace/server/model/common/request"
	"kubespace/server/model/system"
)

type SysDictionaryDetailSearch struct {
	system.SysDictionaryDetail
	request.PageInfo
}
