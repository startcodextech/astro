package types

import os "astro/os/types"

type VPNConfig struct {
	IPs      []os.Iface `json:"ips"`
	PublicIP string     `json:"public_ip"`
	Domains  []string   `json:"domains"`
}
