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
1. Starting local docker kafka `make kafka-start`
2. Crete deposit topics `make kafka-create-topics`
3. Running services `make dev`
4. After finished using the app, you can delete the kafka cluster with command `make kafka-stop`