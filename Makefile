VERSION?="0.0.1"
PROJECT := github.com/towa48/go-libre-storage
GOFMT_FILES?=$$(find . -name '*.go' | grep -v vendor)
DEPLOY_DIR = ./dist
DEPLOY_FILES = ./bin ./configs/ ./views ./web ./LICENSE ./README.md

tools:
	npm i

build-dev-backend:
	@mkdir -p ./bin
	go build -o ./bin/go-libre-storage ./cmd/go-libre-storage

build-dev-frontend:
	npm run build

build-dev: build-dev-backend build-dev-frontend

build-arm7hf:
	@mkdir -p ./bin
	GOGC=40 GOOS=linux GOARCH=arm GOARM=7 CGO_ENABLED=1 CC=arm-linux-gnueabihf-gcc CXX=arm-linux-gnueabihf-g++ go build -ldflags "-linkmode external -extldflags -static" -i -o ./bin/go-libre-storage ./cmd/go-libre-storage
	npm run build

build: build-dev

deploy:
	@mkdir -p $(DEPLOY_DIR)
	@cp -f -r $(DEPLOY_FILES) ./$(DEPLOY_DIR)

clean:
	@rm -rf ./bin $(DEPLOY_DIR)

# disallow any parallelism (-j) for Make. This is necessary since some
# commands during the build process create temporary files that collide
# under parallel conditions.
.NOTPARALLEL:

.PHONY: bin test build clean tools