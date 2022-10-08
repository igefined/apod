run-local: clean build-swagger local-db-init
	go build -o dist/betera cmd/http/main.go && ./dist/betera

clean:
	rm -rf build
	mkdir -p build

local-db-init:
	docker compose -f ./local/docker-compose.yml up -d

build-swagger:
	go install github.com/swaggo/swag/cmd/swag@v1.7.8
	swag init -g cmd/http/main.go

build-local: local-db-init
	docker build -t betera . && \
	docker run -d --rm -p 8080:8080 --name betera --network=local_betera_db_network betera