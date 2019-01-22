package main

import (
	"bytes"
	"os"
	"time"

	"github.com/gen2brain/malgo"
	log "github.com/sirupsen/logrus"
)

// Continuously grab a duration of audio.
func mikePoll() {
	ticker := time.NewTicker(time.Duration(config.Mike.PollDuration) * time.Second)
	quit := make(chan struct{})

	go func() {
		for {
			select {
			case <-ticker.C:
				mikeCheck()
			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()
}

// Start the microphone and capture a single second of audio. If the captured audio
// is empty then we can assume that the microphone is muted.
func mikeCheck() {
	ctx, err := malgo.InitContext(nil, malgo.ContextConfig{}, func(message string) {
		log.Infof("Log <%v>\n", message)
	})
	if err != nil {
		log.Error("Could not initiate malgo context: ", err)
		os.Exit(1)
	}
	defer func() {
		_ = ctx.Uninit()
		ctx.Free()
	}()

	deviceConfig := malgo.DefaultDeviceConfig()
	deviceConfig.Format = malgo.FormatS16
	deviceConfig.Channels = config.Mike.Channels
	deviceConfig.SampleRate = config.Mike.SampleRate
	deviceConfig.Alsa.NoMMap = config.Mike.NoMMap

	var capturedSampleCount uint32
	var audioReceived bool
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
		log.Error("Could not initiate microphone: ", err)
		os.Exit(1)
	}

	err = device.Start()
	if err != nil {
		log.Error("Could not start microphone: ", err)
		os.Exit(1)
	}

	time.Sleep(time.Second)

	if audioReceived {
		mikeOn = true
	} else {
		mikeOn = false
	}

	device.Uninit()
}
