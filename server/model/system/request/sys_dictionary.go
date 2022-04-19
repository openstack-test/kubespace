package request

import (
	"kubespace/server/model/common/request"
	"kubespace/server/model/system"
)

type SysDictionarySearch struct {
	system.SysDictionary
	request.PageInfo
}
