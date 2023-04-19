package os

import (
	"astro/os/types"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func CheckOS() (types.LinuxDistro, error) {
	var distro types.LinuxDistro

	if _, err := os.Stat("/etc/debian_version"); err == nil {
		distro.Name = "debian"
	} else if _, err := os.Stat("/etc/system-release"); err == nil {
	} else if _, err := os.Stat("/etc/arch-release"); err == nil {
		distro.Name = "arch"
	} else {
		return distro, fmt.Errorf("It looks like you are not running this installer on a Debian, Ubuntu, Fedora, CentOS, Amazon Linux 2, Oracle Linux 8, or Arch Linux system")
	}

	content, err := os.ReadFile("/etc/os-release")
	if err != nil {
		return distro, fmt.Errorf("could not read /etc/os-release: %v", err)
	}

	lines := strings.Split(string(content), "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "ID=") {
			distro.Name = strings.Trim(line, "ID=")
		} else if strings.HasPrefix(line, "VERSION_ID=") {
			distro.Version = strings.Trim(line, "VERSION_ID=")
			distro.MajorVersion = strings.Split(distro.Version, ".")[0]
		}
	}

	switch distro.Name {
	case "debian", "raspbian":
		if distro.MajorVersion < "9" {
			distro.Unsupported = true
		}
	case "ubuntu":
		if distro.MajorVersion < "16" {
			distro.Unsupported = true
		}
	case "centos", "rocky", "almalinux":
		if distro.MajorVersion < "7" {
			distro.Unsupported = true
			distro.OnlySupported = "CentOS 7 y CentOS 8"
		}
	case "ol":
		if distro.MajorVersion != "8" {
			distro.Unsupported = true
			distro.OnlySupported = "Oracle Linux 8"
		}
	case "amzn":
		if distro.MajorVersion != "2" {
			distro.Unsupported = true
			distro.OnlySupported = "Amazon Linux 2"
		}
	}

	return distro, nil
}

func GetLinuxOsFamily() string {
	cmdOutput, _ := exec.Command("cat", "/etc/os-release").Output()
	osRelease := string(cmdOutput)
	if strings.Contains(osRelease, "Debian") {
		return "debian"
	} else if strings.Contains(osRelease, "Ubuntu") {
		return "ubuntu"
	} else if strings.Contains(osRelease, "CentOS") {
		return "centos"
	} else if strings.Contains(osRelease, "Amazon Linux") {
		return "amzn"
	} else if strings.Contains(osRelease, "Oracle Linux") {
		return "oracle"
	} else if strings.Contains(osRelease, "Fedora") {
		return "fedora"
	} else if strings.Contains(osRelease, "Arch Linux") {
		return "arch"
	}
	return "unknown"
}
