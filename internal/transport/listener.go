package transport

import (
	"fmt"
	"net"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

type Mode string

const (
	ModeHTTP Mode = "http"
	ModeUnix Mode = "unix"
	ModePipe Mode = "pipe"
)

func ParseMode(s string) (Mode, error) {
	switch s {
	case "http":
		return ModeHTTP, nil
	case "unix":
		return ModeUnix, nil
	case "pipe":
		return ModePipe, nil
	default:
		return "", fmt.Errorf("unsupported mode: %s (supported: http, unix, pipe)", s)
	}
}

func NewListener(mode Mode, address string) (net.Listener, error) {
	switch mode {
	case ModeHTTP:
		return net.Listen("tcp", address)
	case ModeUnix:
		return newUnixListener(address)
	case ModePipe:
		return newPipeListener(address)
	default:
		return nil, fmt.Errorf("unsupported mode: %s", mode)
	}
}

func newUnixListener(address string) (net.Listener, error) {
	if runtime.GOOS == "windows" {
		return nil, fmt.Errorf("unix socket mode is not supported on Windows, use 'pipe' mode instead")
	}
	if err := validateUnixSocketPath(address); err != nil {
		return nil, fmt.Errorf("invalid unix socket path: %w", err)
	}
	os.Remove(address)
	return net.Listen("unix", address)
}

func validateUnixSocketPath(path string) error {
	if !filepath.IsAbs(path) {
		return fmt.Errorf("path must be absolute: %s", path)
	}
	dir := filepath.Dir(path)
	info, err := os.Stat(dir)
	if err != nil {
		return fmt.Errorf("parent directory does not exist: %s", dir)
	}
	if !info.IsDir() {
		return fmt.Errorf("parent path is not a directory: %s", dir)
	}
	return nil
}

func ValidatePipePath(path string) error {
	if !strings.HasPrefix(path, `\\.\pipe\`) {
		return fmt.Errorf("path must start with '\\\\.\\pipe\\': %s", path)
	}
	pipeName := strings.TrimPrefix(path, `\\.\pipe\`)
	if pipeName == "" {
		return fmt.Errorf("pipe name cannot be empty")
	}
	if strings.ContainsAny(pipeName, `\`) {
		return fmt.Errorf("pipe name cannot contain backslashes: %s", pipeName)
	}
	return nil
}
