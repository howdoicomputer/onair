build:
	go build

	export GOOS="linux"
	export GOARM="6"
	export GOARCH="arm"

	go build -o onair_server ./server/server.go
