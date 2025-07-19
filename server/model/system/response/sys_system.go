package response

import "ygpt/server/config"

type SysConfigResponse struct {
	Config config.Server `json:"config"`
}
