package vpn

import (
	"archive/zip"
	os2 "astro/os"
	"astro/vpn/types"
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func ListUsersForPki() ([]types.User, error) {
	lines, err := getPkiIndexFile()
	if err != nil {
		return []types.User{}, err
	}

	users := make([]types.User, 0)

	for _, line := range lines {
		user, err := parseLine(line)
		if err != nil {
			return []types.User{}, err
		}
		users = append(users, user)
	}

	return users, nil
}

func GeneratePathFileOvpn(username string) (string, error) {

	exists, err := validateUserExists(username)
	if err != nil {
		return "", err
	}

	if !exists {
		return "", fmt.Errorf("User %s does not exists", username)
	}

	path := fmt.Sprintf("%s%s.ovpn", os2.ASTRO_OVPN_PATH, username)

	zipbuf := new(strings.Builder)

	zipwriter := zip.NewWriter(zipbuf)

	file, err := os.Open(path)
	if err != nil {
		return "", err
	}

	defer file.Close()

	info, err := file.Stat()
	if err != nil {
		return "", err
	}

	header, err := zip.FileInfoHeader(info)
	if err != nil {
		return "", err
	}
	header.Name = filepath.Base(path)
	header.Method = zip.Store

	writer, err := zipwriter.CreateHeader(header)
	if err != nil {
		return "", err
	}
	if _, err = io.Copy(writer, file); err != nil {
		return "", err
	}

	if err = zipwriter.Close(); err != nil {
		return "", err
	}

	return zipbuf.String(), nil

}

func getPkiIndexFile() ([]string, error) {
	file, err := os.Open(os2.ASTRO_VPN_PKI_PATH)
	if err != nil {
		return []string{}, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	lines := []string{}

	for scanner.Scan() {
		line := scanner.Text()
		lines = append(lines, line)
	}

	if err := scanner.Err(); err != nil {
		return []string{}, fmt.Errorf("Error reading file: %s", err)
	}

	return lines, nil
}

func parseLine(line string) (types.User, error) {
	fields := strings.Fields(line)

	expire, err := time.Parse("060102150405Z", fields[1])
	if err != nil {
		return types.User{}, err
	}

	if fields[0] == "R" {

		revoke, err := time.Parse("060102150405Z", fields[2])
		if err != nil {
			return types.User{}, err
		}

		return types.User{
			Status:    fields[0],
			RevokedAt: &revoke,
			ExpiredAt: expire,
			Serial:    fields[3],
			Username:  strings.Split(fields[5], "=")[1],
		}, nil
	}

	return types.User{
		Status:    fields[0],
		ExpiredAt: expire,
		Serial:    fields[2],
		Username:  strings.Split(fields[4], "=")[1],
	}, nil
}
