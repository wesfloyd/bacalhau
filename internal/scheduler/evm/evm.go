package evm

import (
	"context"

	"github.com/filecoin-project/bacalhau/internal/types"
)

type EVMScheduler struct {
	Ctx context.Context

	Jobs map[string]*types.Job
}
