# Goka User Deposit
This project is a service to handle user deposits.

## Dependencies
1. Go version 1.18.4
2. [Goka](https://github.com/lovoo/goka)

## Installation
1. Clone this repo
2. Starting local docker kafka `make kafka-start`
3. Create deposit topics `make kafka-create-topics`
4. Running services `make start`
5. After finished using the app, you can delete the kafka cluster with command `make kafka-stop`


## Generating protobuf model
```bash
protoc -I=pb --go_out=. pb/*.proto
```

## Running the project
1. Starting local docker kafka `make kafka-start`
2. Crete deposit topics `make kafka-create-topics`
3. Running services `make start`
4. After finished using the app, you can delete the kafka cluster with command `make kafka-stop`

## Testing the project
- sample request for deposit
```
curl --location --request POST 'localhost:8000/deposit' \
--header 'Content-Type: application/json' \
--data-raw '{
    "wallet_id": "a",
    "amount": 7000
}'
```
- sample request for getting wallet balance
```
curl --location --request GET 'localhost:8000/details' \
--header 'Content-Type: application/json' \
--data-raw '{
    "wallet_id": "a"
}'
```

## Useful tools
- [Kafka IDE](https://kafkaide.com/): to query kafka data that would let us visualize the content of our kafka service

## Common issues
- Goka Cache error
```
2022/07/15 13:52:15 [Processor aboveThreshold] setup generation 1, claims=map[string][]int32{"deposits":[]int32{0}}
2022/07/15 13:52:15 1 error occurred:
        * error consuming from group consumer: 1 error occurred:
        * error setting up (partition=0): Setup failed. Cannot start processor for partition 0: 1 error occurred:
        * kafka tells us there's no message in the topic, but our cache has one. The table might be gone. Try to delete your local cache! Topic aboveThreshold-table, partition 0, hwm 0, local offset 12
```
Solution: delete everything in `/tmp/goka` (goka cache location)