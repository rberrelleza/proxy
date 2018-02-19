# kahlys/proxy

[![godoc](https://godoc.org/github.com/kahlys/proxy?status.svg)](https://godoc.org/github.com/kahlys/proxy) 
[![build](https://api.travis-ci.org/kahlys/proxy.svg?branch=master)](https://travis-ci.org/kahlys/proxy)
[![go report](https://goreportcard.com/badge/github.com/kahlys/proxy)](https://goreportcard.com/report/github.com/kahlys/proxy)

Simple tcp proxy package and executable binary in Golang. The executable provides both TCP and TCP/TLS connection.

## Installation

With a correctly configured [Go toolchain](https://golang.org/doc/install):
```
go get -u github.com/kahlys/proxy/cmd/tcpproxy
```

## Usage

By default, the proxy address is *localhost:4444* and the target address is *localhost:80*.
```
$ tcpproxy
```
You can specify some options.
```
$ tcpproxy -h
Usage of tcpproxy:

  -lhost string
    	proxy local address (default ":4444")

  -lcert string
    	proxy certificate x509 file for tls/ssl use

  -lkey string
    	proxy key x509 file for tls/ssl use
      
  -ltls
    	tls/ssl between client and proxy
      
  -rhots string
    	proxy remote address (default ":80")
      
  -rtls
    	tls/ssl between proxy and target
```
