package install

import (
	"astro/dto"
	os "astro/os"
)

func GetConfig() (dto.VPNConfig, error) {

	ips, err := os.GetIps()
	if err != nil {
		return dto.VPNConfig{}, err
	}

	publicIP, err := os.GetPublicIP()
	if err != nil {
		return dto.VPNConfig{}, err
	}

	domains, err := os.GetDomainForIP(publicIP)
	if err != nil {
		return dto.VPNConfig{}, err
	}

	return dto.VPNConfig{
		IPs:      ips,
		PublicIP: publicIP,
		Domains:  domains,
	}, nil
}
