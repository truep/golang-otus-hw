package main

import (
	"errors"
	"io"
	"net"
	"time"
)

const network = "tcp"

var ErrConnactionEstablished = errors.New("there is no established coonection for client")

type TelnetClient interface {
	Connect() error
	io.Closer
	Send() error
	Receive() error
}

type Telnet struct {
	conn    net.Conn
	address string
	timeout time.Duration
	in      io.ReadCloser
	out     io.Writer
}

func (t *Telnet) Connect() error {
	var err error
	t.conn, err = net.DialTimeout(network, t.address, t.timeout)
	return err
}

func (t *Telnet) Close() error {
	if err := t.conn.Close(); err != nil {
		return err
	}

	if err := t.in.Close(); err != nil {
		return err
	}
	return nil
}

func (t *Telnet) Send() error {
	if t.conn == nil {
		return ErrConnactionEstablished
	}
	_, err := io.Copy(t.conn, t.in)
	if err != nil {
		return err
	}
	return nil
}

func (t *Telnet) Receive() error {
	if t.conn == nil {
		return ErrConnactionEstablished
	}

	_, err := io.Copy(t.out, t.conn)
	if err != nil {
		return err
	}
	return nil
}

func NewTelnetClient(address string, timeout time.Duration, in io.ReadCloser, out io.Writer) TelnetClient {
	return &Telnet{
		nil,
		address,
		timeout,
		in,
		out,
	}
}
