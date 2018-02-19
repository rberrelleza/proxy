package proxy

import (
	"crypto/tls"
	"fmt"
	"net"
	"os"
	"sync"
	"time"
)

// Server is a TCP server that takes an incoming request and sends it to another
// server, proxying the response back to the client.
type Server struct {
	// Target address
	Target *net.TCPAddr

	// Local address
	Addr *net.TCPAddr

	// Director must be a function which modifies the request into a new request
	// to be sent. Its response is then copied back to the client unmodified.
	Director func(b *[]byte)

	// If config is not nil, the proxy connects to the target address and then
	// initiates a TLS handshake.
	Config *tls.Config

	// Timeout is the duration the proxy is staying alive without activity from
	// both client and target. Also, if a pipe is closed, the proxy waits 'timeout'
	// seconds before closing the other one. By default timeout is 60 seconds.
	Timeout time.Duration
}

// NewServer created a new proxy which sends all packet to target. The function dir
// intercept and can change the packet before sending it to the target.
func NewServer(target *net.TCPAddr, dir func(*[]byte), config *tls.Config) *Server {
	p := &Server{
		Target:   target,
		Director: dir,
		Config:   config,
	}
	return p
}

// ListenAndServe listens on the TCP network address laddr and then handle packets
// on incoming connections.
func (p *Server) ListenAndServe(laddr *net.TCPAddr) {
	p.Addr = laddr

	var listener net.Listener
	listener, err := net.ListenTCP("tcp", laddr)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	p.serve(listener)
}

// ListenAndServeTLS acts identically to ListenAndServe, except that it uses TLS
// protocol. Additionally, files containing a certificate and matching private key
// for the server must be provided.
func (p *Server) ListenAndServeTLS(laddr *net.TCPAddr, certFile, keyFile string) {
	p.Addr = laddr

	var listener net.Listener
	cer, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		fmt.Println(err)
		return
	}
	config := &tls.Config{Certificates: []tls.Certificate{cer}}
	listener, err = tls.Listen("tcp", laddr.String(), config)
	if err != nil {
		fmt.Println(err)
		return
	}

	p.serve(listener)
}

func (p *Server) serve(ln net.Listener) {
	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println(err)
			continue
		}

		go p.handleConn(conn)
	}
}

// handleConn handles connection.
func (p *Server) handleConn(conn net.Conn) {
	// connects to target server
	var rconn net.Conn
	var err error
	if p.Config == nil {
		rconn, err = net.Dial("tcp", p.Target.String())
	} else {
		rconn, err = tls.Dial("tcp", p.Target.String(), p.Config)
	}
	if err != nil {
		fmt.Println(err)
		return
	}

	var pipeMux sync.Mutex
	var pipeDone = false

	// write to dst what it reads from src
	var pipe = func(src, dst net.Conn, filter func(b *[]byte)) {
		defer func() {
			pipeMux.Lock()
			// if first pipe to end, closing conn will end the other pipe.
			if !pipeDone {
				conn.Close()
				rconn.Close()
			}
			pipeDone = true
			pipeMux.Unlock()
		}()

		buff := make([]byte, 65535)
		for {
			src.SetReadDeadline(time.Now().Add(10 * time.Second))
			n, err := src.Read(buff)
			if err, ok := err.(net.Error); ok && err.Timeout() {
				continue
			}
			if err != nil {
				return
			}
			b := buff[:n]

			if filter != nil {
				filter(&b)
			}

			_, err = dst.Write(b)
			if err != nil {
				return
			}
		}
	}

	go pipe(conn, rconn, p.Director)
	go pipe(rconn, conn, nil)
}
