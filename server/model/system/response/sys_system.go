package response

import "kubespace/server/config"

type SysConfigResponse struct {
	Config config.Server `json:"config"`
}
