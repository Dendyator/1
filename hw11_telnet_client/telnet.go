package main

import (
	"errors"
	"io"
	"net"
	"time"
)

type TelnetClient interface {
	Connect() error
	io.Closer
	Send() error
	Receive() error
}

type telnetClient struct {
	conn    net.Conn
	address string
	in      io.ReadCloser
	out     io.Writer
	timeout time.Duration
}

func NewTelnetClient(address string, timeout time.Duration, in io.ReadCloser, out io.Writer) TelnetClient {
	if address == "" || timeout <= 0 || in == nil || out == nil {
		return nil
	}
	return &telnetClient{address: address, in: in, out: out, timeout: timeout}
}

func (c *telnetClient) Connect() error {
	var err error
	c.conn, err = net.DialTimeout("tcp", c.address, c.timeout)
	return err
}

func (c *telnetClient) Send() error {
	buf := make([]byte, 1024)
	n, err := c.in.Read(buf)
	if err != nil && !errors.Is(err, io.EOF) {
		return err
	}
	if n > 0 {
		_, err = c.conn.Write(buf[:n])
	}
	return err
}

func (c *telnetClient) Receive() error {
	buf := make([]byte, 1024)
	n, err := c.conn.Read(buf)
	if err != nil {
		return err
	}
	if n > 0 {
		_, err = c.out.Write(buf[:n])
	}
	return err
}

func (c *telnetClient) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}
