package client

import (
	"context"
	"fmt"
	"io"
	"net"
	"os"
	"os/signal"
	"time"
)

type client struct{}

func (c client) Start() error {
	ctx, _ := signal.NotifyContext(context.Background(), os.Interrupt)
	d := net.Dialer{
		Timeout:   time.Second,
		KeepAlive: time.Minute,
	}

	conn, err := d.DialContext(ctx, "tcp", "localhost:8002")
	if err != nil {
		return err
	}

	defer conn.Close()

	go func() {
		io.Copy(os.Stdout, conn)
	}()
	io.Copy(conn, os.Stdin) // until you send ^Z
	fmt.Printf("%s: exit", conn.LocalAddr())
	return nil
}

func NewClient() client {
	return client{}
}
