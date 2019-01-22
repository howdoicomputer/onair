# About

OnAir is an application that monitors microphone activity and Discord mute/unmute states to determine whether or not I am speaking into my microphone. This is to solve a *very* specific problem but feel free to peruse.

# How it Functions

There are two parts of OnAir: a server and a client.

The client is responsible for:

* Periodically capturing microphone input and checking if that input actually contains any audio data (checking to see if my physical mic is muted).
* Receiving VoiceStateUpdate events from the Discord API to check whether or not I have muted or unmuted myself in the Discord client itself.
* Sending my 'mute/unmuted' state to an OnAir server in accordance with the aforementioned client responsiblities.

The server, which runs on a Raspberry Pi Zero, is setup to receive a "speaking" value. If that value is true, the server will then control the GPIO pins that are hooked up to a power relay that a physical "on air" sign is plugged into. If I'm on air then the sign will receive power and thus light up.

# Installing

I would be super surprised if anyone else got any use out of this thing but *maaaaaybe** someone might want to run this on their personal computer.

## Client Prereqs

* A Windows host. While I don't *think* there is anything that stops this from being run on Linux, I'm only developing this on Windows for now because that is what I'm gaming on.
* Golang over `v1.11` installed.
* GCC installed on said Windows host. This project relies on [malgo](https://github.com/gen2brain/malgo) and that project relies on [cgo](https://github.com/golang/go/wiki/cgo#windows).

If you need to install golang or gcc, then I highly recommend using [scoop](https://github.com/lukesampson/scoop) to do so.

## Building

To build inside PowerShell,

```
$env:GO111MODULE="on"
go mod vendor
go build
go build -o onair_server.exe ./server/server.go
```

## Running

Start the server in one console tab,

```
onair master *% $ .\onair_server.exe
time="2019-01-21T20:00:54-08:00" level=info msg="Server started."
```

In another console tab,

```
onair master *% $ .\onair.exe
time="2019-01-21T20:26:37-08:00" level=info msg="Last known Discord muted/unmuted state: Unmuted."
time="2019-01-21T20:27:39-08:00" level=info msg="Muted on Discord."
```

If you want to run OnAir on boot then you might want to look into [nssm](https://nssm.cc/).

# Why

I'm almost always in a Discord channel when I'm playing video games. The sign helps my girlfriend know when our conversations might be listened to if she walks up and starts talking to me. Also, it'll help *me* keep track of whether or not I'm on air.

---
