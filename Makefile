kafka-start:
	docker-compose up -d

kafka-create-topics:
	docker-compose exec broker kafka-topics --create --bootstrap-server localhost:9092 --replication-factor 1 --partitions 1 --topic deposits
	docker-compose exec broker kafka-topics --create --bootstrap-server localhost:9092 --replication-factor 1 --partitions 1 --topic balance-table
	docker-compose exec broker kafka-topics --create --bootstrap-server localhost:9092 --replication-factor 1 --partitions 1 --topic aboveThreshold-table

kafka-stop:
	docker-compose down

dep:
	go mod tidy
	go mod vendor

start:
	go run main.go