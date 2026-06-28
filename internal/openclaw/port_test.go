package openclaw

import (
	"net"
	"strconv"
	"testing"
)

func TestIsPortAvailableStubBusyPorts(t *testing.T) {
	t.Setenv("OPENCLAW_STUB_BUSY_PORTS", "18789,18790")
	if isPortAvailable(18789) {
		t.Fatal("expected stub busy port 18789 to be unavailable")
	}
	if isPortAvailable(18790) {
		t.Fatal("expected stub busy port 18790 to be unavailable")
	}
}

func TestIsPortAvailableStubLog(t *testing.T) {
	t.Setenv("OPENCLAW_STUB_LOG", "1")
	t.Setenv("OPENCLAW_STUB_BUSY_PORTS", "")
	if !isPortAvailable(18789) {
		t.Fatal("expected OPENCLAW_STUB_LOG to make ports available")
	}
}

func TestIsPortInUseDetectsListener(t *testing.T) {
	t.Setenv("OPENCLAW_STUB_BUSY_PORTS", "")
	t.Setenv("OPENCLAW_STUB_LOG", "")

	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	defer ln.Close()

	_, portStr, err := net.SplitHostPort(ln.Addr().String())
	if err != nil {
		t.Fatal(err)
	}
	port, err := strconv.Atoi(portStr)
	if err != nil {
		t.Fatal(err)
	}

	if isPortAvailable(port) {
		t.Fatalf("expected port %d with active listener to be unavailable", port)
	}
}