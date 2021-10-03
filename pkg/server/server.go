package server

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"github.com/bernardosecades/feeder/pkg/service"
	"log"
	"net"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

var (
	ErrClientIndicateTerminate = errors.New("client indicate 'terminate'")
)

// Config pending text
type Config struct {
	Protocol  string
	Host      string
	Port      string
	KeepAlive time.Duration
	MaxConn   int
}

type Server interface {
	Start(ctx context.Context) error
}

// Server pending text
type server struct {
	cf     Config
	feeder service.Feeder
	stopCh chan bool // To control input "terminate" and disconnect all clients and perform a clean shutdown.
	connCh chan bool // Semaphore to control max concurrency in connections (with buffered channel).
}

// NewServer create new instance of server with config and service
func NewServer(cf Config, feeder service.Feeder) Server {
	return &server{
		cf:     cf,
		feeder: feeder,
		stopCh: make(chan bool),
		connCh: make(chan bool, cf.MaxConn),
	}
}

// Start start the server and running until detect timeout/cancel signals and 'terminate' message from some client
func (s *server) Start(ctx context.Context) error {
	defer close(s.stopCh)
	defer close(s.connCh)

	ctx, cancelSignal := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM)
	defer cancelSignal()
	ctx, cancelTimeout := context.WithTimeout(ctx, s.cf.KeepAlive)
	defer cancelTimeout()

	fmt.Println("Starting " + s.cf.Protocol + " server on " + s.cf.Host + ":" + s.cf.Port)
	l, err := net.Listen(s.cf.Protocol, s.cf.Host+":"+s.cf.Port)
	if err != nil {
		log.Println(err)
		return err
	}
	defer l.Close()

	go s.connectionsHandler(l, ctx)

	select {
	case <-ctx.Done(): // We detect context done by timeout or cancel signals from the system.
		s.stop()
		return ctx.Err()
	case <-s.stopCh:  // Client send 'terminate' to disconnect all clients and perform a clean shutdown.
		s.stop()
		return ErrClientIndicateTerminate
	}
}

// stop it will be called when server stop (by context=signal, timeout or message 'terminate' from client)
// It will get report and persist that report from that execution.
func (s *server) stop() {

	// Log unique SKUs
	s.feeder.Log()

	// Print report in stdout
	totalUnique, totalDuplicated, totalInvalid := s.feeder.Report()

	log.Println("total number of unique product skus received for this run of the Application:", totalUnique)
	log.Println("total number of duplicated products skus received for this run of the Application:", totalDuplicated)
	log.Println("total number of invalid Feeder format received for this run of the Application:", totalInvalid)

	// Persist unique SKUs in running in storage if already were not inserted
	totalInserted, totalSkipped, err := s.feeder.Persist()
	if err != nil {
		log.Panic(err)
	}

	log.Println("total feeder persisted in storage:", totalInserted)
	log.Println("total feeder skipped to persist in storage:", totalSkipped)
}

// isLimitConnReached it will check if connCh channel is filled
func (s *server) isLimitConnReached() bool {
	return len(s.connCh) == cap(s.connCh)
}

// connectionsHandler it will handle connections to limit number of concurrency connections
func (s *server) connectionsHandler(listener net.Listener, ctx context.Context) {
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("error connection", err)
			return
		}

		if s.isLimitConnReached() {
			log.Println("concurrent connections were reached")
			_, err = conn.Write([]byte("limit connections reached\n"))
			if err != nil {
				log.Panic(err)
			}

			err = conn.Close()
			if err != nil {
				log.Panic(err)
			}

			continue
		}

		// "increment" connection in buffered channel sending a boolean.
		s.connCh <- true

		go s.requestsHandler(conn, ctx)
	}
}

// requestsHandler it will handle the request from client. It will add the sku using the feeder service and
// controle if some client send message 'terminate' to stop the application
func (s *server) requestsHandler(conn net.Conn, ctx context.Context) {
	buf := bufio.NewReader(conn)
	for {
		input, err := buf.ReadString('\n')
		if err != nil {
			log.Println("client disconnected", conn.RemoteAddr().String())
			break
		}

		input = strings.ReplaceAll(input, "\n", "")
		input = strings.ReplaceAll(input, "\r", "")

		if input == "terminate" {
			s.stopCh <- true
		} else {
			s.feeder.AddSku(input)
		}

		_, err = conn.Write([]byte("OK\n"))
		if err != nil {
			log.Panic(err)
		}

		err = conn.Close()
		if err != nil {
			log.Panic(err)
		}
	}

	// We decrement connection in buffered channel getting the boolean
	// (release resource concurrent connections).
	<-s.connCh
}
