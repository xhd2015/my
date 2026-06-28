package openclaw

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type DataDirEntry struct {
	Path    string `json:"path"`
	Note    string `json:"note,omitempty"`
	AddedAt string `json:"added_at"`
}

type ImageMeta struct {
	SpecHash string `json:"spec_hash"`
	BuiltAt  string `json:"built_at"`
}

type Registry struct {
	DataDirs []DataDirEntry `json:"data_dirs"`
	Image    *ImageMeta     `json:"image,omitempty"`
}

func loadRegistry() (*Registry, error) {
	path, err := registryPath()
	if err != nil {
		return nil, err
	}
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return &Registry{}, nil
		}
		return nil, err
	}
	var reg Registry
	if err := json.Unmarshal(data, &reg); err != nil {
		return nil, err
	}
	return &reg, nil
}

func saveRegistry(reg *Registry) error {
	path, err := registryPath()
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(reg, "", "  ")
	if err != nil {
		return err
	}
	data = append(data, '\n')
	return os.WriteFile(path, data, 0o644)
}

func findDataDir(reg *Registry, path string) (int, *DataDirEntry) {
	for i := range reg.DataDirs {
		if reg.DataDirs[i].Path == path {
			return i, &reg.DataDirs[i]
		}
	}
	return -1, nil
}