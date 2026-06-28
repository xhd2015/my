# `my openclaw` CLI — Doc-Style Test Tree

Tests for the `my` CLI OpenClaw subcommand group: data-dir registry (`add`, `list`),
Podman gateway launcher (`run-in-podman`), and local host gateway launcher (`run`).

## Version

0.0.4

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
- **OpenClaw stub** — fake `openclaw` on `PATH` records invocations (argv plus `OPENCLAW_STATE_DIR` and
  `OPENCLAW_GATEWAY_TOKEN`) to `OPENCLAW_STUB_LOG`; exits 0 unless `--port` matches
  `OPENCLAW_STUB_FAIL_PORTS` (comma-separated), then prints an `EADDRINUSE` message and exits 1.
  `gateway stop` clears `OPENCLAW_STUB_GATEWAY_STATE` when set. Stub `gateway status` remains for
  compatibility but is **not** used by `my openclaw run` for running detection (see gateway bookkeeping).
- **Gateway bookkeeping** — per data dir, `<data-dir>/.my/gateway.json` records `{pid, port, started_at, kind}`.
  `localGatewayRunning` treats a gateway as running when bookkeeping exists, `pid` is alive (`Signal(0)`),
  and `port` is in use (`OPENCLAW_STUB_BUSY_PORTS` in tests). Stale bookkeeping (dead pid or free port) is
  removed lazily on read. Written on detached `--restart` / foreground launch; removed on stop or exit.
  Tests simulate running via `writeGatewayBookkeeping(t, dataDir, os.Getpid(), port)` plus busy-port hook.
- **Gateway running hook (deprecated)** — `OPENCLAW_STUB_GATEWAY_RUNNING=1` and
  `OPENCLAW_STUB_RUNNING_DATA_DIRS` previously drove stub `gateway status` for running simulation; replaced
  by gateway bookkeeping for `my openclaw run` tests. `OPENCLAW_STUB_GATEWAY_STATE` still tracks stub
  stop/start for `--restart` flows.
- **Slack stub hook** — `MY_OPENCLAW_SLACK_STUB=1` (test-only) skips real Slack HTTP; send and channel
  lookup succeed without network.
- **Npm stub** — fake `npm` on `PATH` records invocations to `NPM_STUB_LOG`; on
  `install -g openclaw@latest` writes the openclaw stub into `NPM_STUB_BIN_DIR` and exits 0.
- **Local gateway launcher** — `my openclaw run` validates data dir, resolves token, selects port, execs
  `openclaw gateway --bind lan --port <port>` with `OPENCLAW_STATE_DIR` and `OPENCLAW_GATEWAY_TOKEN`.
- **Install prompt hook** — when stdin is a TTY (`MY_OPENCLAW_EXEC_INTERACTIVE=1`), missing `openclaw`
  triggers `Install openclaw now? [y/N]: `; canned answer via `MY_OPENCLAW_INSTALL_ANSWER=yes|no`.
- **Port occupancy hook** — `OPENCLAW_STUB_BUSY_PORTS=18789,18790` (test-only) marks ports unavailable for
  local port selection logic.
- **Runtime port fail hook** — `OPENCLAW_STUB_FAIL_PORTS=18789` makes the openclaw stub fail at launch with
  a port-in-use error so auto-port mode can retry the next port.

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
├── run/
│   ├── prerequisites/
│   │   ├── missing-data-dir
│   │   ├── missing-openclaw-json
│   │   ├── missing-token
│   │   ├── missing-openclaw-binary
│   │   ├── install-prompt-decline
│   │   └── install-prompt-accept
│   ├── launch/
│   │   ├── happy-path
│   │   ├── auto-port-bump
│   │   ├── runtime-port-retry
│   │   ├── port-in-use-explicit
│   │   ├── token-from-json
│   │   ├── token-from-env
│   │   ├── token-json-overrides-env
│   │   └── custom-port
│   ├── manage/
│   │   ├── status-not-running
│   │   ├── status-running
│   │   ├── restart-stopped
│   │   ├── restart-running
│   │   ├── import-local-grok
│   │   ├── test-slack-not-running
│   │   └── test-slack-running
│   ├── gateway-bookkeeping/
│   │   ├── stale-pid-cleanup
│   │   └── write-on-restart
│   └── resolve-data-dir/
│       ├── no-registry
│       ├── none-running
│       ├── multiple-running
│       ├── one-running-auto
│       ├── port-busy-not-running
│       └── explicit-override
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
| run missing data dir | `run/prerequisites/missing-data-dir` |
| run missing openclaw.json | `run/prerequisites/missing-openclaw-json` |
| run missing token | `run/prerequisites/missing-token` |
| run missing openclaw binary | `run/prerequisites/missing-openclaw-binary` |
| run install prompt decline | `run/prerequisites/install-prompt-decline` |
| run install prompt accept | `run/prerequisites/install-prompt-accept` |
| run happy path | `run/launch/happy-path` |
| run auto port bump | `run/launch/auto-port-bump` |
| run runtime port retry | `run/launch/runtime-port-retry` |
| run port in use explicit | `run/launch/port-in-use-explicit` |
| run token from json | `run/launch/token-from-json` |
| run token from env | `run/launch/token-from-env` |
| run token json overrides env | `run/launch/token-json-overrides-env` |
| run custom port | `run/launch/custom-port` |
| run status not running | `run/manage/status-not-running` |
| run status running | `run/manage/status-running` |
| run restart stopped | `run/manage/restart-stopped` |
| run restart running | `run/manage/restart-running` |
| run import local grok | `run/manage/import-local-grok` |
| run test slack not running | `run/manage/test-slack-not-running` |
| run test slack running | `run/manage/test-slack-running` |
| run resolve no registry | `run/resolve-data-dir/no-registry` |
| run resolve none running | `run/resolve-data-dir/none-running` |
| run resolve multiple running | `run/resolve-data-dir/multiple-running` |
| run resolve one running auto | `run/resolve-data-dir/one-running-auto` |
| run resolve port busy not running | `run/resolve-data-dir/port-busy-not-running` |
| run resolve explicit override | `run/resolve-data-dir/explicit-override` |
| run stale pid cleanup | `run/gateway-bookkeeping/stale-pid-cleanup` |
| run write on restart | `run/gateway-bookkeeping/write-on-restart` |

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
	PodmanStateFile string
	OpenClawLog string
	OpenClawGatewayStateFile string
	NpmLog      string

	Subcommand string

	DataDirPath string
	Note        string

	RunDataDir    string
	Rebuild       bool
	Stop          bool
	Restart       bool
	Logs          bool
	Status        bool
	ShowTokens    bool
	Dashboard     bool
	Exec            bool
	ExecArgs        []string
	ImportLocalGrok        bool
	TestSlack              bool
	SlackChannel           string
	InstallGrok            bool
	GrokAuthPath           string
	PodmanContainerDataDir string
	ContainerName          string
	Port          string
	BusyPorts     string
	FailPorts     string

	InstallPrompt bool
	InstallAnswer string
	BinDirOnlyPATH bool
	SlackStub              bool

	PodmanImageExists      bool
	PodmanMachineRunning   bool
	PodmanContainerRunning bool
	RegistrySeed           []byte
}

type Response struct {
	Stdout   string
	Stderr   string
	ExitCode int

	RegistryJSON  []byte
	PodmanCalls   []string
	OpenClawCalls []string
	NpmCalls      []string
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
	interactive := "0"
	if req.InstallPrompt {
		interactive = "1"
	}
	pathVal := req.BinDir + string(os.PathListSeparator) + os.Getenv("PATH")
	if req.BinDirOnlyPATH {
		pathVal = req.BinDir
	}
	env := append(os.Environ(),
		"MY_CONFIG_DIR="+req.ConfigDir,
		"MY_OPENCLAW_NO_OPEN=1",
		"MY_OPENCLAW_EXEC_INTERACTIVE="+interactive,
		"PODMAN_STUB_LOG="+req.PodmanLog,
		"OPENCLAW_STUB_LOG="+req.OpenClawLog,
		"NPM_STUB_LOG="+req.NpmLog,
		"NPM_STUB_BIN_DIR="+req.BinDir,
		"PATH="+pathVal,
	)
	if req.BusyPorts != "" {
		env = append(env, "OPENCLAW_STUB_BUSY_PORTS="+req.BusyPorts)
	}
	if req.FailPorts != "" {
		env = append(env, "OPENCLAW_STUB_FAIL_PORTS="+req.FailPorts)
	}
	if req.OpenClawGatewayStateFile != "" {
		env = append(env, "OPENCLAW_STUB_GATEWAY_STATE="+req.OpenClawGatewayStateFile)
	}
	if req.SlackStub {
		env = append(env, "MY_OPENCLAW_SLACK_STUB=1")
	}
	if req.InstallAnswer != "" {
		env = append(env, "MY_OPENCLAW_INSTALL_ANSWER="+req.InstallAnswer)
	}
	if req.GrokAuthPath != "" {
		env = append(env, "MY_GROK_AUTH_PATH="+req.GrokAuthPath)
	}
	if req.PodmanContainerDataDir != "" {
		env = append(env, "PODMAN_STUB_DATA_DIR="+req.PodmanContainerDataDir)
	}
	if v := os.Getenv("PODMAN_STUB_HAS_CURL"); v != "" {
		env = append(env, "PODMAN_STUB_HAS_CURL="+v)
	}
	return env
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

func readOpenclawLog(path string) ([]string, error) {
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

func readNpmLog(path string) ([]string, error) {
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

func npmCallsContain(calls []string, substr string) bool {
	for _, call := range calls {
		if strings.Contains(call, substr) {
			return true
		}
	}
	return false
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
		args := []string{"openclaw", "run-in-podman"}
		if req.Exec {
			args = append(args, "--exec")
			args = append(args, req.ExecArgs...)
			return args
		}
		if req.ImportLocalGrok {
			args = append(args, "--import-local-grok")
			if req.RunDataDir != "" {
				args = append(args, "--data-dir", req.RunDataDir)
			}
			return args
		}
		if req.InstallGrok {
			args = append(args, "--install-grok")
			if req.ContainerName != "" {
				args = append(args, "--container-name", req.ContainerName)
			}
			return args
		}
		if req.Stop {
			args = append(args, "--stop")
		} else if req.Restart {
			args = append(args, "--restart")
			if req.RunDataDir != "" {
				args = append(args, "--data-dir", req.RunDataDir)
			}
		} else if req.Logs {
			args = append(args, "--logs")
		} else if req.Status {
			args = append(args, "--status")
		} else if req.ShowTokens {
			args = append(args, "--show-tokens")
			if req.RunDataDir != "" {
				args = append(args, "--data-dir", req.RunDataDir)
			}
		} else if req.Dashboard {
			args = append(args, "--dashboard")
			if req.RunDataDir != "" {
				args = append(args, "--data-dir", req.RunDataDir)
			}
		} else {
			args = append(args, "--data-dir", req.RunDataDir)
		}
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
	case "run":
		args := []string{"openclaw", "run"}
		if req.ImportLocalGrok {
			args = append(args, "--import-local-grok")
			if req.RunDataDir != "" {
				args = append(args, "--data-dir", req.RunDataDir)
			}
			return args
		}
		if req.Restart {
			args = append(args, "--restart")
			if req.RunDataDir != "" {
				args = append(args, "--data-dir", req.RunDataDir)
			}
			if req.Port != "" {
				args = append(args, "--port", req.Port)
			}
			return args
		}
		if req.Status {
			args = append(args, "--status")
			if req.RunDataDir != "" {
				args = append(args, "--data-dir", req.RunDataDir)
			}
			if req.Port != "" {
				args = append(args, "--port", req.Port)
			}
			return args
		}
		if req.TestSlack {
			args = append(args, "--test-slack")
			if req.RunDataDir != "" {
				args = append(args, "--data-dir", req.RunDataDir)
			}
			if req.SlackChannel != "" {
				args = append(args, "--slack-channel", req.SlackChannel)
			}
			return args
		}
		args = append(args, "--data-dir", req.RunDataDir)
		if req.Port != "" {
			args = append(args, "--port", req.Port)
		}
		return args
	default:
		return nil
	}
}

func openclawCallsContain(calls []string, substr string) bool {
	for _, call := range calls {
		if strings.Contains(call, substr) {
			return true
		}
	}
	return false
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
	if req.PodmanLog != "" {
		if err := os.WriteFile(req.PodmanLog, nil, 0o644); err != nil && !os.IsNotExist(err) {
			return nil, err
		}
	}
	if req.OpenClawLog != "" {
		if err := os.WriteFile(req.OpenClawLog, nil, 0o644); err != nil && !os.IsNotExist(err) {
			return nil, err
		}
	}
	if req.NpmLog != "" {
		if err := os.WriteFile(req.NpmLog, nil, 0o644); err != nil && !os.IsNotExist(err) {
			return nil, err
		}
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

	openclawCalls, err := readOpenclawLog(req.OpenClawLog)
	if err != nil {
		return nil, fmt.Errorf("read openclaw log: %w", err)
	}

	npmCalls, err := readNpmLog(req.NpmLog)
	if err != nil {
		return nil, fmt.Errorf("read npm log: %w", err)
	}

	return &Response{
		Stdout:        stdout,
		Stderr:        stderr,
		ExitCode:      exitCode,
		RegistryJSON:  registryJSON,
		PodmanCalls:   podmanCalls,
		OpenClawCalls: openclawCalls,
		NpmCalls:      npmCalls,
	}, nil
}
```