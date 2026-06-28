package openclaw

import (
	"fmt"
	"net"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var portAlreadyInUsePattern = regexp.MustCompile(`(?i)port .* is already in use`)

const (
	defaultGatewayPort = 18789
	maxPortScanTries   = 100
)

func parsePortSetEnv(name string) map[int]bool {
	ports := make(map[int]bool)
	raw := strings.TrimSpace(os.Getenv(name))
	if raw == "" {
		return ports
	}
	for _, part := range strings.Split(raw, ",") {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}
		port, err := strconv.Atoi(part)
		if err != nil {
			continue
		}
		ports[port] = true
	}
	return ports
}

func stubBusyPorts() map[int]bool {
	return parsePortSetEnv("OPENCLAW_STUB_BUSY_PORTS")
}

func isPortInUse(port int) bool {
	if stubBusyPorts()[port] {
		return true
	}
	if os.Getenv("OPENCLAW_STUB_LOG") != "" {
		return false
	}

	addr := fmt.Sprintf("127.0.0.1:%d", port)
	conn, err := net.DialTimeout("tcp", addr, 500*time.Millisecond)
	if err == nil {
		_ = conn.Close()
		return true
	}

	ln, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", port))
	if err != nil {
		return true
	}
	_ = ln.Close()
	return false
}

func isPortAvailable(port int) bool {
	return !isPortInUse(port)
}

func selectRestartPort(explicit string) (port int, bumped bool, err error) {
	if explicit != "" {
		return selectPort(explicit)
	}
	return defaultGatewayPort, false, nil
}

func selectPort(explicit string) (port int, bumped bool, err error) {
	if explicit != "" {
		port, err = strconv.Atoi(explicit)
		if err != nil {
			return 0, false, fmt.Errorf("invalid port: %s", explicit)
		}
		if !isPortAvailable(port) {
			return 0, false, fmt.Errorf("port %d is in use", port)
		}
		return port, false, nil
	}

	for i := 0; i < maxPortScanTries; i++ {
		candidate := defaultGatewayPort + i
		if isPortAvailable(candidate) {
			return candidate, candidate != defaultGatewayPort, nil
		}
	}
	return 0, false, fmt.Errorf("no available port found in range %d-%d", defaultGatewayPort, defaultGatewayPort+maxPortScanTries-1)
}

func isPortInUseError(output string) bool {
	lower := strings.ToLower(output)
	if strings.Contains(lower, "eaddrinuse") {
		return true
	}
	if strings.Contains(lower, "already in use") {
		return true
	}
	return portAlreadyInUsePattern.MatchString(output)
}