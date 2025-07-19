package router

import (
	cloud "ygpt/server/router/cloudpods"
	"ygpt/server/router/compute"
	"ygpt/server/router/example"
	"ygpt/server/router/system"
	"ygpt/server/router/yunguan"
)

var RouterGroupApp = new(RouterGroup)

type RouterGroup struct {
	System  system.RouterGroup
	Example example.RouterGroup
	YunGuan yunguan.RouterGroup
	Compute compute.RouterGroup
	Cloud   cloud.RouterGroup
}
