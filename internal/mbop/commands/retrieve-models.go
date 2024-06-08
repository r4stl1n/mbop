package commands

import (
	"github.com/r4stl1n/mbop/internal/mbop/dto"
	"github.com/r4stl1n/mbop/pkg/api/llm"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

type RetrieveModels struct {
	ctx *dto.Context
}

func (c *RetrieveModels) Cmd(ctx *dto.Context) *cobra.Command {
	c.ctx = ctx

	cmd := &cobra.Command{
		Use:     "retrieveModels",
		Aliases: []string{"rm"},
		Args:    cobra.ExactArgs(0),
		Short:   "Retrieve models from OpenAI",
		Run:     c.Init(ctx).Run,
	}

	return cmd
}

func (c *RetrieveModels) Init(ctx *dto.Context) *RetrieveModels {
	c.ctx = ctx
	return c
}

func (c *RetrieveModels) Run(_ *cobra.Command, _ []string) {

	openaiClient, openaiClientError := new(llm.OpenAIAPI).Init()

	if openaiClientError != nil {
		zap.L().Fatal("failed to connect to openai api", zap.Error(openaiClientError))
	}

	connectionError := openaiClient.TestConnection()

	if connectionError != nil {
		zap.L().Fatal("failed to connect to openai api", zap.Error(connectionError))
	} else {
		zap.L().Info("connected to openai api")
	}

	zap.L().Info("attempting to get models from openai")

	models, modelsError := openaiClient.GetModels()

	if modelsError != nil {
		zap.L().Fatal("failed to retrieve models", zap.Error(modelsError))
	}

	for _, x := range models.Data {
		zap.L().Info("", zap.String("name", x.ID),
			zap.String("object", x.Object), zap.String("owned_by", x.OwnedBy))
	}

}
