main: build

build:
ifdef TAGS
	go build -tags ${TAGS}
else
	go build
endif

	GOOS="linux" GOARM="6" GOARCH="arm" go build -o onair_server ./server/server.go

debug:
	$(MAKE) build TAGS="mal_debug"
