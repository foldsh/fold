package runtime

import (
	"fmt"

	uuid "github.com/satori/go.uuid"
)

// A helper function to create a random address for a unix domain socket.
func newAddr() string {
	sockId := uuid.NewV4().String()
	return fmt.Sprintf("/tmp/fold.%s.sock", sockId)
}
