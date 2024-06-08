package commands

import (
	"fmt"
	"github.com/r4stl1n/mbop/internal/mbop/dto"
	"github.com/r4stl1n/mbop/internal/mbop/version"
	"github.com/spf13/cobra"
)

type Version struct {
	ctx *dto.Context
}

func (c *Version) Cmd(ctx *dto.Context) *cobra.Command {
	c.ctx = ctx

	cmd := &cobra.Command{
		Use:     "version",
		Aliases: []string{"v"},
		Args:    cobra.ExactArgs(0),
		Short:   "Show current build and version information",
		Run:     c.Init(ctx).Run,
	}

	return cmd
}

func (c *Version) Init(ctx *dto.Context) *Version {
	c.ctx = ctx
	return c
}

func (c *Version) Run(_ *cobra.Command, _ []string) {
	fmt.Printf("---------------------------------\n")
	fmt.Printf("Version: %s\n", version.Version)
	fmt.Printf("OS Arch: %s\n", version.OsArch)
	fmt.Printf("Go Version: %s\n", version.GoVersion)
	fmt.Printf("Git Commit: %s\n", version.GitCommit)
	fmt.Printf("Build Date: %s\n", version.BuildDate)
	fmt.Printf("---------------------------------")
}
