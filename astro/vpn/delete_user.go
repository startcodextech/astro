package vpn

import (
	os2 "astro/os"
	"astro/vpn/types"
	"fmt"
	"github.com/dgraph-io/badger/v3"
	"os"
	"os/exec"
	"strings"
)

func DeleteUser(db *badger.DB, username string) error {
	exists, err := validateUserExists(username)
	if err != nil {
		return err
	}

	if !exists {
		return fmt.Errorf("user %s does not exist", username)
	}

	return revokeClientCertificate(db, username)
}

func revokeClientCertificate(db *badger.DB, username string) error {
	cmds := []*exec.Cmd{
		exec.Command(os2.ASTRO_VPN_EASYRSA_BIN_PATH, "--batch", "revoke", username),
		exec.Command(os2.ASTRO_VPN_EASYRSA_BIN_PATH, "gen-crl"),
		exec.Command("cp", os2.ASTRO_VPN_EASYRSA_CRL_PATH, os2.ASTRO_VPN_CRL_PATH),
		exec.Command("chmod", "644", os2.ASTRO_VPN_CRL_PATH),
		exec.Command("find", os2.ASTRO_VPN_EASYRSA_PATH, "-maxdepth", "2", "-name", fmt.Sprintf("%s.ovpn", username), "-delete"),
		exec.Command("sed", "-i", fmt.Sprintf("/^%s,.*/d", username), os2.ASTRO_VPN_IPP),
		exec.Command("cp", os2.ASTRO_VPN_PKI_PATH, os2.ASTRO_VPN_PKI_PATH+".bk"),
	}

	for _, cmd := range cmds {
		cmd.Dir = os2.ASTRO_VPN_EASYRSA_PATH
		err := cmd.Run()
		if err != nil {
			return err
		}
	}

	err := deleteCert(db, username)
	if err != nil {
		return err
	}

	return nil
}

func deleteCert(db *badger.DB, username string) error {

	linesPki, err := getPkiIndexFile()
	if err != nil {
		return err
	}

	var user types.User
	lines := ""

	for _, line := range linesPki {
		if strings.Contains(line, "V") {
			lines += line + "\n"
		}
		if strings.Contains(line, username) && strings.Contains(line, "R") {
			user, err = parseLine(line)
			if err != nil {
				return err
			}
		}
	}

	err = os.WriteFile(os2.ASTRO_VPN_PKI_PATH, []byte(lines), 0600)
	if err != nil {
		return err
	}

	filesRemoved := []string{
		fmt.Sprintf(os2.ASTRO_VPN_EASYRSA_REVOKED_CRT_PATH, user.Serial),
		fmt.Sprintf(os2.ASTRO_VPN_EASYRSA_REVOKED_KEY_PATH, user.Serial),
		fmt.Sprintf(os2.ASTRO_VPN_EASYRSA_REVOKED_REQ_PATH, user.Serial),
	}

	for _, file := range filesRemoved {
		err = os.Remove(file)
		if err != nil {
			return err
		}
	}

	return RemoveUser(db, username)
}
