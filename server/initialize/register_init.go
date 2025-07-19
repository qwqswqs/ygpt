package initialize

import (
	_ "ygpt/server/source/example"
	_ "ygpt/server/source/system"
)

func init() {
	// do nothing,only import source package so that inits can be registered
}
