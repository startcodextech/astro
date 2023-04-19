package install

import (
	os2 "astro/os"
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"regexp"
)

func installUnbound() error {

	configPath := "/etc/unbound/unbound.conf"

	OS := os2.GetLinuxOsFamily()

	_, err := os.Stat(configPath)

	if os.IsNotExist(err) {
		switch OS {
		case "debian", "ubuntu":

			err := exec.Command("apt-get", "install", "-y", "unbound").Run()
			if err != nil {
				return err
			}

			file, err := os.OpenFile(configPath, os.O_CREATE|os.O_WRONLY, 0644)
			if err != nil {
				return err
			}
			defer file.Close()
			config := `interface: 10.8.0.1
access-control: 10.8.0.1/24 allow
hide-identity: yes
hide-version: yes
use-caps-for-id: yes
prefetch: yes`
			writer := bufio.NewWriter(file)
			_, err = writer.WriteString(config)
			if err != nil {
				return err
			}
			err = writer.Flush()
			if err != nil {
				return err
			}
		case "centos", "amzn", "oracle":
			err := exec.Command("yum", "install", "-y", "unbound").Run()
			if err != nil {
				return err
			}
			file, err := os.ReadFile(configPath)
			if err != nil {
				return err
			}
			config := []*regexp.Regexp{
				regexp.MustCompile(`# interface: 0\.0\.0\.0`),
				regexp.MustCompile(`# access-control: 127\.0\.0\.0/8 allow`),
				regexp.MustCompile(`# hide-identity: no`),
				regexp.MustCompile(`# hide-version: no`),
				regexp.MustCompile(`use-caps-for-id: no`),
			}

			values := [][]byte{
				[]byte("interface: 10.8.0.1"),
				[]byte("access-control: 10.8.0.1/24 allow"),
				[]byte("hide-identity: yes"),
				[]byte("hide-version: yes"),
				[]byte("use-caps-for-id: yes"),
			}

			for i, r := range config {
				file = r.ReplaceAll(file, values[i])
			}

			err = os.WriteFile(configPath, file, 0644)
			if err != nil {
				return err
			}
		case "fedora":
			err := exec.Command("dnf", "install", "-y", "unbound").Run()
			if err != nil {
				return err
			}
			file, err := os.ReadFile(configPath)
			if err != nil {
				return err
			}
			config := []*regexp.Regexp{
				regexp.MustCompile(`# interface: 0\.0\.0\.0`),
				regexp.MustCompile(`# access-control: 127\.0\.0\.0/8 allow`),
				regexp.MustCompile(`# hide-identity: no`),
				regexp.MustCompile(`# hide-version: no`),
				regexp.MustCompile(`use-caps-for-id: no`),
			}

			values := [][]byte{
				[]byte("interface: 10.8.0.1"),
				[]byte("access-control: 10.8.0.1/24 allow"),
				[]byte("hide-identity: yes"),
				[]byte("hide-version: yes"),
				[]byte("use-caps-for-id: yes"),
			}

			for i, r := range config {
				file = r.ReplaceAll(file, values[i])
			}

			err = os.WriteFile(configPath, file, 0644)
			if err != nil {
				return err
			}
		case "arch":
			err := exec.Command("pacman", "--noconfirm", "-S", "unbound").Run()
			if err != nil {
				return err
			}

			err = exec.Command("curl", "-o", "/etc/unbound/root.hints", "https://www.internic.net/domain/named.cache").Run()
			if err != nil {
				return err
			}

			if _, err := os.Stat("/etc/unbound/unbound.conf.old"); os.IsNotExist(err) {
				err = os.Rename("/etc/unbound/unbound.conf", "/etc/unbound/unbound.conf.old")
				if err != nil {
					log.Fatal(err)
				}
			}

			unboundConf := `server:
use-syslog: yes
do-daemonize: no
username: "unbound"
directory: "/etc/unbound"
trust-anchor-file: trusted-key.key
root-hints: root.hints
interface: 10.8.0.1
access-control: 10.8.0.1/24 allow
port: 53
num-threads: 2
use-caps-for-id: yes
harden-glue: yes
hide-identity: yes
hide-version: yes
qname-minimisation: yes
prefetch: yes
	`
			err = os.WriteFile(configPath, []byte(unboundConf), 0644)
			if err != nil {
				log.Fatal(err)
			}
		default:
			config := []string{
				"private-address: 10.0.0.0/8",
				"private-address: fd42:42:42:42::/112",
				"private-address: 172.16.0.0/12",
				"private-address: 192.168.0.0/16",
				"private-address: 169.254.0.0/16",
				"private-address: fd00::/8",
				"private-address: fe80::/10",
				"private-address: 127.0.0.0/8",
				"private-address: ::ffff:0:0/96",
			}

			file, err := os.OpenFile(configPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
			if err != nil {
				return err
			}
			defer file.Close()

			writer := bufio.NewWriter(file)
			for _, line := range config {
				_, err = writer.WriteString(line + "\n")
				if err != nil {
					return err
				}
			}
		}
	} else {
		file, err := os.OpenFile(configPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return err
		}

		defer file.Close()
		if _, err = file.WriteString("include: /etc/unbound/openvpn.conf"); err != nil {
			return err
		}

		content := []byte(`server:
interface: 10.8.0.1
access-control: 10.8.0.1/24 allow
hide-identity: yes
hide-version: yes
use-caps-for-id: yes
prefetch: yes
private-address: 10.0.0.0/8
private-address: fd42:42:42:42::/112
private-address: 172.16.0.0/12
private-address: 192.168.0.0/16
private-address: 169.254.0.0/16
private-address: fd00::/8
private-address: fe80::/10
private-address: 127.0.0.0/8
private-address: ::ffff:0:0/96`)
		err = os.WriteFile("/etc/unbound/openvpn.conf", content, 0644)
		if err != nil {
			return err
		}
	}

	err = exec.Command("systemctl", "enable", "unbound").Run()
	if err != nil {
		return err
	}

	return exec.Command("systemctl", "restart", "unbound").Run()
}

func removeUnbound() error {

	OS := os2.GetLinuxOsFamily()

	// Remove OpenVPN-related config
	err := exec.Command("sed", "-i", "/include: \\/etc\\/unbound\\/openvpn.conf/d", "/etc/unbound/unbound.conf").Run()
	if err != nil {
		return err
	}

	err = os.Remove("/etc/unbound/openvpn.conf")
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	// Stop Unbound
	err = exec.Command("systemctl", "stop", "unbound").Run()
	if err != nil {
		return err
	}

	switch OS {
	case "debian", "ubuntu":
		err = exec.Command("apt-get", "remove", "--purge", "-y", "unbound").Run()
	case "arch":
		err = exec.Command("pacman", "--noconfirm", "-R", "unbound").Run()
	case "centos", "amzn", "oracle":
		err = exec.Command("yum", "remove", "-y", "unbound").Run()
	case "fedora":
		err = exec.Command("dnf", "remove", "-y", "unbound").Run()
	default:
		return fmt.Errorf("unsupported OS: %s", OS)
	}

	if err != nil {
		return err
	}

	err = os.RemoveAll("/etc/unbound/")
	if err != nil {
		return err
	}

	log.Println("Unbound removed!")

	return nil
}
