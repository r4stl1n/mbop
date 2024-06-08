package client

import (
	"github.com/r4stl1n/mbop/internal/mbop/commands"
	"github.com/r4stl1n/mbop/internal/mbop/dto"
	"github.com/spf13/cobra"
)

type CLI struct {
	RootCmd *cobra.Command
	Context *dto.Context
	Error   error
}

func (c *CLI) Init() *CLI {
	*c = *new(CLI)

	c.RootCmd = &cobra.Command{
		Use:           "mbop",
		Short:         "Merry Band of Pirates",
		SilenceErrors: false,
		SilenceUsage:  false,
	}

	c.Context = new(dto.Context).Init(c.RootCmd.Version)

	newCommands := []*cobra.Command{
		new(commands.ChatCompletion).Cmd(c.Context),
		new(commands.RetrieveModels).Cmd(c.Context),
		new(commands.Quit).Cmd(c.Context),
		new(commands.Version).Cmd(c.Context),
	}

	cobra.OnInitialize(func() {

		c.RootCmd.ResetFlags()

	})

	c.RootCmd.AddCommand(newCommands...)

	return c
}

func (c *CLI) Run() error {
	if c.Error != nil {
		return c.Error
	}

	c.Error = c.RootCmd.Execute()
	if c.Error != nil {
		return c.Error
	}

	return c.Error
}
