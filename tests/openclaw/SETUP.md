# Scenario

**Feature**: `my openclaw` CLI manages OpenClaw data-dir registry and launches gateway in Podman.

```
# isolated config dir + stub podman + built my binary
MY_CONFIG_DIR -> registry JSON; stub podman -> recorded invocations; my openclaw <subcommand>
```

## Preconditions

- Module root is `filepath.Join(DOCTEST_ROOT, "..", "..")`.
- Every test uses an isolated `MY_CONFIG_DIR` (never the real user config).
- A stub `podman` script is prepended to `PATH` and logs invocations to `PODMAN_STUB_LOG`.
- The `my` CLI is built to `binDir/my` before each test run.

## Steps

1. Create `binDir`, `configDir`, and `podmanLog` under `t.TempDir()`.
2. Install the podman stub script into `binDir/podman`.
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
	req.BinPath = filepath.Join(req.BinDir, "my")
	req.PodmanMachineRunning = true

	if err := os.MkdirAll(req.BinDir, 0o755); err != nil {
		return err
	}
	if err := os.MkdirAll(req.ConfigDir, 0o755); err != nil {
		return err
	}
	if err := installPodmanStub(req.BinDir); err != nil {
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

func installPodmanStub(binDir string) error {
	script := `#!/bin/sh
log="${PODMAN_STUB_LOG:-/dev/stderr}"
printf '%s\n' "podman $*" >> "$log"
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
build|stop|rm|run)
  exit 0
  ;;
esac
exit 0
`
	path := filepath.Join(binDir, "podman")
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

func hashFileSHA256(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	sum := sha256.Sum256(data)
	return "sha256:" + hex.EncodeToString(sum[:]), nil
}

```