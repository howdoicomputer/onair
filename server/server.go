package main

import (
	"net"
	"net/rpc"
	"os"
	"os/signal"
	"syscall"

	log "github.com/sirupsen/logrus"
)

// Which port to bind to
var bindPort = "42586"

// OnAir An exported OnAir object.
//
// We don't actually care about the
// object itself - only the exported
// Speaking method.
type OnAir int

// Speaking Receiver for a speaking event.
func (o *OnAir) Speaking(speaking bool, ack *bool) error {
	if speaking {
		log.Info("Speaking detected.")
	}

	return nil
}

func main() {
	addr, err := net.ResolveTCPAddr("tcp", "localhost:"+bindPort)
	if err != nil {
		log.Errorf("Error binding to port %s: %s", bindPort, err)
	}

	listener, err := net.ListenTCP("tcp", addr)
	if err != nil {
		log.Error("Unable to listen for TCP.")
	}

	onAir := new(OnAir)
	rpc.Register(onAir)
	go rpc.Accept(listener)

	log.Info("Server started.")

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc
}
