package vpn

import (
	os2 "astro/os"
	"astro/vpn/types"
	"bufio"
	"bytes"
	"fmt"
	"github.com/creack/pty"
	"github.com/dgraph-io/badger/v3"
	"io"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"time"
)

func validateClientName(name string) bool {
	validName := regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)
	return validName.MatchString(name)
}

func validateUserExists(name string) (bool, error) {
	indexData, err := os.ReadFile(os2.ASTRO_VPN_PKI_PATH)
	if err != nil {
		return false, err
	}

	clientExists := strings.Contains(string(indexData), fmt.Sprintf("/CN=%s", name))
	return clientExists, nil
}

func createPassword(userName, password string) error {

	var cmdArgs []string
	var cmd *exec.Cmd

	if len(password) == 0 {
		cmdArgs = []string{"--batch", "build-client-full", userName, "nopass"}
		cmd = exec.Command(os2.ASTRO_VPN_EASYRSA_BIN_PATH, cmdArgs...)
		cmd.Dir = os2.ASTRO_VPN_EASYRSA_PATH

		var stderr bytes.Buffer
		cmd.Stderr = &stderr

		err := cmd.Run()
		if err != nil {
			return fmt.Errorf("%v, stderr: %s", err, stderr.String())
		}
	} else {
		cmdArgs = []string{"--batch", "build-client-full", userName}
		cmd = exec.Command(os2.ASTRO_VPN_EASYRSA_BIN_PATH, cmdArgs...)

		var bufout, buferr bytes.Buffer
		cmd.Dir = os2.ASTRO_VPN_EASYRSA_PATH

		ptmx, err := pty.Start(cmd)
		if err != nil {
			return err
		}
		defer ptmx.Close()

		go func() {
			io.Copy(&bufout, ptmx)
		}()

		go func() {
			io.Copy(&buferr, ptmx)
		}()

		bufin := bufio.NewWriter(ptmx)
		bufin.WriteString(password + "\n")
		bufin.Flush()
		bufin.WriteString(password + "\n")
		bufin.Flush()

		err = cmd.Wait()
		if err != nil {
			return fmt.Errorf("%v, exit: %s, error: %s", err, bufout.String(), buferr.String())
		}
	}

	return nil
}

func getTlsSignature() (string, error) {
	file, err := os.Open(os2.ASTRO_VPN_SERVER_CONF_PATH)
	if err != nil {
		return "", err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "tls-crypt") {
			return "tls-crypt", nil
		} else if strings.HasPrefix(line, "tls-auth") {
			return "tls-auth", nil
		}
	}

	if scanner.Err() != nil {
		return "", fmt.Errorf("%v", scanner.Err())
	}

	return "", fmt.Errorf("No tls-auth or tls-crypt found in server.conf")
}

func generateFile(userName string) error {

	_, err := os.Stat(os2.ASTRO_OVPN_PATH)

	if os.IsNotExist(err) {
		err := os.MkdirAll(os2.ASTRO_OVPN_PATH, 0755)
		if err != nil {
			return err
		}
	}

	template, err := os.ReadFile(os2.ASTRO_VPN_CLIENT_TEMPLATE_PATH)
	if err != nil {
		return err
	}

	ca, err := os.ReadFile(os2.ASTRO_VPN_EASYRSA_CA_PATH)
	if err != nil {
		return err
	}

	cert, err := os.ReadFile(fmt.Sprintf(os2.ASTRO_VPN_EASYRSA_CRT_PATH, userName))
	if err != nil {
		return err
	}

	key, err := os.ReadFile(fmt.Sprintf(os2.ASTRO_VPN_EASYRSA_KEY_PATH, userName))
	if err != nil {
		return err
	}

	tlsSignature, err := getTlsSignature()
	if err != nil {
		return err
	}

	var tlsBlock string

	switch tlsSignature {
	case "tls-crypt":
		tls, err := os.ReadFile(os2.ASTRO_VPN_EASYRSA_CRYPT_PATH)
		if err != nil {
			return err
		}
		tlsBlock = fmt.Sprintf("<tls-crypt>\n%s\n</tls-crypt>", tls)
	case "tls-auth":
		tls, err := os.ReadFile(os2.ASTRO_VPN_EASYRSA_AUTH_PATH)
		if err != nil {
			return err
		}
		tlsBlock = fmt.Sprintf("key-direction 1\n<tls-auth>\n%s\n</tls-auth>", tls)
	}

	ovpnContent := fmt.Sprintf("%s\n<ca>\n%s\n</ca>\n<cert>\n%s\n</cert>\n<key>\n%s\n</key>\n%s", template, ca, cert, key, tlsBlock)
	ovpnFilePath := fmt.Sprintf("%s%s.ovpn", os2.ASTRO_OVPN_PATH, userName)

	if err := os.WriteFile(ovpnFilePath, []byte(ovpnContent), 0644); err != nil {
		return fmt.Errorf("Error writing file %s: %v", ovpnFilePath, err)
	}
	return nil
}

func CreateUser(db *badger.DB, userName, password string) (types.User, error) {
	if !validateClientName(userName) {
		return types.User{}, fmt.Errorf("Invalid client name")
	}
	exists, err := validateUserExists(userName)
	if err != nil {
		return types.User{}, err
	}

	if exists {
		return types.User{}, fmt.Errorf("User already exists")
	}

	err = createPassword(userName, password)
	if err != nil {
		return types.User{}, err
	}
	err = generateFile(userName)
	if err != nil {
		return types.User{}, err
	}

	users, err := ListUsersForPki()
	if err != nil {
		return types.User{}, err
	}

	var user types.User

	for _, u := range users {
		if u.Username == userName {
			user = u
			break
		}
	}

	user.CreatedAt = time.Now()
	user.UsedPassword = true
	if password == "" {
		user.UsedPassword = false
	}

	err = SaveUser(db, user)
	if err != nil {
		return types.User{}, err
	}

	return user, nil
}
