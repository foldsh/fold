package ctl

import (
	"context"

	"github.com/foldsh/fold/ctl/config"
	"github.com/foldsh/fold/ctl/output"
	"github.com/foldsh/fold/logging"
)

type CmdCtx struct {
	context.Context
	*output.Output

	Logger logging.Logger
	Config *config.Config
}

func NewCmdCtx(
	ctx context.Context,
	logger logging.Logger,
	cfg *config.Config,
	out *output.Output,
) *CmdCtx {
	return &CmdCtx{
		Context: ctx,
		Config:  cfg,
		Logger:  logger,
		Output:  out,
	}
}
