VERSION?="0.0.1"
PROJECT := github.com/towa48/go-libre-storage
GOFMT_FILES?=$$(find . -name '*.go' | grep -v vendor)

tools:
	go get -u github.com/kardianos/govendor

build:
	@mkdir -p ./bin
	GOGC=off go build -i -o ./bin/go-libre-storage ./cmd/go-libre-storage

clean:
	@rm -rf ./bin

vendor-list:
	@govendor list

vendor-update:
	@govendor update +vendor

fmt:
	@govendor fmt +local

# disallow any parallelism (-j) for Make. This is necessary since some
# commands during the build process create temporary files that collide
# under parallel conditions.
.NOTPARALLEL:

.PHONY: bin test build clean tools