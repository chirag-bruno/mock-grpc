//go:build !windows

package transport

import (
	"fmt"
	"net"
)

func newPipeListener(_ string) (net.Listener, error) {
	return nil, fmt.Errorf("pipe mode is only supported on Windows, use 'unix' mode instead")
}
