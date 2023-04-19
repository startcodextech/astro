package auth

import (
	"astro/auth/types"
)

func SystemModules() ([]types.ModuleSystem, error) {
	return []types.ModuleSystem{
		types.ModuleSystem{
			Name: "Auth",
			Actions: []string{
				"create_user",
				"create_any_user",
				"update_user",
				"update_any_user",
				"delete_user",
				"delete_any_user",
				"list_users",
				"get_user",
				"get_any_user",
				"change_password",
				"change_any_password",
			},
		},
		types.ModuleSystem{
			Name: "VPN",
			Actions: []string{
				"create_client",
				"create_any_client",
				"remove_client",
				"remove_any_client",
				"list_clients",
				"get_client",
				"get_any_client",
				"install_vpn",
				"uninstall_vpn",
				"download_ovpn",
				"download_any_ovpn",
				"get_config",
				"set_config",
			},
		},
		types.ModuleSystem{
			Name: "System",
			Actions: []string{
				"reboot",
				"shutdown",
				"restart_service_vpn",
				"restart_service_astro",
				"check_port",
			},
		},
	}, nil
}
