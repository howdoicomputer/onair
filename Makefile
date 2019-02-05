# This Makefile operates under the assumption that the Raspberry Pi has been
# configured as a host named 'onair' in ssh config. It also assumes that
# the ssh keys are properly setup.

main: clean build

# Pass in mal_debug as a tag to get malgo debugging information.
.PHONY: build
build:
ifdef TAGS
	go build -tags ${TAGS}
else
	go build
endif

	GOOS="linux" GOARM="6" GOARCH="arm" go build -o onair_server ./server/server.go

	tar cf build/onair.tar.gz scripts/deploy.sh scripts/onair.service onair_server

clean:
	rm -rf onair.exe onair_server build/*

deploy: clean build
	scp build/onair.tar.gz onair:/home/pi
	ssh -t onair "tar xf /home/pi/onair.tar.gz && sudo bash /home/pi/scripts/deploy.sh"
