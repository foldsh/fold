package gateway

import "fmt"

type Gateway struct {
	Port int
}

func (gw *Gateway) ImageName() string {
	return fmt.Sprintf("foldsh/foldgw:%s", GatewayVersion.String())
}
