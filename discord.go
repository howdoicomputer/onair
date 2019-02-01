package main

import (
	"encoding/json"
	"io/ioutil"
	"os"

	"github.com/bwmarrin/discordgo"
	log "github.com/sirupsen/logrus"
)

// Have discordgo handle any incoming events pushed to us by the Discord API.
func discordEventReceiver() {
	dg, err := discordgo.New(config.Discord.Token)
	if err != nil {
		log.Error("Error creating Discord session: ", err)
		return
	}

	dg.AddHandler(muting)

	err = dg.Open()
	if err != nil {
		log.Error("Error opening Discord session: ", err)
		return
	}

	log.Info("Discord event receiver active.")
	defer dg.Close()
}

// If Discord let's us know that we've been muted or unmuted
// then write the mute state to a file and set discordUnmuted to
// the appropriate value.
//
// Also, you have to be in a channel in order for Discord to emit
// a VoiceStateUpdate event.
func muting(s *discordgo.Session, v *discordgo.VoiceStateUpdate) {
	if v.UserID == s.State.User.ID {
		lkm := newLkm()

		if v.Deaf || v.Mute || v.SelfDeaf || v.SelfMute {
			log.Info("Muted on Discord.")
			discordUnmuted = false
			lkm.write(true)
		} else {
			log.Info("Unmuted on Discord.")
			discordUnmuted = true
			lkm.write(false)
		}
	}

	return
}

// Discord's API is entirely push based. You cannot query for information.
// Therefore, there is no way to tell if you're muted or unmuted before
// a VoiceStateUpdate event is emitted. I.e, if you muted yourself and then quit
// Discord then there would be no way for OnAir to know that you're muted.
//
// Let's keep track of that muted state in a file.
//
// There is no locking or anything. Pretty simple.
type lastKnownMute struct {
	Muted    bool `json:"muted"`
	filePath string
}

func newLkm() *lastKnownMute {
	return &lastKnownMute{
		filePath: stateFilePath,
	}
}

func (l *lastKnownMute) read() bool {
	dat, err := ioutil.ReadFile(l.filePath)
	if err != nil {
		log.WithFields(log.Fields{
			"filePath": l.filePath,
		}).Error("Could not read file: ", err)
	}

	json.Unmarshal(dat, l)
	return l.Muted
}

func (l *lastKnownMute) write(state bool) {
	l.Muted = state

	m, err := json.Marshal(l)
	if err != nil {
		log.Error("Could not marshal last known muted state JSON: ", err)
	}

	err = os.Truncate(l.filePath, 0)
	if err != nil {
		log.WithFields(log.Fields{
			"filePath": l.filePath,
		}).Error("Could not truncate file: ", err)
	}

	err = ioutil.WriteFile(l.filePath, m, 0644)
	if err != nil {
		log.WithFields(log.Fields{
			"filePath": l.filePath,
		}).Error("Could not write file: ", err)
	}
}
