package install

import (
	osastro "astro/os"
	"fmt"
	"os"
)

func InstallVPN() error {
	if !osastro.IsRoot() {
		return fmt.Errorf("This script must be run as root")
	}

	if !tunAvailable() {
		return fmt.Errorf("The TUN device is not available. You need to enable TUN before running this script")
	}

	_, err := osastro.CheckOS()
	if err != nil {
		return err
	}

	err = installUnbound()
	if err != nil {
		return err
	}

	return nil
}

func UninstallVPN() error {
	if !osastro.IsRoot() {
		return fmt.Errorf("This script must be run as root")
	}

	err := removeUnbound()
	if err != nil {
		return err
	}

	return nil
}

func tunAvailable() bool {
	_, err := os.Stat("/dev/net/tun")
	if err != nil {
		if os.IsNotExist(err) {
			return false
		}
		fmt.Println("Error checking tun availability:", err)
		return false
	}

	return true
}
