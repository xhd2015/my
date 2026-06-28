package openclaw

import (
	"fmt"
	"os"
	"strings"
)

const (
	containerGrokDir      = "/home/node/.grok"
	containerGrokAuthPath = "/home/node/.grok/auth.json"
	containerGrokBin      = "/home/node/.grok/bin/grok"
	grokCLIInstallURL     = "https://x.ai/cli/install.sh"
)

const grokCLIInstallShell = "curl -fsSL " + grokCLIInstallURL + " | bash"

// Tolerates Podman VM / container clock skew during apt bootstrap.
const grokCLIInstallBootstrap = "DEBIAN_FRONTEND=noninteractive apt-get -o Acquire::Check-Valid-Until=false -o Acquire::Check-Date=false update -qq && " +
	"DEBIAN_FRONTEND=noninteractive apt-get install -y -qq curl ca-certificates && " +
	`su -s /bin/bash node -c "` + grokCLIInstallShell + `"`

type grokInstallPlan struct {
	execArgs []string
}

func containerHasCurl(containerName string) (bool, error) {
	out, err := podmanOutput("exec", containerName, "sh", "-c", "command -v curl")
	if err != nil {
		return false, err
	}
	return strings.TrimSpace(out) != "", nil
}

func planGrokInstall(containerName string) (grokInstallPlan, error) {
	hasCurl, err := containerHasCurl(containerName)
	if err != nil {
		return grokInstallPlan{}, err
	}
	if hasCurl {
		return grokInstallPlan{
			execArgs: []string{"exec", containerName, "sh", "-c", grokCLIInstallShell},
		}, nil
	}
	return grokInstallPlan{
		execArgs: []string{"exec", "--user", "root", containerName, "sh", "-c", grokCLIInstallBootstrap},
	}, nil
}

func (p grokInstallPlan) previewCommand() string {
	return formatPodmanCommand(p.execArgs...)
}

func copyGrokAuthToContainer(containerName, grokAuthPath string) error {
	running, err := containerIsRunning(containerName)
	if err != nil {
		return err
	}
	if !running {
		return nil
	}

	absPath, err := resolvePath(grokAuthPath)
	if err != nil {
		return err
	}
	if _, err := os.Stat(absPath); err != nil {
		return fmt.Errorf("grok auth file not found: %s", absPath)
	}

	if err := podmanRunPreviewed("exec", containerName, "mkdir", "-p", containerGrokDir); err != nil {
		return err
	}
	dest := containerName + ":" + containerGrokAuthPath
	return podmanRunPreviewed("cp", absPath, dest)
}

func grokRunHint(containerName string) string {
	return fmt.Sprintf("  podman exec -it %s %s", containerName, containerGrokBin)
}