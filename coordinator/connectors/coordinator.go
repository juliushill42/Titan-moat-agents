package connectors

import (
	"context"
	"encoding/json"

	"titanuniversalai.com/shared/planner"
)

var Coordinator = planner.NewClient("coordinator")

func Start(ctx context.Context, stack string, handler func(context.Context, *planner.Task) (json.RawMessage, error)) error {
	return Coordinator.RunLoop(ctx, stack, handler)
}
