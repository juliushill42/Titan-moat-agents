package connectors

import (
	"context"

	"titanuniversalai.com/shared/llm"
)

var llmRouter = llm.NewRouterFromEnv()

func Infer(ctx context.Context, prompt string) (string, error) {
	res, err := llmRouter.Complete(ctx, llm.CompletionRequest{Prompt: prompt})
	if err != nil {
		return "", err
	}
	return res.Text, nil
}
