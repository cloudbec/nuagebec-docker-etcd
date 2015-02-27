package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"time"

	"github.com/coreos/go-etcd/etcd"
)

// http://blog.gopheracademy.com/advent-2013/day-06-service-discovery-with-etcd/

func main() {

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	go func() {
		for sig := range c {
			log.Printf("captured %v, stopping profiler and exiting..", sig)

			os.Exit(1)
		}
	}()
	// used for test
	//	os.Setenv("etcd_discovery", "testxs")

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

	getEtcdCtlMemberAPI("infra5")

	//	curl http://172.17.42.1:4002/v2/members -XPOST -H "Content-Type: application/json" -d '{"peerURLs":["http://10.0.0.10:2380"]}'
	//	postData := `{"name":"infra5", "peerURLs":["http://172.17.42.1:4005"] }

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

func printCommand(cmd *exec.Cmd) {
	fmt.Printf("==> Executing: %s\n", strings.Join(cmd.Args, " "))
}

func printError(err error) {
	if err != nil {
		fmt.Print("error!!!")
		os.Stderr.WriteString(fmt.Sprintf("==> Error: %s\n", err.Error()))
	}
}

func printOutput(outs []byte) {
	if len(outs) > 0 {
		fmt.Printf("==> Output: %s\n", string(outs))
	}
}

func getEtcdMemberAPI() {

	url := "http://172.17.42.1:4002/v2/members"
	httpClient := &http.Client{}
	//req, err := http.NewRequest("POST", url, strings.NewReader(postData))
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatalln(err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := httpClient.Do(req)
	if err != nil {
		log.Fatalln(err)
	}
	defer resp.Body.Close()
	var f interface{}
	body, _ := ioutil.ReadAll(resp.Body)
	err = json.Unmarshal(body, &f)
	m := f.(map[string]interface{})

	for k, v := range m {
		switch vv := v.(type) {
		case string:
			fmt.Println(k, "is string", vv)
		case int:
			fmt.Println(k, "is int", vv)
		case []interface{}:
			fmt.Println(k, "is an array:")
			for i, u := range vv {
				fmt.Println(i, u)
			}
		default:
			fmt.Println(k, "is of a type I don't know how to handle")
		}
	}
	//fmt.Println(string(body))

}

func getEtcdCtlMemberAPI(machineName string) {
	// machineName := "infra5"
	cmd := exec.Command("etcdctl", "-C", "http://172.17.42.1:4002", "member", "add", machineName, "http://172.17.42.1:4005")
	printCommand(cmd)
	printMachineStatus(cmd.CombinedOutput())

}

func (string) printMachineStatus(stdOut []byte, stdErr error) {
	if strings.Contains(string(stdOut), "Added member named "+machineName) {
		printOutput(
			stdOut,
		)

	} else {

		printError(stdErr)

	}
}

type machine struct {
	status string
	name   string
}
