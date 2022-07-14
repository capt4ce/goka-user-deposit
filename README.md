# Goka User Deposit
This project is a service to handle user deposits.

## Dependencies
1. Go version 1.18.4
2. [Goka](https://github.com/lovoo/goka)

## Installation
1. Clone this repo
2. Run `go mod` to install dependencies


## Generating protobuf model
```bash
protoc -I=pb --go_out=. pb/*.proto
```

## Running the project
1. Starting local docker kafka `docker-compose up -d`
2. Running services `go run main.go`