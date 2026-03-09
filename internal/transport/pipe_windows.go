//go:build windows

package transport

import (
	"fmt"
	"net"

	"github.com/Microsoft/go-winio"
)

func newPipeListener(address string) (net.Listener, error) {
	if err := ValidatePipePath(address); err != nil {
		return nil, fmt.Errorf("invalid pipe path: %w", err)
	}
	return winio.ListenPipe(address, nil)
}
