package main

type mikeConfig struct {
	// How long between microphone activations.
	PollDuration int

	// Malgo stuff.
	Channels   uint32
	SampleRate uint32
	NoMMap     uint32
	BufferSize uint32
}

type onAirConfig struct {
	// The address to the OnAir RPC server.
	RPCServerAddr string

	// How long to poll for changes to
	// mikeOn and discordUnmuted.
	ActivityPollDuration int
}

type discordConfig struct {
	// The Discord authentication token.
	Token string
}

type configuration struct {
	Discord discordConfig
	Mike    mikeConfig
	OnAir   onAirConfig
}
