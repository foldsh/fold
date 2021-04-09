package commands

import (
	"context"
	"io"

	"github.com/foldsh/fold/ctl/output"
	"github.com/foldsh/fold/logging"
)

type CmdCtx struct {
	Context context.Context
	Logger  logging.Logger
	Out     *output.Multiplexer
}

func NewCmdCtx(ctx context.Context, logger logging.Logger, out io.Writer) *CmdCtx {
	return &CmdCtx{Context: ctx, Logger: logger, Out: output.NewMultiplexer(out)}
}
