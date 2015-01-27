package main

import (
	"errors"
	"io/ioutil"
	"net"
	"os"
)

// http://blog.gopheracademy.com/advent-2013/day-06-service-discovery-with-etcd/

func main() {
	// used for test only
	os.Setenv("ETCD_DISCOVERY", "testXS")

	filename := ".env_etcd"
	if len(os.Getenv("ETCD_DISCOVERY")) > 5 {
		if _, err := os.Stat(filename); os.IsNotExist(err) {
			ioutil.WriteFile(filename, []byte(os.Getenv("ETCD_DISCOVERY")), 0600)
		}
	}
	extIP, _ := externalIP()
	print(extIP)
	// then
	// cat > /data/.env_etcd << EOF
	// $ETCD_DISCOVERY
	// EOF
	// fi

}

func externalIP() (string, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return "", err
	}
	for _, iface := range ifaces {
		if iface.Flags&net.FlagUp == 0 {
			continue // interface down
		}
		if iface.Flags&net.FlagLoopback != 0 {
			continue // loopback interface
		}

		addrs, err := iface.Addrs()
		if err != nil {
			return "", err
		}
		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			if ip == nil || ip.IsLoopback() {
				continue
			}
			ip = ip.To4()
			if ip == nil {
				continue // not an ipv4 address
			}
			return ip.String(), nil
		}
	}
	return "", errors.New("are you connected to the network?")
}
