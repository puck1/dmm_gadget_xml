package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
)

var port = flag.Int64("p", 0, "The port the server should listen at")

func getRemoteIp(w http.ResponseWriter, r *http.Request) {
	// get client ip address
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	log.Println("remote: ", ip)

	// print out the ip address
	fmt.Fprintf(w,ip)
}

func main() {
	flag.Parse()

	http.HandleFunc("/", getRemoteIp)

	l, err := net.Listen("tcp", fmt.Sprint(":", *port))
	if err != nil {
		log.Fatal("error listening: ", err)
	}
	log.Println("listening at: ", l.Addr())

	log.Fatal(http.Serve(l, nil))
}