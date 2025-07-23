package exec

import (
	"context"

	"github.com/bestruirui/bestsub/internal/utils/log"
)

type Instance interface {
	Init() error
	Run(ctx context.Context, log *log.Logger) error
}
