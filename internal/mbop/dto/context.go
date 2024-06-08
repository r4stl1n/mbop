package dto

import "github.com/r4stl1n/mbop/pkg/util"

type Context struct {
	Version       string
	ConfigFile    string
	Utils         *util.Utils
	IsInteractive bool
}

func (c *Context) Init(version string) *Context {

	*c = Context{
		Version:       version,
		IsInteractive: false,
		Utils:         new(util.Utils).Init(),
	}

	return c
}
