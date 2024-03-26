B=$(shell git rev-parse --abbrev-ref HEAD)
BRANCH=$(subst /,-,$(B))
GITREV=$(shell git describe --abbrev=7 --always --tags)
REV=$(GITREV)-$(BRANCH)-$(shell date +%Y%m%d-%H:%M:%S)


.PHONY: dockerx
dockerx:
	docker buildx build --progress=plain --platform linux/amd64,linux/arm/v7,linux/arm64 --no-cache -t jtrw/micro-pinger:latest --push .

build: info
	- cd app && CGO_ENABLED=0 go build -ldflags "-X main.revision=$(REV) -s -w" -o ../dist/micro-pinger

info:
	- @echo "revision $(REV)"
