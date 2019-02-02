package main

import (
	"bytes"
	"os"
	"time"

	"github.com/gen2brain/malgo"
	log "github.com/sirupsen/logrus"
)

var mikeCtx malgo.AllocatedContext
var audioReceived bool

// Continuously grab a duration of audio.
func mikePoll(device *malgo.Device) {
	ticker := time.NewTicker(time.Duration(config.Mike.PollDuration) * time.Second)
	quit := make(chan struct{})

	go func() {
		for {
			select {
			case <-ticker.C:
				mikeCheck(device)
			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()
}

func initMikeContext() *malgo.AllocatedContext {
	// For some reason, when I use the Direct Sound backend... microphone activity is picked up regardless.
	// backend := []malgo.Backend{malgo.BackendDsound}
	ctx, err := malgo.InitContext(nil, malgo.ContextConfig{}, func(message string) {
		log.Infof("Log <%v>\n", message)
	})
	if err != nil {
		log.Error("Could not initiate malgo context: ", err)
		os.Exit(1)
	}

	return ctx
}

func initMikeDevice(ctx *malgo.AllocatedContext) *malgo.Device {
	deviceConfig := malgo.DefaultDeviceConfig()
	deviceConfig.Format = malgo.FormatS16
	deviceConfig.Channels = config.Mike.Channels
	deviceConfig.SampleRate = config.Mike.SampleRate
	deviceConfig.Alsa.NoMMap = config.Mike.NoMMap
	deviceConfig.BufferSizeInMilliseconds = config.Mike.BufferSize

	var capturedSampleCount uint32
	pCapturedSamples := make([]byte, 0)

	sizeInBytes := uint32(malgo.SampleSizeInBytes(deviceConfig.Format))
	onRecvFrames := func(framecount uint32, pSamples []byte) {
		sampleCount := framecount * deviceConfig.Channels * sizeInBytes

		newCapturedSampleCount := capturedSampleCount + sampleCount

		pCapturedSamples = append(pCapturedSamples, pSamples...)
		emptyByte := make([]byte, len(pSamples))

		if bytes.Equal(pSamples, emptyByte) {
			audioReceived = false
		} else {
			audioReceived = true
		}

		capturedSampleCount = newCapturedSampleCount
	}

	captureCallbacks := malgo.DeviceCallbacks{
		Recv: onRecvFrames,
	}
	device, err := malgo.InitDevice(ctx.Context, malgo.Capture, nil, deviceConfig, captureCallbacks)
	if err != nil {
		log.Error("Could not initialize microphone: ", err)
		os.Exit(1)
	}

	return device
}

// Start the microphone and capture a single second of audio. If the captured audio
// is empty then we can assume that the microphone is muted.
func mikeCheck(device *malgo.Device) {
	err := device.Start()
	if err != nil {
		log.Error("Could not start microphone: ", err)
		os.Exit(1)
	}

	time.Sleep(time.Second)

	lastMikeStatus := mikeOn
	if audioReceived {
		mikeOn = true

		if lastMikeStatus != mikeOn {
			log.Info("Mic is unmuted.")
		}
	} else {
		mikeOn = false

		if lastMikeStatus != mikeOn {
			log.Info("Mic is muted.")
		}
	}

	device.Stop()
}
