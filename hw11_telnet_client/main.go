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
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	tc := NewTelnetClient(net.JoinHostPort(host, port), timeout, os.Stdin, os.Stdout)
	if err := tc.Connect(); err != nil {
		cancel()
		log.Println(err)
	}
	defer func() {
		if err := tc.Close(); err != nil {
			log.Fatalln(err)
		}
	}()

	go func() {
		select {
		case <-ctx.Done():
			return
		default:
			if err := tc.Send(); err != nil {
				log.Fatalln(err)
			}
			cancel()
			return
		}
	}()

	go func() {
		select {
		case <-ctx.Done():
			return
		default:
			if err := tc.Receive(); err != nil {
				log.Fatalln(err)
			}
			cancel()
			return
		}
	}()

	<-ctx.Done()
}

func checkArgs(host, port string) bool {
	if host == "" || port == "" {
		return false
	}
	if _, err := strconv.Atoi(port); err != nil {
		return false
	}

	if _, err := net.LookupHost(host); err != nil {
		return false
	}

	return true
}
