package openclaw

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

const slackTestMessagePrefix = "OpenClaw test message from my CLI"

type slackDoctorCheck struct {
	name    string
	ok      bool
	detail  string
	missing bool
}

type slackConfigView struct {
	enabled       bool
	mode          string
	botToken      string
	appToken      string
	signingSecret string
	relayURL      string
	relayAuth     string
	testTarget    string
	channelIDs    []string
}

func runDoctorSlack(dataDir, containerName string) int {
	if err := ensurePodmanMachine(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		return 1
	}

	resolvedDataDir, err := requireRunningContainerDataDir(dataDir, containerName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		return 1
	}

	absPath, err := resolvePath(resolvedDataDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		return 1
	}

	view, checks, err := diagnoseSlackConfig(absPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		return 1
	}

	fmt.Printf("Checking Slack configuration in %s\n", absPath)
	fmt.Printf("Mode: %s\n", view.mode)
	if view.testTarget != "" {
		fmt.Printf("Test target: %s\n", view.testTarget)
	}
	fmt.Println()

	for _, check := range checks {
		mark := "ok"
		if !check.ok {
			mark = "missing"
		}
		fmt.Printf("[%s] %s", mark, check.name)
		if check.detail != "" {
			fmt.Printf(": %s", check.detail)
		}
		fmt.Println()
	}

	for _, check := range checks {
		if check.missing {
			fmt.Println()
			fmt.Println("Slack is not correctly configured.")
			return 1
		}
	}

	fmt.Println()
	fmt.Printf("Probing Slack via running container %s\n", containerName)
	execArgs := []string{"openclaw", "channels", "status", "--probe", "--channel", "slack"}
	printCommand(formatPodmanCommand(append([]string{"exec", containerName}, execArgs...)...))
	if err := podmanExecInteractive(containerName, execArgs...); err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			return exitErr.ExitCode()
		}
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		return 1
	}

	fmt.Println()
	fmt.Println("Slack configuration looks correct.")
	return 0
}

func runTestSlack(dataDir, containerName, slackChannel string) int {
	if err := ensurePodmanMachine(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		return 1
	}

	resolvedDataDir, err := requireRunningContainerDataDir(dataDir, containerName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		return 1
	}

	absPath, err := resolvePath(resolvedDataDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		return 1
	}

	return runTestSlackCore(absPath, slackChannel, containerName)
}

func runTestSlackLocal(dataDir, slackChannel string) int {
	resolvedDataDir, err := requireRunningLocalGateway(dataDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		return 1
	}

	absPath, err := resolvePath(resolvedDataDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		return 1
	}

	return runTestSlackCore(absPath, slackChannel, "")
}

func runTestSlackCore(absPath, slackChannel, containerName string) int {
	view, checks, err := diagnoseSlackConfig(absPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		return 1
	}
	for _, check := range checks {
		if check.missing {
			fmt.Fprintf(os.Stderr, "error: Slack is not correctly configured (%s)\n", check.name)
			if containerName != "" {
				fmt.Fprintf(os.Stderr, "hint: run 'my openclaw run-in-podman --doctor-slack --container-name %s'\n", containerName)
			}
			return 1
		}
	}

	if view.botToken == "" {
		fmt.Fprintf(os.Stderr, "error: slack bot token not found in %s\n", absPath)
		return 1
	}

	targetID, targetLabel, err := resolveSlackTestTarget(view, view.botToken, slackChannel)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		printSlackTestTargetHelp(absPath)
		return 1
	}

	message := fmt.Sprintf("%s (%s).", slackTestMessagePrefix, time.Now().Format(time.RFC3339))
	fmt.Printf("Sending Slack test message to %s (%s)\n", targetLabel, targetID)
	if err := sendSlackChannelMessage(view.botToken, targetID, message); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		return 1
	}

	fmt.Println()
	fmt.Println("Test message sent.")
	fmt.Printf("Open Slack and check %s for the test message.\n", targetLabel)
	return 0
}

func diagnoseSlackConfig(dataDir string) (slackConfigView, []slackDoctorCheck, error) {
	cfg, env, err := loadOpenClawConfigMap(dataDir)
	if err != nil {
		return slackConfigView{}, nil, err
	}

	view := parseSlackConfigView(cfg, env)
	checks := buildSlackDoctorChecks(view)
	return view, checks, nil
}

func loadOpenClawConfigMap(dataDir string) (map[string]any, map[string]string, error) {
	configPath := filepath.Join(dataDir, "openclaw.json")
	data, err := os.ReadFile(configPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil, fmt.Errorf("openclaw.json not found in %s", dataDir)
		}
		return nil, nil, err
	}

	var cfg map[string]any
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, nil, fmt.Errorf("parse openclaw.json: %w", err)
	}

	env, err := readDotEnv(filepath.Join(dataDir, ".env"))
	if err != nil {
		return nil, nil, err
	}
	return cfg, env, nil
}

func readDotEnv(path string) (map[string]string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return map[string]string{}, nil
		}
		return nil, err
	}

	env := make(map[string]string)
	for _, line := range strings.Split(string(data), "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		key, value, ok := strings.Cut(line, "=")
		if !ok {
			continue
		}
		key = strings.TrimSpace(key)
		value = strings.TrimSpace(value)
		if key != "" && value != "" {
			env[key] = value
		}
	}
	return env, nil
}

func parseSlackConfigView(cfg map[string]any, env map[string]string) slackConfigView {
	channels, _ := cfg["channels"].(map[string]any)
	slack, _ := channels["slack"].(map[string]any)
	if slack == nil {
		return slackConfigView{mode: "socket"}
	}

	view := slackConfigView{
		enabled:    slackEnabled(slack),
		mode:       slackMode(slack),
		botToken:   resolveSlackCredential(slack, "botToken", env, "SLACK_BOT_TOKEN"),
		appToken:   resolveSlackCredential(slack, "appToken", env, "SLACK_APP_TOKEN"),
		signingSecret: resolveSlackCredential(slack, "signingSecret", env, "SLACK_SIGNING_SECRET"),
		testTarget: strings.TrimSpace(stringFromAny(slack["testTarget"])),
	}

	if relay, ok := slack["relay"].(map[string]any); ok {
		view.relayURL = strings.TrimSpace(stringFromAny(relay["url"]))
		view.relayAuth = resolveSlackCredential(relay, "authToken", env, "SLACK_RELAY_AUTH_TOKEN")
	}

	if channelsMap, ok := slack["channels"].(map[string]any); ok {
		for id := range channelsMap {
			id = strings.TrimSpace(id)
			if id != "" {
				view.channelIDs = append(view.channelIDs, id)
			}
		}
	}

	return view
}

func buildSlackDoctorChecks(view slackConfigView) []slackDoctorCheck {
	checks := []slackDoctorCheck{
		{name: "channels.slack configured", ok: view.enabled || view.botToken != "" || view.appToken != "" || view.signingSecret != "" || len(view.channelIDs) > 0},
	}

	if !checks[0].ok {
		checks[0].missing = true
		checks[0].detail = "add channels.slack to openclaw.json"
		return append(checks, slackDoctorCheck{
			name:    "channels.slack enabled",
			ok:      false,
			missing: true,
			detail:  "set channels.slack.enabled = true",
		})
	}

	enabledCheck := slackDoctorCheck{name: "channels.slack enabled", ok: view.enabled}
	if !enabledCheck.ok {
		enabledCheck.missing = true
		enabledCheck.detail = "set channels.slack.enabled = true"
	}
	checks = append(checks, enabledCheck)

	switch view.mode {
	case "http":
		checks = append(checks, credentialCheck("bot token", view.botToken))
		checks = append(checks, credentialCheck("signing secret", view.signingSecret))
	case "relay":
		checks = append(checks, credentialCheck("bot token", view.botToken))
		relayURL := slackDoctorCheck{name: "relay URL", ok: view.relayURL != ""}
		if !relayURL.ok {
			relayURL.missing = true
			relayURL.detail = "set channels.slack.relay.url"
		}
		checks = append(checks, relayURL)
		checks = append(checks, credentialCheck("relay auth token", view.relayAuth))
	default:
		checks = append(checks, credentialCheck("bot token", view.botToken))
		checks = append(checks, credentialCheck("app token", view.appToken))
	}

	return checks
}

func credentialCheck(name, value string) slackDoctorCheck {
	check := slackDoctorCheck{name: name, ok: value != ""}
	if !check.ok {
		check.missing = true
		check.detail = "configure in openclaw.json or .env"
	}
	return check
}

func printSlackTestTargetHelp(configPath string) {
	fmt.Fprintln(os.Stderr)
	fmt.Fprintln(os.Stderr, "--test-slack could not find a destination.")
	fmt.Fprintf(os.Stderr, "By default it sends to #%s when no target is configured.\n", defaultSlackTestChannel)
	fmt.Fprintln(os.Stderr)
	fmt.Fprintln(os.Stderr, "Override the destination in openclaw.json:")
	fmt.Fprintln(os.Stderr, `  "channels": {`)
	fmt.Fprintln(os.Stderr, `    "slack": {`)
	fmt.Fprintln(os.Stderr, `      "testTarget": "C01234567"`)
	fmt.Fprintln(os.Stderr, `    }`)
	fmt.Fprintln(os.Stderr, `  }`)
	fmt.Fprintln(os.Stderr)
	fmt.Fprintf(os.Stderr, "Or pass a channel on the command line:\n")
	fmt.Fprintf(os.Stderr, "  my openclaw run-in-podman --test-slack --slack-channel general\n")
	fmt.Fprintf(os.Stderr, "  my openclaw run-in-podman --resolve-slack-channel \"#release-notes\"\n")
	fmt.Fprintln(os.Stderr)
	fmt.Fprintf(os.Stderr, "Config file: %s\n", configPath)
}

func slackEnabled(slack map[string]any) bool {
	if enabled, ok := slack["enabled"].(bool); ok {
		return enabled
	}
	return true
}

func slackMode(slack map[string]any) string {
	mode := strings.TrimSpace(stringFromAny(slack["mode"]))
	if mode == "" {
		return "socket"
	}
	return mode
}

func resolveSlackCredential(section map[string]any, key string, env map[string]string, fallbackEnvKey string) string {
	if value := credentialValue(section[key], env); value != "" {
		return value
	}
	return strings.TrimSpace(env[fallbackEnvKey])
}

func credentialValue(raw any, env map[string]string) string {
	switch value := raw.(type) {
	case string:
		return strings.TrimSpace(value)
	case map[string]any:
		id := strings.TrimSpace(stringFromAny(value["id"]))
		if id == "" {
			return ""
		}
		source := strings.TrimSpace(stringFromAny(value["source"]))
		if source == "" || source == "env" {
			return strings.TrimSpace(env[id])
		}
	}
	return ""
}

func stringFromAny(raw any) string {
	switch value := raw.(type) {
	case string:
		return value
	default:
		return ""
	}
}