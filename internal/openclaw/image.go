package openclaw

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/xhd2015/my/internal/openclaw/container"
)

const imageName = "my-openclaw:local"

func containerSpecHash() string {
	sum := sha256.Sum256(container.Containerfile)
	return "sha256:" + hex.EncodeToString(sum[:])
}

func imageExists() (bool, error) {
	out, err := podmanOutput("images")
	if err != nil {
		return false, err
	}
	return strings.Contains(out, imageName), nil
}

func needsRebuild(reg *Registry, force bool) (bool, error) {
	if force {
		return true, nil
	}
	exists, err := imageExists()
	if err != nil {
		return false, err
	}
	if !exists {
		return true, nil
	}
	current := containerSpecHash()
	if reg.Image == nil || reg.Image.SpecHash != current {
		return true, nil
	}
	return false, nil
}

func buildImage(reg *Registry) error {
	tmpDir, err := os.MkdirTemp("", "my-openclaw-build-*")
	if err != nil {
		return err
	}
	defer os.RemoveAll(tmpDir)

	dockerfilePath := filepath.Join(tmpDir, "Containerfile")
	if err := os.WriteFile(dockerfilePath, container.Containerfile, 0o644); err != nil {
		return err
	}

	if err := podmanRunPreviewed("build", "-f", dockerfilePath, "-t", imageName, tmpDir); err != nil {
		return fmt.Errorf("podman build: %w", err)
	}

	reg.Image = &ImageMeta{
		SpecHash: containerSpecHash(),
		BuiltAt:  time.Now().Format(time.RFC3339),
	}
	return saveRegistry(reg)
}