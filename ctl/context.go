package ctl

import (
	"context"
	"io"

	"github.com/foldsh/fold/ctl/config"
	"github.com/foldsh/fold/ctl/output"
	"github.com/foldsh/fold/logging"
)

type CmdCtx struct {
	context.Context
	logging.Logger
	*config.Config
	*output.Multiplexer
}

func NewCmdCtx(
	ctx context.Context,
	logger logging.Logger,
	cfg *config.Config,
	out io.Writer,
) *CmdCtx {
	return &CmdCtx{
		Context:     ctx,
		Config:      cfg,
		Logger:      logger,
		Multiplexer: output.NewMultiplexer(out),
	}
}
