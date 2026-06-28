package openclaw

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

const exportSettingsZipRoot = "openclaw"

func zipOpenClawDataDir(sourceDir, outputZip string) (int, error) {
	absSource, err := resolvePath(sourceDir)
	if err != nil {
		return 0, err
	}
	info, err := os.Stat(absSource)
	if err != nil {
		return 0, err
	}
	if !info.IsDir() {
		return 0, fmt.Errorf("path is not a directory: %s", absSource)
	}

	absOutput, err := resolvePath(outputZip)
	if err != nil {
		return 0, err
	}
	if err := os.MkdirAll(filepath.Dir(absOutput), 0o755); err != nil {
		return 0, err
	}

	out, err := os.Create(absOutput)
	if err != nil {
		return 0, err
	}
	defer out.Close()

	writer := zip.NewWriter(out)
	defer writer.Close()

	fileCount := 0
	err = filepath.WalkDir(absSource, func(path string, entry os.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}

		rel, err := filepath.Rel(absSource, path)
		if err != nil {
			return err
		}
		if rel == "." {
			return nil
		}

		if entry.Type()&os.ModeSymlink != 0 {
			return nil
		}
		if isSQLiteSidecar(entry.Name()) {
			return nil
		}

		zipName := filepath.ToSlash(filepath.Join(exportSettingsZipRoot, rel))
		if entry.IsDir() {
			_, err := writer.Create(zipName + "/")
			return err
		}

		if err := addFileToZip(writer, path, zipName, entry); err != nil {
			return err
		}
		fileCount++
		return nil
	})
	if err != nil {
		return 0, err
	}
	if err := writer.Close(); err != nil {
		return 0, err
	}
	if err := out.Close(); err != nil {
		return 0, err
	}
	return fileCount, nil
}

func addFileToZip(writer *zip.Writer, sourcePath, zipName string, entry os.DirEntry) error {
	info, err := entry.Info()
	if err != nil {
		return err
	}
	if !info.Mode().IsRegular() {
		return nil
	}

	header, err := zip.FileInfoHeader(info)
	if err != nil {
		return err
	}
	header.Name = zipName
	header.Method = zip.Deflate

	dest, err := writer.CreateHeader(header)
	if err != nil {
		return err
	}

	src, err := os.Open(sourcePath)
	if err != nil {
		return err
	}
	defer src.Close()

	_, err = io.Copy(dest, src)
	return err
}