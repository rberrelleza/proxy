package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"net"
	"os"

	"github.com/kahlys/proxy"
)

var (
	localAddr  = flag.String("lhost", ":4444", "proxy local address")
	remoteAddr = flag.String("rhost", ":80", "proxy remote address")
	localTLS   = flag.Bool("ltls", false, "tls/ssl between client and proxy")
	localCert  = flag.String("lcert", "", "proxy certificate x509 file for tls/ssl use")
	localKey   = flag.String("lkey", "", "proxy key x509 file for tls/ssl use")
	remoteTLS  = flag.Bool("rtls", false, "tls/ssl between proxy and target")
)

func main() {
	flag.Parse()

	laddr, err := net.ResolveTCPAddr("tcp", *localAddr)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	raddr, err := net.ResolveTCPAddr("tcp", *remoteAddr)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	if *localTLS && !exists(*localCert) && !exists(*localKey) {
		fmt.Println("certificate and key file required")
		os.Exit(1)
	}

	var p = new(proxy.Server)
	if *remoteTLS {
		// Testing only. You needs to specify config.ServerName insteand of InsecureSkipVerify
		p = proxy.NewServer(raddr, nil, &tls.Config{InsecureSkipVerify: true})
	} else {
		p = proxy.NewServer(raddr, nil, nil)
	}

	fmt.Println("Proxying from " + laddr.String() + " to " + p.Target.String())
	if *localTLS {
		p.ListenAndServeTLS(laddr, *localCert, *localKey)
	} else {
		p.ListenAndServe(laddr)
	}
}

func exists(filename string) bool {
	_, err := os.Stat(filename)
	return !os.IsNotExist(err)
}
