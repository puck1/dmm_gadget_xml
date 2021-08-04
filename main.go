package main

import (
	"flag"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net"
	"net/http"
)

type address struct {
	IP 		string
	Port 	int64
}

const gadgetXml =
`<?xml version="1.0" encoding="UTF-8" ?>
<Module>
<ModulePrefs title="Sample Application">
<Link rel="event.addapp" href="http://{{.IP}}:{{.Port}}/addapp" method="POST" />
<Link rel="event.suspendapp" href="http://{{.IP}}:{{.Port}}/suspendapp" method="POST" />
<Link rel="event.resumeapp" href="http://{{.IP}}:{{.Port}}/resumeapp" method="POST" />
<Link rel="event.removeapp" href="http://{{.IP}}:{{.Port}}/removeapp" method="GET" />
</ModulePrefs>
</Module>`

var (
	gadgetXmlTmpl *template.Template
	serverAddr    address

	serverIP = flag.String("i", "", "The IP address of the server")
	serverPort = flag.Int64("p", 0, "The port of the server address")
	ipServiceAddr = flag.String("e", "", "The IP address of external IP Service from where we can get our IP address")
)

////////////////////////////////////////////////////////////////////////////////////////////
func handleIndex(w http.ResponseWriter, r *http.Request) {
	log.Println("handle Index")
	err := gadgetXmlTmpl.Execute(w, serverAddr)
	if err != nil {
		log.Println("error writing gadget xml")
	}
}

func handleAddApp(w http.ResponseWriter, r *http.Request) {
	log.Println("handle addApp")
	if r.Method != http.MethodPost {
		log.Println("invalid http method in addApp")
		w.WriteHeader(http.StatusMethodNotAllowed)
	} else {
		fmt.Fprint(w, "success")
	}
}

func handleSuspendApp(w http.ResponseWriter, r *http.Request) {
	log.Println("handle suspendApp")
	if r.Method != http.MethodPost {
		log.Println("invalid http method in suspendApp")
		w.WriteHeader(http.StatusMethodNotAllowed)
	} else {
		fmt.Fprint(w, "success")
	}
}

func handleResumeApp(w http.ResponseWriter, r *http.Request) {
	log.Println("handle resumeApp")
	if r.Method != http.MethodPost {
		log.Println("invalid http method in resumeApp")
		w.WriteHeader(http.StatusMethodNotAllowed)
	} else {
		fmt.Fprint(w, "success")
	}
}

func handleRemoveApp(w http.ResponseWriter, r *http.Request) {
	log.Println("handle removeApp")
	if r.Method != http.MethodGet {
		log.Println("invalid http method in removeApp")
		w.WriteHeader(http.StatusMethodNotAllowed)
	} else {
		fmt.Fprint(w, "success")
	}
}

////////////////////////////////////////////////////////////////////////////////////////////
func getServerIp() (string ,error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "", err
	}

	for _, addr := range addrs {
		// 检查 ip 地址判断是否回环地址
		if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() && IsPublicIP(ipnet.IP) {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String(), nil
			}
		}
	}

	res, err := http.Get("http://" + *ipServiceAddr)
	if err == nil && res.StatusCode == http.StatusOK {
		defer res.Body.Close()
		body, err := ioutil.ReadAll(res.Body)
		if err == nil {
			if IsPublicIP(body) {
				return string(body), nil
			}
		}
	}

	return "", fmt.Errorf("can not find the server ip address")
}

func IsPublicIP(IP net.IP) bool {
	if IP.IsLoopback() || IP.IsLinkLocalMulticast() || IP.IsLinkLocalUnicast() {
		return false
	}
	if ip4 := IP.To4(); ip4 != nil {
		switch true {
		case ip4[0] == 10:
			return false
		case ip4[0] == 172 && ip4[1] >= 16 && ip4[1] <= 31:
			return false
		case ip4[0] == 192 && ip4[1] == 168:
			return false
		default:
			return true
		}
	}
	return false
}

////////////////////////////////////////////////////////////////////////////////////////////
func main() {
	flag.Parse()

	if *serverIP == "" {
		ip, err := getServerIp()
		if err != nil {
			log.Fatal("error get server IP: ", err)
		}
		*serverIP = ip
	}
	serverAddr.IP = *serverIP

	http.HandleFunc("/", handleIndex)
	http.HandleFunc("/addapp", handleAddApp)
	http.HandleFunc("/suspendapp", handleSuspendApp)
	http.HandleFunc("/resumeapp", handleResumeApp)
	http.HandleFunc("/removeapp", handleRemoveApp)

	listenAddr := ":0"
	if *serverPort != 0 {
		listenAddr = fmt.Sprint(":", *serverPort)
	}

	l, err := net.Listen("tcp", listenAddr)
	if err != nil {
		log.Fatal("error listening: ", err)
	}
	log.Println("listening at: ", l.Addr())
	serverAddr.Port = int64(l.Addr().(*net.TCPAddr).Port)

	gadgetXmlTmpl, err = template.New("gadgetXml").Parse(gadgetXml)
	if err != nil {
		log.Fatal("error paring literal tmpl: ", err)
	}

	log.Fatal(http.Serve(l, nil))
}