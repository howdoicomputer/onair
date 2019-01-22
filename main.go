package main

import (
	"net/rpc"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/mitchellh/go-homedir"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var mikeOn bool
var discordUnmuted bool

var stateFilePath string
var configFilePath string
var dataDirPath string

var config configuration

func activityTimer(client *rpc.Client) {
	ticker := time.NewTicker(time.Duration(config.OnAir.ActivityPollDuration) * time.Second)
	quit := make(chan struct{})

	go func() {
		for {
			select {
			case <-ticker.C:
				if mikeOn && discordUnmuted {
					var reply bool

					err := client.Call("OnAir.Speaking", true, &reply)
					if err != nil {
						log.Fatal("Could not call Speaking function on OnAir RPC server: ", err)
					}
				}
			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()
}

// Let's make sure that we have a place to put our state file and our configuration file.
func init() {
	home, err := homedir.Dir()
	if err != nil {
		log.Error("Error locating home directory: ", err)
	}

	dataDirPath, err = filepath.Abs(home + "/AppData/Local/OnAir")
	if err != nil {
		log.Error("Error determining absolute path to OnAir AppData directory: ", err)
	}

	stateFilePath = filepath.Join(dataDirPath + "/last_known_mute")

	err = os.MkdirAll(dataDirPath, os.ModePerm)
	if err != nil {
		log.Error("Could not make directory: ", err)
	}

	configFilePath = filepath.Join(dataDirPath + "/config.toml")

	viper.SetConfigName("config")
	viper.AddConfigPath(dataDirPath)

	err = viper.ReadInConfig()
	if err != nil {
		log.WithFields(log.Fields{
			"configPath": configFilePath,
		}).Panic("Could not read config file: ", err)
	}

	err = viper.Unmarshal(&config)
	if err != nil {
		log.Panic("Unable to decode configuration: ", err)
	}
}

func main() {
	// Determine the last known 'mute/unmute' status for Discord.
	lkm := newLkm()
	discordUnmuted = !lkm.read()

	var lkmState string
	if discordUnmuted {
		lkmState = "Unmuted."
	} else {
		lkmState = "Muted."
	}

	log.Info("Last known Discord muted/unmuted state: " + lkmState)

	// Establish a connection to the OnAir RPC server.
	client, err := rpc.Dial("tcp", config.OnAir.RPCServerAddr)
	if err != nil {
		log.Fatal("Could not connect to OnAir RPC server: ", err)
	}

	// Start polling for microphone activity.
	mikePoll()

	// Start waiting for Discord mute/unmute events.
	discordEventReceiver()

	// Continuously send our 'on air' status to the OnAir RPC server.
	activityTimer(client)

	// Wait for some sort of CTRL+C like exit.
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc
}
