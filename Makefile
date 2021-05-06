.PHONY: all 
all:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o ./bin/eth-scanner ./cmd/eth-scanner