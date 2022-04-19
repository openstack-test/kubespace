package request

import (
	"kubespace/server/model/{{.Package}}"
	"kubespace/server/model/common/request"
)

type {{.StructName}}Search struct{
    {{.Package}}.{{.StructName}}
    request.PageInfo
}
