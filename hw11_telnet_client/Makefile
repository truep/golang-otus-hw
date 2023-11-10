app_name = go-telnet
build_flags = -ldflags "-s -w" 

build:
	go mod tidy
	go build $(build_flags) -o $(app_name) .

.PHONY: clean

test: 
	go test -v -count=1 -race -timeout=1m ./...
	./test.sh

clean:
	rm -rf ./$(app_name)

all: test build