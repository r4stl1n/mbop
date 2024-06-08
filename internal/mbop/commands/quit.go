package commands

import (
	"github.com/r4stl1n/mbop/internal/mbop/dto"
	"github.com/spf13/cobra"
	"os"
)

type Quit struct {
	ctx *dto.Context
}

func (c *Quit) Cmd(ctx *dto.Context) *cobra.Command {
	c.ctx = ctx

	cmd := &cobra.Command{
		Use:     "quit",
		Aliases: []string{"exit", "bye", "x", "q"},
		Args:    cobra.ExactArgs(0),
		Short:   "A command the quit the utility when running in interactive mode",
		Run:     c.Init(ctx).Run,
	}

	return cmd
}

func (c *Quit) Init(ctx *dto.Context) *Quit {
	c.ctx = ctx
	return c
}

func (c *Quit) Run(_ *cobra.Command, _ []string) {
	os.Exit(0)
}
