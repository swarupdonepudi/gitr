package upgrade

import (
	"archive/tar"
	"archive/zip"
	"bufio"
	"compress/gzip"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/swarupdonepudi/gitr/pkg/ui"
)

// UpgradeViaDirect upgrades the CLI by downloading directly from GitHub
func UpgradeViaDirect(version string) error {
	goos, goarch := GetPlatformInfo()
	downloadURL := BuildDownloadURL(version, goos, goarch)
	checksumURL := BuildChecksumURL(version)

	fmt.Println()
	ui.Info(fmt.Sprintf("Downloading gitr %s...", version))

	tempArchive, err := downloadToTemp(downloadURL)
	if err != nil {
		return fmt.Errorf("failed to download archive: %w", err)
	}
	defer os.Remove(tempArchive)

	ui.Success(fmt.Sprintf("Downloaded gitr %s", version))

	ui.Info("Verifying checksum...")
	if err := verifyChecksum(tempArchive, checksumURL, version, goos, goarch); err != nil {
		return fmt.Errorf("checksum verification failed: %w", err)
	}
	ui.Success("Checksum verified")

	ui.Info("Extracting binary...")
	tempDir, err := os.MkdirTemp("", "gitr-upgrade-*")
	if err != nil {
		return fmt.Errorf("failed to create temp directory: %w", err)
	}
	defer os.RemoveAll(tempDir)

	binaryPath, err := extractBinary(tempArchive, tempDir, goos)
	if err != nil {
		return fmt.Errorf("failed to extract binary: %w", err)
	}
	ui.Success("Extracted binary")

	currentBinary, err := os.Executable()
	if err != nil {
		return fmt.Errorf("failed to determine current binary path: %w", err)
	}

	currentBinary, err = filepath.EvalSymlinks(currentBinary)
	if err != nil {
		return fmt.Errorf("failed to resolve binary path: %w", err)
	}

	ui.Info("Installing...")
	if err := replaceBinary(binaryPath, currentBinary); err != nil {
		return err
	}
	ui.Success("Installed new binary")

	if runtime.GOOS == "darwin" {
		_ = exec.Command("xattr", "-dr", "com.apple.quarantine", currentBinary).Run()
	}

	return nil
}

func downloadToTemp(url string) (string, error) {
	var tempFile *os.File
	var err error
	if strings.HasSuffix(url, ".zip") {
		tempFile, err = os.CreateTemp("", "gitr-upgrade-*.zip")
	} else {
		tempFile, err = os.CreateTemp("", "gitr-upgrade-*.tar.gz")
	}
	if err != nil {
		return "", fmt.Errorf("failed to create temp file: %w", err)
	}
	tempPath := tempFile.Name()

	client := &http.Client{Timeout: httpTimeout * 6}
	resp, err := client.Get(url)
	if err != nil {
		tempFile.Close()
		os.Remove(tempPath)
		return "", fmt.Errorf("failed to download: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		tempFile.Close()
		os.Remove(tempPath)
		return "", fmt.Errorf("download failed: HTTP %d", resp.StatusCode)
	}

	_, err = io.Copy(tempFile, resp.Body)
	tempFile.Close()
	if err != nil {
		os.Remove(tempPath)
		return "", fmt.Errorf("failed to write downloaded file: %w", err)
	}

	return tempPath, nil
}

func extractBinary(archivePath, destDir, goos string) (string, error) {
	if goos == "windows" {
		return extractFromZip(archivePath, destDir)
	}
	return extractFromTarGz(archivePath, destDir)
}

func extractFromTarGz(archivePath, destDir string) (string, error) {
	file, err := os.Open(archivePath)
	if err != nil {
		return "", fmt.Errorf("failed to open archive: %w", err)
	}
	defer file.Close()

	gzReader, err := gzip.NewReader(file)
	if err != nil {
		return "", fmt.Errorf("failed to create gzip reader: %w", err)
	}
	defer gzReader.Close()

	tarReader := tar.NewReader(gzReader)

	var binaryPath string
	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return "", fmt.Errorf("failed to read tar entry: %w", err)
		}

		if header.Typeflag == tar.TypeReg && header.Name == "gitr" {
			binaryPath = filepath.Join(destDir, header.Name)
			outFile, err := os.OpenFile(binaryPath, os.O_CREATE|os.O_WRONLY, os.FileMode(header.Mode))
			if err != nil {
				return "", fmt.Errorf("failed to create binary file: %w", err)
			}
			if _, err := io.Copy(outFile, tarReader); err != nil {
				outFile.Close()
				return "", fmt.Errorf("failed to write binary file: %w", err)
			}
			outFile.Close()
			break
		}
	}

	if binaryPath == "" {
		return "", fmt.Errorf("binary not found in archive")
	}

	return binaryPath, nil
}

func extractFromZip(archivePath, destDir string) (string, error) {
	reader, err := zip.OpenReader(archivePath)
	if err != nil {
		return "", fmt.Errorf("failed to open zip archive: %w", err)
	}
	defer reader.Close()

	var binaryPath string
	for _, file := range reader.File {
		if file.Name == "gitr.exe" {
			binaryPath = filepath.Join(destDir, file.Name)
			rc, err := file.Open()
			if err != nil {
				return "", fmt.Errorf("failed to open file in archive: %w", err)
			}
			outFile, err := os.OpenFile(binaryPath, os.O_CREATE|os.O_WRONLY, file.Mode())
			if err != nil {
				rc.Close()
				return "", fmt.Errorf("failed to create binary file: %w", err)
			}
			if _, err := io.Copy(outFile, rc); err != nil {
				rc.Close()
				outFile.Close()
				return "", fmt.Errorf("failed to write binary file: %w", err)
			}
			rc.Close()
			outFile.Close()
			break
		}
	}

	if binaryPath == "" {
		return "", fmt.Errorf("binary not found in archive")
	}

	return binaryPath, nil
}

func verifyChecksum(archivePath, checksumURL, version, goos, goarch string) error {
	client := &http.Client{Timeout: httpTimeout}
	resp, err := client.Get(checksumURL)
	if err != nil {
		return fmt.Errorf("failed to download checksums: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to fetch checksums: HTTP %d", resp.StatusCode)
	}

	versionNum := strings.TrimPrefix(version, "v")
	var archiveName string
	if goos == "windows" {
		archiveName = fmt.Sprintf("gitr_%s_%s_%s.zip", versionNum, goos, goarch)
	} else {
		archiveName = fmt.Sprintf("gitr_%s_%s_%s.tar.gz", versionNum, goos, goarch)
	}

	var expectedChecksum string
	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Fields(line)
		if len(parts) >= 2 && parts[1] == archiveName {
			expectedChecksum = parts[0]
			break
		}
	}

	if expectedChecksum == "" {
		return fmt.Errorf("checksum not found for %s", archiveName)
	}

	file, err := os.Open(archivePath)
	if err != nil {
		return fmt.Errorf("failed to open downloaded archive: %w", err)
	}
	defer file.Close()

	hash := sha256.New()
	if _, err := io.Copy(hash, file); err != nil {
		return fmt.Errorf("failed to calculate checksum: %w", err)
	}

	actualChecksum := hex.EncodeToString(hash.Sum(nil))
	if actualChecksum != expectedChecksum {
		return fmt.Errorf("checksum mismatch: expected %s, got %s", expectedChecksum, actualChecksum)
	}

	return nil
}

func replaceBinary(newBinary, currentBinary string) error {
	if err := os.Chmod(newBinary, 0755); err != nil {
		return fmt.Errorf("failed to make binary executable: %w", err)
	}

	if err := os.Rename(newBinary, currentBinary); err != nil {
		if err := copyFile(newBinary, currentBinary); err != nil {
			if os.IsPermission(err) {
				return &PermissionError{Path: currentBinary, OrigErr: err}
			}
			return fmt.Errorf("failed to install new binary: %w", err)
		}
	}

	return nil
}

func copyFile(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	destFile, err := os.OpenFile(dst, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, sourceFile)
	return err
}

// PermissionError represents a permission error
type PermissionError struct {
	Path    string
	OrigErr error
}

func (e *PermissionError) Error() string {
	return fmt.Sprintf("permission denied: cannot write to %s", e.Path)
}
