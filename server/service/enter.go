package service

import (
	"kubespace/server/service/system"
)

type ServiceGroup struct {
	SystemServiceGroup  system.ServiceGroup
}

var ServiceGroupApp = new(ServiceGroup)
