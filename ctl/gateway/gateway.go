package gateway

import (
	"fmt"

	"github.com/foldsh/fold/version"
)

type Gateway struct {
	Port int
}

func (gw *Gateway) ImageName() string {
	return fmt.Sprintf("foldsh/foldgw:%s", version.FoldVersion.String())
}
