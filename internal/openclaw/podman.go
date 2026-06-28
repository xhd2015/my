package openclaw

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func stdinIsTerminal() bool {
	switch os.Getenv("MY_OPENCLAW_EXEC_INTERACTIVE") {
	case "0", "false", "no":
		return false
	case "1", "true", "yes":
		return true
	}
	fi, err := os.Stdin.Stat()
	if err != nil {
		return false
	}
	return fi.Mode()&os.ModeCharDevice != 0
}

func podmanRun(args ...string) error {
	cmd := exec.Command("podman", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func podmanExecInteractive(containerName string, execArgs ...string) error {
	args := []string{"exec"}
	if stdinIsTerminal() {
		args = append(args, "-it")
	}
	args = append(args, containerName)
	args = append(args, execArgs...)

	cmd := exec.Command("podman", args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func podmanExecInteractivePreviewed(containerName string, execArgs ...string) error {
	previewArgs := []string{"exec"}
	if stdinIsTerminal() {
		previewArgs = append(previewArgs, "-it")
	}
	previewArgs = append(previewArgs, containerName)
	previewArgs = append(previewArgs, execArgs...)
	printCommand(formatPodmanCommand(previewArgs...))
	return podmanExecInteractive(containerName, execArgs...)
}

func podmanOutput(args ...string) (string, error) {
	cmd := exec.Command("podman", args...)
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return string(out), nil
}

func ensurePodmanMachine() error {
	if os.Getenv("MY_OPENCLAW_CHECK_PODMAN_MACHINE") != "1" {
		return nil
	}
	out, err := podmanOutput("machine", "info")
	if err != nil {
		return err
	}
	if strings.Contains(out, "Running: false") {
		return podmanRun("machine", "start")
	}
	return nil
}

func containerIsRunning(name string) (bool, error) {
	out, err := podmanOutput("ps", "--filter", "name=^"+name+"$", "--format", "{{.Names}}")
	if err != nil {
		return false, err
	}
	for _, line := range strings.Split(out, "\n") {
		if strings.TrimSpace(line) == name {
			return true, nil
		}
	}
	return false, nil
}

func stopAndRemoveContainer(name string) error {
	running, err := containerIsRunning(name)
	if err != nil {
		return err
	}
	if !running {
		return nil
	}
	if err := podmanRunPreviewed("stop", name); err != nil {
		return err
	}
	return podmanRunPreviewed("rm", name)
}

func showContainerLogs(name string) error {
	return podmanRunPreviewed("logs", name)
}

func containerHostPort(name, containerPort string) (string, error) {
	format := fmt.Sprintf(`{{with (index .NetworkSettings.Ports "%s/tcp")}}{{(index . 0).HostPort}}{{end}}`, containerPort)
	out, err := podmanOutput("inspect", name, "--format", format)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(out), nil
}

func resolveGatewayPort(explicitPort, containerName string) string {
	if explicitPort != defaultPort {
		return explicitPort
	}
	hostPort, err := containerHostPort(containerName, containerPort)
	if err == nil && hostPort != "" {
		return hostPort
	}
	return explicitPort
}

func resolveGatewayToken(dataDir, containerName string) string {
	if dataDir != "" {
		if sources, err := resolveTokenSources(dataDir); err == nil {
			if token := sources.Effective(); token != "" {
				return token
			}
		}
	}
	if running, err := containerIsRunning(containerName); err == nil && running {
		if token, err := containerEnvToken(containerName); err == nil {
			return token
		}
	}
	return ""
}

func containerEnvToken(name string) (string, error) {
	out, err := podmanOutput("exec", name, "sh", "-c", "printf %s \"$OPENCLAW_GATEWAY_TOKEN\"")
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(out), nil
}

func runningContainerError(name string) error {
	return fmt.Errorf("container %s is already running (use 'my openclaw run-in-podman --stop' or '--restart')", name)
}