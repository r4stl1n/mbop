package commands

import (
	"github.com/r4stl1n/mbop/internal/mbop/dto"
	"github.com/r4stl1n/mbop/pkg/api/llm"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

type ChatCompletion struct {
	ctx *dto.Context

	model string
	query string
}

func (c *ChatCompletion) Cmd(ctx *dto.Context) *cobra.Command {
	c.ctx = ctx

	cmd := &cobra.Command{
		Use:     "chatCompletion",
		Aliases: []string{"cc"},
		Args:    cobra.ExactArgs(0),
		Short:   "Attempt a chat completion",
		Run:     c.Init(ctx).Run,
	}

	cmd.Flags().StringVarP(&c.model, "model", "m", "gpt-3.5-turbo", "model to use for generation")
	cmd.Flags().StringVarP(&c.query, "query", "q", "", "query to use for generation")

	_ = cmd.MarkFlagRequired("query")

	return cmd
}

func (c *ChatCompletion) Init(ctx *dto.Context) *ChatCompletion {
	c.ctx = ctx
	return c
}

func (c *ChatCompletion) Run(_ *cobra.Command, _ []string) {

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

	zap.L().Info("attempting to get completion", zap.String("model", c.model), zap.String("query", c.query))

	chatCompletion, _, completionError := openaiClient.GetCompletion(llm.Completion{
		Model:        c.model,
		SystemPrompt: "You are a helpful assistant.",
		UserPrompt:   c.query,
	})

	if completionError != nil {
		zap.L().Fatal("failed to retrieve chat completion", zap.Error(completionError))
	}

	zap.L().Info("chat completion", zap.String("query", c.query), zap.String("message", chatCompletion))
}
