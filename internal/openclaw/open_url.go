package openclaw

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
)

func openURL(url string) error {
	if os.Getenv("MY_OPENCLAW_NO_OPEN") == "1" {
		return nil
	}
	switch runtime.GOOS {
	case "darwin":
		return exec.Command("open", url).Run()
	case "linux":
		if _, err := exec.LookPath("xdg-open"); err == nil {
			return exec.Command("xdg-open", url).Run()
		}
	}
	return fmt.Errorf("no URL opener available on %s; open manually: %s", runtime.GOOS, url)
}