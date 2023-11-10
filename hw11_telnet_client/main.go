package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
)

var timeout time.Duration

func init() {
	flag.Usage = func() {
		name := os.Args[0]
		fmt.Fprintf(os.Stderr, "Usage of %s [--timeout=duration] host port:\n", name)
		fmt.Fprintf(os.Stderr, "  host\n\tserver address(ip/domain)\n  port\n\tserver port\n")

		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "Example:\n %s --timeout=20s 127.0.0.1 8080\n %s 127.0.0.1 8080\n", name, name)
	}

	flag.DurationVar(&timeout, "timeout", 10*time.Second, "server connection timeout")
	log.SetFlags(log.Ldate | log.Lmicroseconds | log.Llongfile)
}

func main() {
	flag.Parse()
	fargs := flag.Args()

	if len(fargs) != 2 {
		flag.Usage()
		log.Fatalln("Not enough or unnecessary arguments")
	}

	host, port := fargs[0], fargs[1]

	if ok := checkArgs(host, port); !ok {
		flag.Usage()
		log.Fatalln("Not valid arguments")
	}

	tc := NewTelnetClient(net.JoinHostPort(host, port), timeout, os.Stdin, os.Stdout)
	if err := tc.Connect(); err != nil {
		log.Fatalln(err)
	}
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()
	defer func() {
		if err := tc.Close(); err != nil {
			cancel()
			log.Fatalln(err)
		}
	}()

	go func() {
		defer cancel()
		if err := tc.Send(); err != nil {
			log.Fatalln(err)
		}
	}()

	go func() {
		defer cancel()
		if err := tc.Receive(); err != nil {
			log.Fatalln(err)
		}
	}()

	<-ctx.Done()
}

func checkArgs(host, port string) bool {
	if host == "" || port == "" {
		return false
	}
	if portNum, err := strconv.ParseUint(port, 10, 32); err != nil || portNum > 65535 || portNum < 1 {
		return false
	}

	if _, err := net.LookupHost(host); err != nil {
		return false
	}

	return true
}
