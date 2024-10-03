package main

import (
	"context"
	"flag"
	"fmt"
	_ "io"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	var timeout time.Duration
	flag.DurationVar(&timeout, "timeout", 10*time.Second, "timeout for connection")
	flag.Parse()

	if len(flag.Args()) != 2 {
		fmt.Println("Usage: go-telnet [--timeout=timeout] host port")
		return
	}

	host, port := flag.Arg(0), flag.Arg(1)
	address := net.JoinHostPort(host, port)

	client := NewTelnetClient(address, timeout, os.Stdin, os.Stdout)

	if err := client.Connect(); err != nil {
		fmt.Fprintf(os.Stderr, "...Connection failed: %v\n", err)
		return
	}
	fmt.Fprintf(os.Stderr, "...Connected to %s\n", address)

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT)
	defer stop()

	go func() {
		<-ctx.Done()
		fmt.Fprintf(os.Stderr, "\n...Connection closed by user\n")
		client.Close()
		os.Exit(0)
	}()

	for {
		if err := client.Send(); err != nil {
			fmt.Fprintln(os.Stderr, "...Connection was closed by peer or error occurred while sending:", err)
			break
		}
		if err := client.Receive(); err != nil {
			fmt.Fprintln(os.Stderr, "...Connection was closed by peer or error occurred while receiving:", err)
			break
		}
	}

	client.Close()
}
