package openclaw

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	lessflags "github.com/xhd2015/less-flags"
)

const runHelpText = `Usage: my openclaw run [action] [options]

Run or manage the OpenClaw gateway on the host.

Lifecycle:
  (default)    Start gateway in foreground
               [--data-dir <path>] [--port PORT]
  --restart    Stop (if running) and start gateway detached
               [--data-dir <path>] [--port PORT]
  --status     Show gateway URLs and info [--data-dir <path>] [--port PORT]

Settings:
  --import-local-grok
               Import ~/.grok OAuth, sync models, reset main session
               [--data-dir <path>]

Slack:
  --test-slack Send a test message (default: #general)
               [--data-dir <path>] [--slack-channel <name|id>]

Shared options:
  --data-dir <path>   Path to .openclaw directory (optional when exactly one
                      registered gateway is running)
  --port <port>       Host port (default: auto-select from 18789)
  -h, --help          Show this help
`

func lookupOpenclaw() (string, error) {
	return exec.LookPath("openclaw")
}

func readInstallAnswer() string {
	fmt.Fprint(os.Stderr, "Install openclaw now? [y/N]: ")
	if ans := strings.TrimSpace(os.Getenv("MY_OPENCLAW_INSTALL_ANSWER")); ans != "" {
		return strings.ToLower(ans)
	}
	line, _ := bufio.NewReader(os.Stdin).ReadString('\n')
	return strings.ToLower(strings.TrimSpace(line))
}

func isYesAnswer(answer string) bool {
	return answer == "y" || answer == "yes"
}

func pathWithSystemBins(pathVal string) string {
	for _, dir := range []string{"/usr/bin", "/bin", "/usr/local/bin"} {
		if dir == "" {
			continue
		}
		sep := string(os.PathListSeparator)
		if !strings.HasPrefix(pathVal, dir+sep) && !strings.Contains(pathVal, sep+dir+sep) && pathVal != dir {
			pathVal = dir + sep + pathVal
		}
	}
	return pathVal
}

func installOpenclawGlobal() error {
	printCommand("npm install -g openclaw@latest")
	cmd := exec.Command("npm", "install", "-g", "openclaw@latest")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = append(os.Environ(), "PATH="+pathWithSystemBins(os.Getenv("PATH")))
	return cmd.Run()
}

func ensureOpenclawAvailable() bool {
	if _, err := lookupOpenclaw(); err == nil {
		return true
	}

	fmt.Fprintln(os.Stderr, "error: openclaw not found on PATH")
	fmt.Fprintln(os.Stderr, "Install with: npm install -g openclaw@latest")
	fmt.Fprintln(os.Stderr, "Then run: openclaw onboard")

	if !stdinIsTerminal() {
		return false
	}

	if !isYesAnswer(readInstallAnswer()) {
		return false
	}

	if err := installOpenclawGlobal(); err != nil {
		return false
	}

	_, err := lookupOpenclaw()
	return err == nil
}

func validateLocalDataDir(absPath string) error {
	info, err := os.Stat(absPath)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("data directory not found: %s", absPath)
		}
		return err
	}
	if !info.IsDir() {
		return fmt.Errorf("path is not a directory: %s", absPath)
	}

	configPath := filepath.Join(absPath, "openclaw.json")
	if _, err := os.Stat(configPath); err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("openclaw.json not found in %s", absPath)
		}
		return err
	}
	return nil
}

func printLocalRunInfo(dataDir string, port int, bumped bool) {
	if bumped {
		fmt.Printf("Port %d in use; using %d instead\n", defaultGatewayPort, port)
	}
	fmt.Printf("Data dir: %s\n", dataDir)
	fmt.Printf("Port: %d\n", port)
	fmt.Printf("http://127.0.0.1:%d/\n", port)
}

func runLocal(args []string) int {
	var (
		dataDir         string
		port            string
		status          bool
		restart         bool
		importLocalGrok bool
		testSlack       bool
		slackChannel    string
	)

	_, err := lessFlags.
		String("--data-dir", &dataDir).
		String("--port", &port).
		Bool("--status", &status).
		Bool("--restart", &restart).
		Bool("--import-local-grok", &importLocalGrok).
		Bool("--test-slack", &testSlack).
		String("--slack-channel", &slackChannel).
		Help("-h,--help", runHelpText).
		HelpNoExit().
		Parse(args)
	if errors.Is(err, lessflags.ErrHelp) {
		return 0
	}
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		return 1
	}

	if slackChannel != "" && !testSlack {
		fmt.Fprintf(os.Stderr, "error: --slack-channel requires --test-slack\n")
		return 1
	}

	actionCount := 0
	if status {
		actionCount++
	}
	if restart {
		actionCount++
	}
	if importLocalGrok {
		actionCount++
	}
	if testSlack {
		actionCount++
	}
	if actionCount > 1 {
		fmt.Fprintf(os.Stderr, "error: use only one action flag at a time\n")
		return 1
	}

	if dataDir == "" {
		reg, regErr := loadRegistry()
		if regErr != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", regErr)
			return 1
		}
		if len(reg.DataDirs) == 0 && (status || testSlack) {
			running, err := localGatewayRunning("")
			if err != nil {
				fmt.Fprintf(os.Stderr, "error: %v\n", err)
				return 1
			}
			if !running {
				fmt.Fprintln(os.Stderr, "error: local gateway is not running")
			}
			fmt.Fprintln(os.Stderr, "error: requires --data-dir <dir>, or add with 'my openclaw add data-dir <dir>'")
			return 1
		}

		resolved, auto, err := resolveRunDataDir(dataDir)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			return 1
		}
		dataDir = resolved
		if auto {
			fmt.Printf("Using data dir: %s\n", dataDir)
		}
	}

	switch {
	case status:
		return runLocalStatus(dataDir, port)
	case restart:
		return runLocalRestart(dataDir, port)
	case importLocalGrok:
		if dataDir == "" {
			fmt.Fprintf(os.Stderr, "error: --data-dir is required\n")
			return 1
		}
		return runImportLocalGrok(dataDir, "")
	case testSlack:
		return runTestSlackLocal(dataDir, slackChannel)
	default:
		return runLocalLaunch(dataDir, port)
	}
}

func runLocalLaunch(dataDir, port string) int {
	if dataDir == "" {
		fmt.Fprintf(os.Stderr, "error: --data-dir is required\n")
		return 1
	}

	absPath, err := resolvePath(dataDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		return 1
	}

	if err := validateLocalDataDir(absPath); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		return 1
	}

	token, err := resolveToken(absPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		return 1
	}

	if !ensureOpenclawAvailable() {
		return 1
	}

	selectedPort, bumped, err := selectPort(port)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		return 1
	}

	autoPort := port == ""
	currentPort := selectedPort
	maxPort := defaultGatewayPort + maxPortScanTries - 1
	printLocalRunInfo(absPath, currentPort, bumped)

	for {
		portStr := fmt.Sprintf("%d", currentPort)
		gatewayArgs := []string{"gateway", "--bind", "lan", "--port", portStr}
		cmd, stderrBuf, err := openclawExecPreviewedForeground(absPath, token, gatewayArgs...)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			return 1
		}

		if err := writeGatewayBookkeeping(absPath, cmd.Process.Pid, currentPort); err != nil {
			_ = cmd.Process.Kill()
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			return 1
		}

		waitErr := cmd.Wait()
		_ = removeGatewayBookkeeping(absPath)

		if waitErr == nil {
			return 0
		}

		stderrOut := stderrBuf.String()

		if autoPort && isPortInUseError(stderrOut) && currentPort < maxPort {
			nextPort := currentPort + 1
			printLocalRunInfo(absPath, nextPort, true)
			currentPort = nextPort
			continue
		}

		if exitErr, ok := waitErr.(*exec.ExitError); ok {
			return exitErr.ExitCode()
		}
		fmt.Fprintf(os.Stderr, "error: %v\n", waitErr)
		return 1
	}
}