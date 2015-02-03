package main

import (
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"strings"
	"time"

	"io/ioutil"

	"github.com/coreos/go-etcd/etcd"
)

// http://blog.gopheracademy.com/advent-2013/day-06-service-discovery-with-etcd/

func main() {
	// used for test
	os.Setenv("etcd_discovery", "testxs")

	filename := ".env_etcd"
	if len(os.Getenv("etcd_discovery")) > 5 {
		if _, err := os.Stat(filename); os.IsNotExist(err) {
			fmt.Printf("no such file or directory: %s", filename)
			return
		}

	}

	// binary, lookErr := exec.LookPath("etcd")
	// if lookErr != nil {
	// 	panic(lookErr)
	// }

	// //	args := []string{"ls", "-a", "-l", "-h"}
	// args := []string{"ls,"}
	// env := os.Environ()

	// execErr := syscall.Exec(binary, args, env)
	// if execErr != nil {
	// 	panic(execErr)
	// }

	client := etcd.NewClient([]string{"http://172.17.42.1:4002"})
	cluster := client.GetCluster()
	fmt.Println(cluster)
	// client

	// curl http://172.17.42.1:4002/v2/members -XPOST -H "Content-Type: application/json" -d '{"peerURLs":["http://10.0.0.10:2380"]}'
	postData := `{"peerURLs":["http://10.0.0.10:2380"]}`
	url := "http://172.17.42.1:4002/v2/members"
	httpClient := &http.Client{}
	req, err := http.NewRequest("POST", url, strings.NewReader(postData))
	if err != nil {
		log.Fatalln(err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := httpClient.Do(req)
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(body))

	Eth0IPv4, err := externalIPFromIf("docker0")
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println(Eth0IPv4)
	// resp, err := client.Get("creds", false, false)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Print("TEST")
	// log.Printf("Current creds: %s: %s\n", resp.Node.Key, resp.Node.Value)
	// for {
	// 	watchChan := make(chan *etcd.Response)
	// 	loopWatch(client, "/creds", watchChan)
	// }

	// then
	// cat > /data/.env_etcd << EOF
	// $ETCD_DISCOVERY
	// EOF
	// fi

}

func loopWatch(client *etcd.Client, key string, watch chan *etcd.Response) {

	time.Sleep(1000)
	go client.Watch(key, 0, false, watch, nil)
	fmt.Println("test")

	log.Println("Waiting for an update...")
	r := <-watch
	log.Printf("Got updated creds: %s: %s\n", r.Node.Key, r.Node.Value)

}

func externalIPFromIf(ipv4Iface string) (string, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return "", err
	}
	for _, iface := range ifaces {
		// fmt.Print("nom de l'interface : \t")
		// fmt.Print(iface.Name)

		if strings.EqualFold(iface.Name, ipv4Iface) {

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
	}
	return "", errors.New("are you connected to the network?")
}
