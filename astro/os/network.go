package os

import (
	"astro/os/types"
	"io"
	"net"
	"net/http"
)

func GetIps() ([]types.Iface, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}

	info := make([]types.Iface, 0)

	gateway := ""

	for _, iface := range ifaces {
		addrs, err := iface.Addrs()
		if err != nil {
			return nil, err
		}

		for _, addr := range addrs {
			var ip net.IP

			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
				gateway = GetGateway(v).String()
			case *net.IPAddr:
				ip = v.IP
			}

			if ip == nil || ip.IsLoopback() || ip.To4() == nil {
				continue
			}

			netIp, ipNet, err := net.ParseCIDR(addr.String())
			if err != nil {
				return nil, err
			}

			info = append(info, types.Iface{
				Interface: iface.Name,
				Gateway:   gateway,
				Network:   netIp.Mask(ipNet.Mask).String(),
				IP:        ip.String(),
				Mask:      net.IP(ipNet.Mask).String(),
				IsPublic:  IsPrivateIP(ip.String()),
			})
		}
	}
	return info, nil
}

func GetPublicIP() (string, error) {
	resp, err := http.Get("https://api.ipify.org")
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

func IsPrivateIP(ip string) bool {
	// Parse la dirección IP
	parsedIP := net.ParseIP(ip)

	// Si la dirección IP no se puede analizar, es inválida
	if parsedIP == nil {
		return false
	}

	// Obtén los bytes de la dirección IP
	ipBytes := parsedIP.To4()

	// Define los rangos de direcciones IP privadas
	privateRanges := [][]byte{
		{10, 0, 0, 0},        // 10.0.0.0/8
		{172, 16, 0, 0},      // 172.16.0.0/12
		{192, 168, 0, 0},     // 192.168.0.0/16
		{169, 254, 0, 0},     // 169.254.0.0/16
		{127, 0, 0, 0},       // 127.0.0.0/8
		{224, 0, 0, 0},       // 224.0.0.0/4
		{240, 0, 0, 0},       // 240.0.0.0/4
		{255, 255, 255, 255}, // 255.255.255.255/32
	}

	// Verificar si la dirección IP se encuentra dentro de un rango privado
	for _, pr := range privateRanges {
		if ipBytes[0] == pr[0] && ipBytes[1] == pr[1] && (ipBytes[2]&^0x1 == pr[2]&^0x1) {
			return true
		}
	}

	return false
}

func GetGateway(ipNet *net.IPNet) net.IP {
	// Obtener las rutas de la tabla de enrutamiento del sistema
	routes, err := net.InterfaceAddrs()
	if err != nil {
		panic(err)
	}

	// Recorrer todas las rutas
	for _, route := range routes {
		// Obtener la dirección IP y la máscara de red de la ruta actual
		routeIP, routeIPNet, err := net.ParseCIDR(route.String())
		if err != nil {
			panic(err)
		}

		// Si la dirección IP actual no es IPv4 o es una dirección de loopback, saltar
		if routeIP.To4() == nil || routeIP.IsLoopback() {
			continue
		}

		// Comprobar si la máscara de red actual y la especificada son iguales
		if routeIPNet.Mask.String() == ipNet.Mask.String() {
			// Comprobar si la dirección IP actual y la especificada pertenecen a la misma red
			if ipNet.Contains(routeIP) {
				// Obtener la dirección IP de la puerta de enlace de la red actual
				return routeIP
			}
		}
	}

	// Si no se encuentra una puerta de enlace, devolver nil
	return nil
}

func GetDomainForIP(ip string) ([]string, error) {
	names, err := net.LookupAddr(ip)
	if err != nil {
		return nil, err
	}

	domain := make([]string, 0)
	for _, name := range names {
		domain = append(domain, name)
	}
	return domain, nil
}

func IsPortAvailable(port string) bool {
	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return false
	}
	defer listener.Close()
	return true
}
