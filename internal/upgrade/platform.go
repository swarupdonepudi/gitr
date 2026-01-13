package upgrade

import (
	"fmt"
	"os/exec"
	"runtime"
	"strings"
)

// UpgradeMethod represents the method used to upgrade the CLI
type UpgradeMethod int

const (
	MethodHomebrew UpgradeMethod = iota
	MethodDirectDownload
)

func (m UpgradeMethod) String() string {
	switch m {
	case MethodHomebrew:
		return "Homebrew"
	case MethodDirectDownload:
		return "Direct Download"
	default:
		return "Unknown"
	}
}

// DetectUpgradeMethod determines the best method to upgrade
func DetectUpgradeMethod() UpgradeMethod {
	if runtime.GOOS != "darwin" {
		return MethodDirectDownload
	}

	if _, err := exec.LookPath("brew"); err != nil {
		return MethodDirectDownload
	}

	cmd := exec.Command("brew", "list", "--cask", "gitr")
	if err := cmd.Run(); err != nil {
		return MethodDirectDownload
	}

	return MethodHomebrew
}

// GetPlatformInfo returns the current OS and architecture
func GetPlatformInfo() (goos string, goarch string) {
	return runtime.GOOS, runtime.GOARCH
}

// BuildDownloadURL constructs the download URL for a specific version and platform
func BuildDownloadURL(version, goos, goarch string) string {
	versionNum := strings.TrimPrefix(version, "v")
	var archiveName string
	if goos == "windows" {
		archiveName = fmt.Sprintf("gitr_%s_%s_%s.zip", versionNum, goos, goarch)
	} else {
		archiveName = fmt.Sprintf("gitr_%s_%s_%s.tar.gz", versionNum, goos, goarch)
	}
	return fmt.Sprintf("https://github.com/swarupdonepudi/gitr/releases/download/%s/%s", version, archiveName)
}

// BuildChecksumURL constructs the checksum file URL
func BuildChecksumURL(version string) string {
	versionNum := strings.TrimPrefix(version, "v")
	return fmt.Sprintf("https://github.com/swarupdonepudi/gitr/releases/download/%s/gitr_%s_checksums.txt", version, versionNum)
}
