package openclaw

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
)

func printCommand(cmd string) {
	fmt.Fprintf(os.Stderr, "$ %s\n", cmd)
}

func shellQuote(s string) string {
	if s == "" {
		return "''"
	}
	if strings.ContainsAny(s, " \t\n'\"$\\") {
		return "'" + strings.ReplaceAll(s, "'", "'\"'\"'") + "'"
	}
	return s
}

func formatPodmanCommand(args ...string) string {
	parts := make([]string, len(args)+1)
	parts[0] = "podman"
	for i, arg := range args {
		parts[i+1] = shellQuote(arg)
	}
	return strings.Join(parts, " ")
}

func podmanRunPreviewed(args ...string) error {
	printCommand(formatPodmanCommand(args...))
	return podmanRun(args...)
}

func formatOpenclawCommand(stateDir, token string, args ...string) string {
	parts := []string{
		"OPENCLAW_STATE_DIR=" + shellQuote(stateDir),
		"OPENCLAW_GATEWAY_TOKEN=" + shellQuote(token),
		"openclaw",
	}
	for _, arg := range args {
		parts = append(parts, shellQuote(arg))
	}
	return strings.Join(parts, " ")
}

func openclawExec(stateDir, token string, args ...string) error {
	cmd := exec.Command("openclaw", args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = append(os.Environ(),
		"OPENCLAW_STATE_DIR="+stateDir,
		"OPENCLAW_GATEWAY_TOKEN="+token,
	)
	return cmd.Run()
}

func openclawExecPreviewed(stateDir, token string, args ...string) error {
	printCommand(formatOpenclawCommand(stateDir, token, args...))
	return openclawExec(stateDir, token, args...)
}

func openclawExecDetached(stateDir, token string, args ...string) (*exec.Cmd, error) {
	printCommand(formatOpenclawCommand(stateDir, token, args...))
	cmd := exec.Command("openclaw", args...)
	cmd.Env = append(os.Environ(),
		"OPENCLAW_STATE_DIR="+stateDir,
		"OPENCLAW_GATEWAY_TOKEN="+token,
	)
	if err := cmd.Start(); err != nil {
		return nil, err
	}
	if os.Getenv("OPENCLAW_STUB_LOG") != "" {
		_ = cmd.Wait()
	}
	return cmd, nil
}

func openclawExecPreviewedForeground(stateDir, token string, args ...string) (*exec.Cmd, *strings.Builder, error) {
	printCommand(formatOpenclawCommand(stateDir, token, args...))

	var stderrBuf strings.Builder
	cmd := exec.Command("openclaw", args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = io.MultiWriter(os.Stderr, &stderrBuf)
	cmd.Env = append(os.Environ(),
		"OPENCLAW_STATE_DIR="+stateDir,
		"OPENCLAW_GATEWAY_TOKEN="+token,
	)
	if err := cmd.Start(); err != nil {
		return nil, &stderrBuf, err
	}
	return cmd, &stderrBuf, nil
}

func openclawExecPreviewedCapture(stateDir, token string, args ...string) (string, error) {
	printCommand(formatOpenclawCommand(stateDir, token, args...))

	var stderrBuf strings.Builder
	cmd := exec.Command("openclaw", args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = io.MultiWriter(os.Stderr, &stderrBuf)
	cmd.Env = append(os.Environ(),
		"OPENCLAW_STATE_DIR="+stateDir,
		"OPENCLAW_GATEWAY_TOKEN="+token,
	)
	err := cmd.Run()
	return stderrBuf.String(), err
}