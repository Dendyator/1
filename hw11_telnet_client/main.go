package main

import (
	"context"
	"flag"
	"log"
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
		log.Fatalf("Usage: go-telnet [--timeout=timeout] host port")
	}

	host, port := flag.Arg(0), flag.Arg(1)
	address := net.JoinHostPort(host, port)

	client := NewTelnetClient(address, timeout, os.Stdin, os.Stdout)
	if client == nil {
		log.Fatalf("Failed to create Telnet client. Check the provided parameters.")
	}

	if err := client.Connect(); err != nil {
		log.Fatalf("Connection failed: %v\n", err)
	}
	log.Printf("Connected to %s\n", address)

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT)
	defer stop()

	go func() {
		<-ctx.Done()
		log.Println("Connection closed by user")
		client.Close()
		os.Exit(0)
	}()

	for {
		if err := client.Send(); err != nil {
			if err.Error() == "EOF" {
				log.Println("EOF received, closing connection")
				break
			}
			log.Println("Connection was closed by peer or error occurred while sending:", err)
			break
		}
		if err := client.Receive(); err != nil {
			log.Println("Connection was closed by peer or error occurred while receiving:", err)
			break
		}
	}

	client.Close()
}
