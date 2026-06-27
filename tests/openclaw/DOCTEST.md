# `my openclaw` CLI — Doc-Style Test Tree

Tests for the `my` CLI OpenClaw subcommand group: data-dir registry (`add`, `list`)
and Podman gateway launcher (`run-in-podman`).

## Version

0.0.2

# DSN (Domain Specific Notion)

- **CLI binary** — `my` built from `cmd/my/` at the module root (`github.com/xhd2015/my`).
- **Config store** — registry JSON at `$MY_CONFIG_DIR/openclaw.json` (override `MY_CONFIG_DIR` in tests).
- **Registry** — explicit bookkeeping of `.openclaw` data directories (`path`, `note`, `added_at`) plus
  embedded container image metadata (`image.spec_hash`, `image.built_at`).
- **Data directory** — user's OpenClaw state dir (the `.openclaw` folder itself); must contain
  `openclaw.json`; may contain `.env` and `workspace/`.
- **Gateway token resolver** — reads `gateway.auth.token` from `openclaw.json`, else
  `OPENCLAW_GATEWAY_TOKEN` from `<data-dir>/.env`; errors when neither is present.
- **Podman stub** — fake `podman` on `PATH` records invocations; simulates machine state and image presence.
- **Container image** — `my-openclaw:local` built from embedded Containerfile; rebuild when image missing,
  `--rebuild`, or stored spec hash differs from current Containerfile SHA256.
- **Podman launcher** — ensures machine running (macOS), validates data dir, resolves token, rebuilds image
  if needed, stops/removes old container, runs gateway with volume mounts and token env, prints dashboard URL.

## How to Run

```sh
doctest vet ./tests/openclaw
doctest test -v ./tests/openclaw
```

## Decision Tree

```
subcommand
├── registry/
│   ├── add/
│   │   ├── path-exists/
│   │   │   ├── new-dir
│   │   │   ├── with-note
│   │   │   └── duplicate-warns
│   │   └── path-missing/
│   │       └── missing-dir
│   └── list/
│       ├── empty
│       └── populated
└── run-in-podman/
    ├── prerequisites/
    │   ├── missing-data-dir
    │   ├── missing-openclaw-json
    │   └── missing-token
    ├── image-rebuild/
    │   ├── rebuild-flag
    │   ├── spec-hash-stale
    │   └── skip-when-image-exists
    ├── launch/
    │   ├── happy-path
    │   ├── token-from-json
    │   ├── token-from-env
    │   ├── token-json-overrides-env
    │   ├── custom-container-name
    │   └── custom-port
    └── podman-machine/
        └── auto-start-when-stopped
```

## Test Leaf Index

| Leaf | Path |
|------|------|
| add new dir | `registry/add/path-exists/new-dir` |
| add with note | `registry/add/path-exists/with-note` |
| add duplicate warns | `registry/add/path-exists/duplicate-warns` |
| add missing dir | `registry/add/path-missing/missing-dir` |
| list empty | `registry/list/empty` |
| list populated | `registry/list/populated` |
| missing data dir | `run-in-podman/prerequisites/missing-data-dir` |
| missing openclaw.json | `run-in-podman/prerequisites/missing-openclaw-json` |
| missing token | `run-in-podman/prerequisites/missing-token` |
| happy path | `run-in-podman/launch/happy-path` |
| token from json | `run-in-podman/launch/token-from-json` |
| token from env | `run-in-podman/launch/token-from-env` |
| token json overrides env | `run-in-podman/launch/token-json-overrides-env` |
| custom container name | `run-in-podman/launch/custom-container-name` |
| custom port | `run-in-podman/launch/custom-port` |
| rebuild flag | `run-in-podman/image-rebuild/rebuild-flag` |
| spec hash stale | `run-in-podman/image-rebuild/spec-hash-stale` |
| skip rebuild | `run-in-podman/image-rebuild/skip-when-image-exists` |
| podman machine auto-start | `run-in-podman/podman-machine/auto-start-when-stopped` |

```go
import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

type Request struct {
	CommandDir string
	BinPath    string
	BinDir     string
	ConfigDir  string
	PodmanLog  string

	Subcommand string

	DataDirPath string
	Note        string

	RunDataDir    string
	Rebuild       bool
	ContainerName string
	Port          string

	PodmanImageExists    bool
	PodmanMachineRunning bool
	RegistrySeed         []byte
}

type Response struct {
	Stdout   string
	Stderr   string
	ExitCode int

	RegistryJSON []byte
	PodmanCalls  []string
}

type registryFile struct {
	DataDirs []struct {
		Path    string `json:"path"`
		Note    string `json:"note"`
		AddedAt string `json:"added_at"`
	} `json:"data_dirs"`
	Image *struct {
		SpecHash string `json:"spec_hash"`
		BuiltAt  string `json:"built_at"`
	} `json:"image,omitempty"`
}

func isolatedEnv(req *Request) []string {
	return append(os.Environ(),
		"MY_CONFIG_DIR="+req.ConfigDir,
		"PODMAN_STUB_LOG="+req.PodmanLog,
		"PATH="+req.BinDir+string(os.PathListSeparator)+os.Getenv("PATH"),
	)
}

func runMyCapture(env []string, dir string, bin string, args ...string) (string, string, int, error) {
	cmd := exec.Command(bin, args...)
	cmd.Dir = dir
	cmd.Env = env
	var stdout, stderr strings.Builder
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	exitCode := 0
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			exitCode = exitErr.ExitCode()
		} else {
			return "", "", 0, err
		}
	}
	return stdout.String(), stderr.String(), exitCode, nil
}

func readRegistry(path string) ([]byte, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	return data, nil
}

func readPodmanLog(path string) ([]string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	lines := strings.Split(strings.TrimSpace(string(data)), "\n")
	if len(lines) == 1 && lines[0] == "" {
		return nil, nil
	}
	return lines, nil
}

func buildArgs(req *Request) []string {
	switch req.Subcommand {
	case "add":
		args := []string{"openclaw", "add", "data-dir", req.DataDirPath}
		if req.Note != "" {
			args = append(args, "--note", req.Note)
		}
		return args
	case "list":
		return []string{"openclaw", "list"}
	case "run-in-podman":
		args := []string{"openclaw", "run-in-podman", "--data-dir", req.RunDataDir}
		if req.Rebuild {
			args = append(args, "--rebuild")
		}
		if req.ContainerName != "" {
			args = append(args, "--container-name", req.ContainerName)
		}
		if req.Port != "" {
			args = append(args, "--port", req.Port)
		}
		return args
	default:
		return nil
	}
}

func podmanCallsContain(calls []string, substr string) bool {
	for _, call := range calls {
		if strings.Contains(call, substr) {
			return true
		}
	}
	return false
}

func podmanCallsCount(calls []string, substr string) int {
	n := 0
	for _, call := range calls {
		if strings.Contains(call, substr) {
			n++
		}
	}
	return n
}

func parseRegistry(data []byte) (*registryFile, error) {
	if len(data) == 0 {
		return &registryFile{}, nil
	}
	var reg registryFile
	if err := json.Unmarshal(data, &reg); err != nil {
		return nil, err
	}
	return &reg, nil
}

func Run(t *testing.T, req *Request) (*Response, error) {
	if err := os.WriteFile(req.PodmanLog, nil, 0o644); err != nil && !os.IsNotExist(err) {
		return nil, err
	}

	env := isolatedEnv(req)
	if req.PodmanImageExists {
		env = append(env, "PODMAN_IMAGE_EXISTS=1")
	} else {
		env = append(env, "PODMAN_IMAGE_EXISTS=0")
	}
	if req.PodmanMachineRunning {
		env = append(env, "PODMAN_MACHINE_RUNNING=1")
	} else {
		env = append(env, "PODMAN_MACHINE_RUNNING=0")
	}
	if runtime.GOOS == "darwin" {
		env = append(env, "MY_OPENCLAW_CHECK_PODMAN_MACHINE=1")
	}

	stdout, stderr, exitCode, err := runMyCapture(env, req.CommandDir, req.BinPath, buildArgs(req)...)
	if err != nil {
		return nil, err
	}

	registryPath := filepath.Join(req.ConfigDir, "openclaw.json")
	registryJSON, err := readRegistry(registryPath)
	if err != nil {
		return nil, fmt.Errorf("read registry: %w", err)
	}

	podmanCalls, err := readPodmanLog(req.PodmanLog)
	if err != nil {
		return nil, fmt.Errorf("read podman log: %w", err)
	}

	return &Response{
		Stdout:       stdout,
		Stderr:       stderr,
		ExitCode:     exitCode,
		RegistryJSON: registryJSON,
		PodmanCalls:  podmanCalls,
	}, nil
}
```