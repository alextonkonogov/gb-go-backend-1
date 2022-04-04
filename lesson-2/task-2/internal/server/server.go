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
	entering       chan client
	leaving        chan client
	messages       chan string
	clients        map[client]bool
	serverColor    *color.Color
	serverColorErr *color.Color
	clientColor    *color.Color
	commandColor   *color.Color
	menuString     string
}

func (s server) Start() error {
	ctx, _ := signal.NotifyContext(context.Background(), os.Interrupt)
	cfg := net.ListenConfig{
		KeepAlive: time.Minute,
	}

	listener, err := cfg.Listen(ctx, "tcp", "localhost:8005")
	if err != nil {
		return err
	}
	fmt.Println("Server is running")
	go s.broadcaster()
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Print(err)
			continue
		}
		go s.handleConn(conn)
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

func (s server) handleConn(conn net.Conn) {
	ch := make(chan string)
	go s.clientWriter(conn, ch)
	who := conn.RemoteAddr().String()

	ch <- s.serverColor.Sprintf("You are %s\n%s", who, s.menuString)
	s.messages <- s.serverColor.Sprintf("%s has arrived", who)
	s.entering <- ch
	input := bufio.NewScanner(conn)
	for input.Scan() {
		msg := strings.TrimSpace(input.Text())

		if msg == "" {
			continue
		}
		if strings.Contains(msg, "--nick set") {
			args := strings.Split(msg, " ")
			if len(args) < 3 {
				ch <- s.serverColorErr.Sprintf("You must enter your nick")
				continue
			}
			nickname := args[2]
			s.messages <- s.serverColor.Sprintf("%s changed nickname to %s", who, nickname)
			who = nickname
			continue
		}
		if msg == "--online get" {
			ch <- s.serverColor.Sprintf("online users %d", len(s.clients))
			continue
		}
		if msg == "--menu get" {
			ch <- s.serverColor.Sprintf("%s", s.menuString)
			continue
		}
		if strings.Contains(msg, "--") {
			ch <- s.serverColorErr.Sprintf("%s is not a command, try again or check the menu "+s.commandColor.Sprintf("--menu get"), msg)
			continue
		}

		s.messages <- s.clientColor.Sprintf("%s:", who) + msg
	}
	s.leaving <- ch
	s.messages <- s.serverColor.Sprintf("%s has left", who)
	conn.Close()
}

func (s server) clientWriter(conn net.Conn, ch <-chan string) {
	for msg := range ch {
		fmt.Fprintln(conn, msg)
	}
}

func NewServer() server {
	blue := color.New(color.FgHiBlue)
	cyan := color.New(color.FgHiCyan)
	return server{
		entering:       make(chan client),
		leaving:        make(chan client),
		messages:       make(chan string),
		clients:        make(map[client]bool),
		serverColor:    blue,
		serverColorErr: color.New(color.FgHiRed),
		clientColor:    color.New(color.FgGreen),
		commandColor:   cyan,
		menuString: "You can use commands:\n" +
			cyan.Sprintf("--online get") + blue.Sprintf(" to get number of online users\n") +
			cyan.Sprintf("--nick set John") + blue.Sprintf(" to change your nickname to John\n") +
			cyan.Sprintf("--menu get") + blue.Sprintf(" to get menu"),
	}
}
