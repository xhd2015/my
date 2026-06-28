# Scenario

**Feature**: `my openclaw` CLI manages OpenClaw data-dir registry and launches gateway in Podman or on the host.

```
# isolated config dir + stub podman/openclaw + built my binary
MY_CONFIG_DIR -> registry JSON; stubs -> recorded invocations; my openclaw <subcommand>
```

## Preconditions

- Module root is `filepath.Join(DOCTEST_ROOT, "..", "..")`.
- Every test uses an isolated `MY_CONFIG_DIR` (never the real user config).
- A stub `podman` script is prepended to `PATH` and logs invocations to `PODMAN_STUB_LOG`.
- A stub `openclaw` script is prepended to `PATH` and logs invocations to `OPENCLAW_STUB_LOG`.
- A stub `npm` script is prepended to `PATH` and logs invocations to `NPM_STUB_LOG`; on
  `install -g openclaw@latest` it writes the openclaw stub into `NPM_STUB_BIN_DIR`.
- The `my` CLI is built to `binDir/my` before each test run.

## Steps

1. Create `binDir`, `configDir`, `podmanLog`, `openclawLog`, and `npmLog` under `t.TempDir()`.
2. Install the podman, openclaw, and npm stub scripts into `binDir/`.
3. Build `my` from `./cmd/my/`.
4. Seed registry JSON when `req.RegistrySeed` is set.
5. Leaf `Setup` customizes `Request` fields.
6. Invoke shared `Run` from `DOCTEST.md`.

## Context

- Registry path: `$MY_CONFIG_DIR/openclaw.json`.
- Data directories in tests are real temp dirs with fixture files copied from leaf `testdata/`.
- Podman machine checks apply on darwin when `MY_OPENCLAW_CHECK_PODMAN_MACHINE=1`.

```go
import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func Setup(t *testing.T, req *Request) error {
	req.CommandDir = filepath.Clean(filepath.Join(DOCTEST_ROOT, "..", ".."))
	req.BinDir = filepath.Join(t.TempDir(), "bin")
	req.ConfigDir = filepath.Join(t.TempDir(), "config")
	req.PodmanLog = filepath.Join(t.TempDir(), "podman.log")
	req.OpenClawLog = filepath.Join(t.TempDir(), "openclaw.log")
	req.OpenClawGatewayStateFile = filepath.Join(t.TempDir(), "openclaw-gateway.state")
	req.NpmLog = filepath.Join(t.TempDir(), "npm.log")
	req.BinPath = filepath.Join(req.BinDir, "my")
	req.PodmanMachineRunning = true
	req.PodmanStateFile = filepath.Join(t.TempDir(), "podman.state")

	if err := os.MkdirAll(req.BinDir, 0o755); err != nil {
		return err
	}
	if err := os.MkdirAll(req.ConfigDir, 0o755); err != nil {
		return err
	}
	if err := installPodmanStub(req.BinDir, req.PodmanStateFile); err != nil {
		return err
	}
	if err := installOpenclawStub(req.BinDir); err != nil {
		return err
	}
	if err := installNpmStub(req.BinDir); err != nil {
		return err
	}

	build := exec.Command("go", "build", "-o", req.BinPath, "./cmd/my/")
	build.Dir = req.CommandDir
	if output, err := build.CombinedOutput(); err != nil {
		return fmt.Errorf("go build my: %w: %s", err, strings.TrimSpace(string(output)))
	}

	if len(req.RegistrySeed) > 0 {
		if err := os.WriteFile(filepath.Join(req.ConfigDir, "openclaw.json"), req.RegistrySeed, 0o644); err != nil {
			return err
		}
	}

	t.Setenv("MY_CONFIG_DIR", req.ConfigDir)
	return nil
}

func installPodmanStub(binDir, stateFile string) error {
	script := fmt.Sprintf(`#!/bin/sh
log="${PODMAN_STUB_LOG:-/dev/stderr}"
state="%s"
printf '%%s\n' "podman $*" >> "$log"

stub_name=""
prev=""
for arg in "$@"; do
  if [ "$prev" = "--name" ]; then
    stub_name="$arg"
  fi
  prev="$arg"
done

case "$1" in
machine)
  case "$2" in
  info)
    if [ "${PODMAN_MACHINE_RUNNING:-1}" = "1" ]; then
      echo "Running: true"
    else
      echo "Running: false"
    fi
    exit 0
    ;;
  start)
    exit 0
    ;;
  esac
  ;;
images)
  if [ "${PODMAN_IMAGE_EXISTS:-0}" = "1" ]; then
    echo "my-openclaw:local"
  fi
  exit 0
  ;;
ps)
  if [ -f "$state" ]; then
    while IFS= read -r line; do
      [ -n "$line" ] && echo "$line"
    done < "$state"
  fi
  exit 0
  ;;
logs)
  target="${2:-${stub_name:-container}}"
  echo "stub log line for $target"
  exit 0
  ;;
exec)
  for arg in "$@"; do
    case "$arg" in
    *"command -v curl"*)
      if [ "${PODMAN_STUB_HAS_CURL:-0}" = "1" ]; then
        echo "/usr/bin/curl"
      fi
      exit 0
      ;;
    esac
  done
  if [ "$3" = "sh" ]; then
    echo "container-gateway-token"
    exit 0
  fi
  exit 0
  ;;
inspect)
  format=""
  prev=""
  for arg in "$@"; do
    if [ "$prev" = "--format" ]; then
      format="$arg"
    fi
    prev="$arg"
  done
  case "$format" in
  *HostPort*)
    echo "${PODMAN_STUB_HOST_PORT:-18789}"
    ;;
  *)
    if [ -n "${PODMAN_STUB_DATA_DIR:-}" ]; then
      echo "$PODMAN_STUB_DATA_DIR"
    fi
    ;;
  esac
  exit 0
  ;;
cp)
  exit 0
  ;;
build|stop|rm|run)
  if [ "$1" = "run" ] && [ -n "$stub_name" ]; then
    echo "$stub_name" >> "$state"
  fi
  if [ "$1" = "stop" ] || [ "$1" = "rm" ]; then
    target="${2:-}"
    if [ -n "$target" ] && [ -f "$state" ]; then
      tmp="$(mktemp)"
      while IFS= read -r line; do
        [ "$line" = "$target" ] && continue
        [ -n "$line" ] && echo "$line" >> "$tmp"
      done < "$state"
      mv "$tmp" "$state"
    fi
  fi
  exit 0
  ;;
esac
exit 0
`, stateFile)
	path := filepath.Join(binDir, "podman")
	if err := os.WriteFile(path, []byte(script), 0o755); err != nil {
		return err
	}
	return nil
}

func openclawStubScript() string {
	return `#!/bin/sh
log="${OPENCLAW_STUB_LOG:-/dev/stderr}"
state="${OPENCLAW_STUB_GATEWAY_STATE:-}"
printf 'openclaw %s\n' "$*" >> "$log"
printf 'OPENCLAW_STATE_DIR=%s\n' "${OPENCLAW_STATE_DIR:-}" >> "$log"
printf 'OPENCLAW_GATEWAY_TOKEN=%s\n' "${OPENCLAW_GATEWAY_TOKEN:-}" >> "$log"

gateway_running() {
  if [ -n "${OPENCLAW_STUB_RUNNING_DATA_DIRS:-}" ]; then
    if [ -z "${OPENCLAW_STATE_DIR:-}" ]; then
      return 1
    fi
    oldifs="$IFS"
    IFS=,
    for d in $OPENCLAW_STUB_RUNNING_DATA_DIRS; do
      d=$(echo "$d" | tr -d ' ')
      if [ "$d" = "$OPENCLAW_STATE_DIR" ]; then
        IFS="$oldifs"
        return 0
      fi
    done
    IFS="$oldifs"
    return 1
  fi
  if [ -n "$state" ] && [ -f "$state" ]; then
    case "$(tr -d '[:space:]' < "$state")" in
    0) return 1 ;;
    1) return 0 ;;
    esac
  fi
  [ "${OPENCLAW_STUB_GATEWAY_RUNNING:-0}" = "1" ]
}

set_gateway_running() {
  if [ -n "$state" ]; then
    printf '%s\n' "$1" > "$state"
  fi
}

gateway_status_require_rpc() {
  for arg in "$@"; do
    if [ "$arg" = "--require-rpc" ]; then
      return 0
    fi
  done
  return 1
}

gateway_status_port_busy() {
  busy_ports="${OPENCLAW_STUB_BUSY_PORTS:-}"
  if [ -z "$busy_ports" ]; then
    return 1
  fi
  default_port="18789"
  oldifs="$IFS"
  IFS=,
  for p in $busy_ports; do
    p=$(echo "$p" | tr -d ' ')
    if [ "$p" = "$default_port" ]; then
      IFS="$oldifs"
      return 0
    fi
  done
  IFS="$oldifs"
  return 1
}

if [ "$1" = "gateway" ]; then
  case "$2" in
  status)
    if gateway_status_require_rpc "$@"; then
      if gateway_running; then
        exit 0
      fi
      echo "gateway not running (rpc check failed)" >&2
      exit 1
    fi
    if gateway_running; then
      exit 0
    fi
    if gateway_status_port_busy; then
      exit 0
    fi
    echo "gateway not running" >&2
    exit 1
    ;;
  stop)
    set_gateway_running 0
    exit 0
    ;;
  esac
fi

port=""
prev=""
for arg in "$@"; do
  if [ "$prev" = "--port" ]; then
    port="$arg"
  fi
  prev="$arg"
done

fail_ports="${OPENCLAW_STUB_FAIL_PORTS:-}"
if [ -n "$port" ] && [ -n "$fail_ports" ]; then
  oldifs="$IFS"
  IFS=,
  for p in $fail_ports; do
    p=$(echo "$p" | tr -d ' ')
    if [ "$p" = "$port" ]; then
      echo "Error: listen EADDRINUSE: address already in use 0.0.0.0:$port" >&2
      exit 1
    fi
  done
  IFS="$oldifs"
fi

if [ "$1" = "gateway" ]; then
  set_gateway_running 1
fi

exit 0
`
}

func installOpenclawStub(binDir string) error {
	path := filepath.Join(binDir, "openclaw")
	if err := os.WriteFile(path, []byte(openclawStubScript()), 0o755); err != nil {
		return err
	}
	return nil
}

func installNpmStub(binDir string) error {
	script := fmt.Sprintf(`#!/bin/sh
log="${NPM_STUB_LOG:-/dev/stderr}"
bindir="${NPM_STUB_BIN_DIR:-}"
printf 'npm %%s\n' "$*" >> "$log"

if [ "$1" = "install" ]; then
  has_g=0
  has_pkg=0
  for arg in "$@"; do
    [ "$arg" = "-g" ] && has_g=1
    [ "$arg" = "openclaw@latest" ] && has_pkg=1
  done
  if [ "$has_g" = "1" ] && [ "$has_pkg" = "1" ] && [ -n "$bindir" ]; then
    cat > "$bindir/openclaw" << 'OPENCLAW_STUB'
%s
OPENCLAW_STUB
    chmod +x "$bindir/openclaw"
  fi
fi
exit 0
`, openclawStubScript())
	path := filepath.Join(binDir, "npm")
	if err := os.WriteFile(path, []byte(script), 0o755); err != nil {
		return err
	}
	return nil
}

func copyDir(src, dst string) error {
	return filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		rel, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}
		target := filepath.Join(dst, rel)
		if info.IsDir() {
			return os.MkdirAll(target, 0o755)
		}
		in, err := os.Open(path)
		if err != nil {
			return err
		}
		defer in.Close()
		if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
			return err
		}
		out, err := os.Create(target)
		if err != nil {
			return err
		}
		defer out.Close()
		_, err = io.Copy(out, in)
		return err
	})
}

func copyFixtureDir(t *testing.T, src string) string {
	t.Helper()
	dst := t.TempDir()
	if err := copyDir(src, dst); err != nil {
		t.Fatalf("copy fixture %s: %v", src, err)
	}
	abs, err := filepath.Abs(dst)
	if err != nil {
		t.Fatalf("abs fixture path: %v", err)
	}
	return abs
}

func writeRegistrySeed(path string, body string) error {
	return os.WriteFile(path, []byte(body), 0o644)
}

func markContainerRunning(req *Request, name string) error {
	if name == "" {
		name = "openclaw-gateway"
	}
	return os.WriteFile(req.PodmanStateFile, []byte(name+"\n"), 0o644)
}

func hashFileSHA256(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	sum := sha256.Sum256(data)
	return "sha256:" + hex.EncodeToString(sum[:]), nil
}

```