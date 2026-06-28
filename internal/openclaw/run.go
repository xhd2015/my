package openclaw

import (
	"errors"
	"fmt"
	"os"

	lessflags "github.com/xhd2015/less-flags"
)

const (
	defaultContainerName = "openclaw-gateway"
	defaultPort          = "18789"
	containerPort        = "18789"
)

func runInPodman(args []string) int {
	if len(args) > 0 && args[0] == "--exec" {
		return runExec(args[1:])
	}

	if helpText := runInPodmanHelpForArgs(args); helpText != "" {
		printHelp(helpText)
		return 0
	}

	var (
		dataDir       string
		rebuild       bool
		stop          bool
		restart       bool
		logs          bool
		status        bool
		showTokens       bool
		dashboard        bool
		importLocalGrok  bool
		installGrok         bool
		installSlackPlugin  bool
		testSlack           bool
		slackChannel        string
		doctorSlack         bool
		resolveSlackChannel string
		exportSettings      string
		containerName       string
		port                string
	)

	_, err := lessFlags.
		String("--data-dir", &dataDir).
		Bool("--rebuild", &rebuild).
		Bool("--stop", &stop).
		Bool("--restart", &restart).
		Bool("--logs", &logs).
		Bool("--status", &status).
		Bool("--show-tokens", &showTokens).
		Bool("--dashboard", &dashboard).
		Bool("--import-local-grok", &importLocalGrok).
		Bool("--install-grok", &installGrok).
		Bool("--install-slack-plugin", &installSlackPlugin).
		Bool("--test-slack", &testSlack).
		String("--slack-channel", &slackChannel).
		Bool("--doctor-slack", &doctorSlack).
		String("--resolve-slack-channel", &resolveSlackChannel).
		String("--export-settings", &exportSettings).
		String("--container-name", &containerName).
		String("--port", &port).
		Help("-h,--help", runInPodmanHelpText).
		HelpNoExit().
		Parse(args)
	if errors.Is(err, lessflags.ErrHelp) {
		return 0
	}
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		return 1
	}

	if containerName == "" {
		containerName = defaultContainerName
	}
	if port == "" {
		port = defaultPort
	}
	if slackChannel != "" && !testSlack {
		fmt.Fprintf(os.Stderr, "error: --slack-channel requires --test-slack\n")
		return 1
	}

	actionCount := 0
	if stop {
		actionCount++
	}
	if restart {
		actionCount++
	}
	if logs {
		actionCount++
	}
	if status {
		actionCount++
	}
	if showTokens {
		actionCount++
	}
	if dashboard {
		actionCount++
	}
	if importLocalGrok {
		actionCount++
	}
	if installGrok {
		actionCount++
	}
	if installSlackPlugin {
		actionCount++
	}
	if testSlack {
		actionCount++
	}
	if doctorSlack {
		actionCount++
	}
	if resolveSlackChannel != "" {
		actionCount++
	}
	if exportSettings != "" {
		actionCount++
	}
	if actionCount > 1 {
		fmt.Fprintf(os.Stderr, "error: use only one action flag at a time\n")
		return 1
	}

	switch {
	case restart:
		return runRestart(dataDir, rebuild, containerName, port)
	case installGrok:
		return runInstallGrok(containerName)
	case installSlackPlugin:
		return runInstallSlackPlugin(dataDir, containerName)
	case resolveSlackChannel != "":
		return runResolveSlackChannel(dataDir, containerName, resolveSlackChannel)
	case doctorSlack:
		return runDoctorSlack(dataDir, containerName)
	case testSlack:
		return runTestSlack(dataDir, containerName, slackChannel)
	case exportSettings != "":
		return runExportSettings(dataDir, containerName, exportSettings)
	case importLocalGrok:
		if err := ensurePodmanMachine(); err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			return 1
		}
		resolvedDataDir, err := resolveImportLocalGrokDataDir(dataDir, containerName)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			return 1
		}
		return runImportLocalGrok(resolvedDataDir, containerName)
	case stop:
		return runStop(containerName)
	case logs:
		return runLogs(containerName)
	case status:
		return runStatus(dataDir, port, containerName)
	case showTokens:
		if err := ensurePodmanMachine(); err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			return 1
		}
		resolvedDataDir, err := resolveContainerDataDir(dataDir, containerName)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			return 1
		}
		return runShowTokens(resolvedDataDir, port, containerName)
	case dashboard:
		return runDashboard(dataDir, port, containerName)
	}

	if dataDir == "" {
		fmt.Fprintf(os.Stderr, "error: --data-dir is required\n")
		return 1
	}

	if err := ensurePodmanMachine(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		return 1
	}

	running, err := containerIsRunning(containerName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		return 1
	}
	if running {
		fmt.Fprintf(os.Stderr, "error: %v\n", runningContainerError(containerName))
		return 1
	}

	return runLaunch(dataDir, rebuild, containerName, port)
}