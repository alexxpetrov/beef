include docker-compose/.env

SCYLLA_HOST=127.0.0.1   # Change to "scylladb" if running inside Docker Compose network
SCYLLA_PORT=9042

ifeq ($(db), pg)
	url := $(DATABASE_URL)
else ifeq ($(db), cassandra)
	url := $(CASSANDRA_URL)
endif

docker-local:
	docker compose -f docker-compose/docker-compose-local.yaml --env-file=docker-compose/.env up --remove-orphans --build

docker-dev:
	docker compose -f docker-compose/docker-compose-dev.yaml --env-file=docker-compose/.env up --remove-orphans --build -d

migrate-create:
	migrate create -ext sql -dir migrations/$(db) -seq $(name)

# Example: > db=pg make migrate-up
migrate-up:
	migrate -path migrations/$(db) -database $(url) up

migrate-down:
	migrate -path migrations/$(db) -database $(url) down

scylla-migrate-up:
	@echo "Running ScyllaDB migrations..."
	@for file in migrations/scylla/*.up.cql; do \
		echo "Applying migration $$file..."; \
		cqlsh $(SCYLLA_HOST) $(SCYLLA_PORT) -f $$file || { echo "Migration $$file failed!" && exit 1; }; \
		echo "Migration $$file applied successfully."; \
	done

scylla-migrate-down:
	@echo "Reverting ScyllaDB migrations..."
	@for file in migrations/scylla/*.down.cql; do \
		echo "Reverting migration $$file..."; \
		cqlsh $(SCYLLA_HOST) $(SCYLLA_PORT) -f $$file || { echo "Migration $$file failed!" && exit 1; }; \
		echo "Migration $$file reverted successfully."; \
	done

create-kafka-topic-local:
	./migrations/kafka/local/create_topic.sh

delete-kafka-topic-local:
	./migrations/kafka/local/delete_topic.sh

create-kafka-topic-dev:
	./migrations/kafka/dev/create_topic.sh

delete-kafka-topic-dev:
	./migrations/kafka/dev/delete_topic.sh