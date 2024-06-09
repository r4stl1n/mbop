package commands

import (
	"github.com/r4stl1n/mbop/internal/mbop/dto"
	"github.com/r4stl1n/mbop/pkg/managers"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

type Sail struct {
	ctx *dto.Context

	model   string
	task    string
	crewDir string
}

func (c *Sail) Cmd(ctx *dto.Context) *cobra.Command {
	c.ctx = ctx

	cmd := &cobra.Command{
		Use:     "sail",
		Aliases: []string{"s"},
		Args:    cobra.ExactArgs(0),
		Short:   "Engage crew and attempt to solve a task",
		Run:     c.Init(ctx).Run,
	}

	cmd.Flags().StringVarP(&c.model, "model", "m", "gpt-3.5-turbo", "model to use for generation")
	cmd.Flags().StringVarP(&c.task, "task", "t", "", "task to work towards")
	cmd.Flags().StringVarP(&c.crewDir, "crewDir", "c", "./crew", "directory of the crew definitions")

	_ = cmd.MarkFlagRequired("task")

	return cmd
}

func (c *Sail) Init(ctx *dto.Context) *Sail {
	c.ctx = ctx
	return c
}

func (c *Sail) Run(_ *cobra.Command, _ []string) {

	sailManager, sailManagerError := new(managers.SailManager).Init(c.model, c.task, c.crewDir)

	if sailManagerError != nil {
		zap.L().Fatal("failed to initialize sail manager", zap.Error(sailManagerError))
	}

	runError := sailManager.Run()

	if runError != nil {
		zap.L().Fatal("failed to run sail", zap.Error(runError))
	}
}
