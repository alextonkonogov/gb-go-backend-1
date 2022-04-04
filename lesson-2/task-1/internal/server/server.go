package server

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/fatih/color"
)

type client chan<- string

type server struct {
	entering    chan client
	leaving     chan client
	messages    chan string
	clients     map[client]bool
	serverColor *color.Color
}

func (s server) Start() error {
	ctx, _ := signal.NotifyContext(context.Background(), os.Interrupt)
	cfg := net.ListenConfig{
		KeepAlive: time.Minute,
	}

	listener, err := cfg.Listen(ctx, "tcp", "localhost:8002")
	if err != nil {
		return err
	}
	fmt.Println("Server is running")
	go s.broadcaster()
	go s.messenger()
	go s.ticker()
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Print(err)
			continue
		}
		go s.handleConn(ctx, conn)
	}
}

func (s server) broadcaster() {
	for {
		select {
		case msg := <-s.messages:
			for cli := range s.clients {
				cli <- msg
			}
		case cli := <-s.entering:
			s.clients[cli] = true
		case cli := <-s.leaving:
			delete(s.clients, cli)
			close(cli)
		}
	}
}

func (s server) ticker() {
	ticker := time.NewTicker(time.Second)
	go func() {
		for {
			select {
			case t := <-ticker.C:
				s.messages <- fmt.Sprintf("now: %s", t)
			}
		}
	}()
}

func (s server) messenger() {
	input := bufio.NewScanner(os.Stdin)
	for input.Scan() {
		msg := strings.TrimSpace(input.Text())
		if msg == "" {
			continue
		}
		s.messages <- s.serverColor.Sprintf(msg)
	}
}

func (s server) handleConn(ctx context.Context, conn net.Conn) {
	ch := make(chan string)
	go s.clientWriter(conn, ch)
	s.entering <- ch
	for {
		select {
		case <-ctx.Done():
			return
		}
	}
	s.leaving <- ch
	conn.Close()
}

func (s server) clientWriter(conn net.Conn, ch <-chan string) {
	for msg := range ch {
		fmt.Fprintln(conn, msg)
	}
}

func NewServer() server {
	return server{
		entering:    make(chan client),
		leaving:     make(chan client),
		messages:    make(chan string),
		clients:     make(map[client]bool),
		serverColor: color.New(color.FgHiBlue),
	}
}
