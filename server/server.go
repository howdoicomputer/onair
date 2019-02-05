package main

import (
	"net"
	"net/rpc"
	"os"
	"os/signal"
	"syscall"

	"gobot.io/x/gobot/drivers/gpio"
	"gobot.io/x/gobot/platforms/raspi"

	log "github.com/sirupsen/logrus"
)

// Which port to bind to
var bindPort = "42586"

// Set up gobot adaptor for writing to Raspi pin
var ada = raspi.NewAdaptor()
var pin = gpio.NewDirectPinDriver(ada, "7")

// OnAir An exported OnAir object.
//
// We don't actually care about the
// object itself - only the exported
// Speaking method.
type OnAir int

// Speaking Receiver for a speaking event.
func (o *OnAir) Speaking(speaking bool, ack *bool) error {
	*ack = true

	if speaking {
		pin.DigitalWrite(1)
	} else {
		pin.DigitalWrite(0)
	}

	return nil
}

func main() {
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
	})

	addr, err := net.ResolveTCPAddr("tcp", "0.0.0.0:"+bindPort)
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
