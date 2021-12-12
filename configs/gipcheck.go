package main

import (
	"crypto/tls"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"syscall"
	"time"
)

var (
	ipv4           bool
	ipv6           bool
	detectProvider string
)

func param() {
	flag.BoolVar(&ipv4, "4", false, "ipv4")
	flag.BoolVar(&ipv6, "6", false, "ipv6")
	flag.StringVar(&detectProvider, "url", "https://ifconfig.io/ip", "detect url")
	flag.Parse()
}

func main() {
	param()
	tr := &http.Transport{
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
			DualStack: false,
			Control: func(network, address string, c syscall.RawConn) error {
				if network == "tcp4" && ipv6 {
					return errors.New("you should not use ipv4")
				}
				return nil
			},
		}).DialContext,
		TLSClientConfig: &tls.Config{ServerName: "ifconfig.io"},
	}
	client := &http.Client{
		Transport: tr,
	}

	request, _ := http.NewRequest("GET", detectProvider, nil)
	resp, err := client.Do(request)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	contents, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Print(resp.StatusCode, " ", string(contents))
}
