HELPERS := ./scripts
OUTPATH := ./bin

EXEC_NAME := go-journal

OS := $(shell uname)
GO_SRC := main.go

ifeq ($(OS),Linux)
	# This is required to get a statically linked binary.
	# Doing this on MacOS breaks something with networking.
	GOVARS := $(GOVARS) CGO_ENABLED=0
endif

.PHONY: build build-linux check-clean ci clean dev format install-deps lint start

install-deps:
	go mod download
	go mod tidy

build:
	$(GOVARS) go build \
		-ldflags '-extldflags "-static"' \
		-o $(OUTPATH)/${EXEC_NAME} \
		${GO_SRC}
	
	# sudo rm -fv /usr/local/bin/${EXEC_NAME}
	# sudo cp $(OUTPATH)/${EXEC_NAME} /usr/local/bin/${EXEC_NAME}

build-linux:
	GOOS=linux $(GOVARS) go build \
				-ldflags '-extldflags "-static"' \
				-o $(OUTPATH)/${EXEC_NAME} 

dev: bin/air
	./bin/air -c .air.toml -d

start: ci
	$(OUTPATH)/${EXEC_NAME}

ci: install-deps build 

check-clean:
	# Ensures working dir is clean
	${HELPERS}/check-clean

format:
	${HELPERS}/format

lint: bin/golangci-lint
	${HELPERS}/lint

clean:
	# Remove files and directories ignored by .gitignore files
	git clean -fdX artifacts bin

docs:
	npm list -g @mermaid-js/mermaid-cli >/dev/null || npm install -g @mermaid-js/mermaid-cli
	mmdc -i docs/diagram.md -o docs/diagram.svg

bin/golangci-lint:
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s v1.52.2

bin/air:
	curl -sSfL https://raw.githubusercontent.com/cosmtrek/air/master/install.sh | sh -s
