SRCS = $(wildcard *.go) $(wildcard */*.go)

all: server

server: $(SRCS)
	go build -o $@

format:
	gofmt -s -w .

clean:
	rm -f server

.PHONY: format clean
