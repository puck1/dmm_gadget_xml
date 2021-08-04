package main

import (
	"fmt"
	"html/template"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"
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
)

////////////////////////////////////////////////////////////////////////////////////////////
func handleIndex(w http.ResponseWriter, r *http.Request) {
	log.Println("handle Index")
	//f, err := os.Open("./gadget.xml")
	//if err != nil {
	//	log.Println("error open gadget xml")
	//	w.WriteHeader(http.StatusInternalServerError)
	//	fmt.Fprint(w, "error open gadget xml: ", err)
	//	return
	//}
	//w.WriteHeader(http.StatusOK)
	//_, err = io.Copy(w, f)
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
		if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() && IsPublicIP(ipnet) {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String(), nil
			}
		}
	}

	return "", fmt.Errorf("can not find the server ip address")
}

func IsPublicIP(IPNet *net.IPNet) bool {
	IP := IPNet.IP
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
	if len(os.Args) > 1 {
		port, err := strconv.ParseInt(os.Args[1], 10, 64)
		if err != nil {
			log.Fatal("err get port from arg: ", err)
		}
		serverAddr.Port = port
	}

	serverIP, err := getServerIp()
	if err != nil {
		log.Fatal("error get server IP: ", err)
	}
	serverAddr.IP = serverIP

	http.HandleFunc("/", handleIndex)
	http.HandleFunc("/addapp", handleAddApp)
	http.HandleFunc("/suspendapp", handleSuspendApp)
	http.HandleFunc("/resumeapp", handleResumeApp)
	http.HandleFunc("/removeapp", handleRemoveApp)

	listenAddr := ":0"
	if serverAddr.Port != 0 {
		listenAddr = fmt.Sprint(":", serverAddr.Port)
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