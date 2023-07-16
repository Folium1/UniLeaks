.SILENT:


build:
ifneq ($(wildcard go.mod),"")
		rm -f go.mod
endif
ifneq ($(wildcard go.sum),"")
		rm -f go.sum
endif
	go mod init leaks
	go mod tidy
	go build cmd/main.go

run:
ifneq ("$(wildcard $(main))","")
	rm main
endif
	docker build -t leaks .
	docker-compose up --build

certs:
	mkdir -p certs

certs/server.crt: certs
	openssl req -new -newkey rsa:2048 -nodes -keyout certs/server.key -out certs/server.csr

certs/server.key: certs/server.crt

.PHONY: generate-certs
generate-certs: certs/server.key
	openssl x509 -req -days 365 -in certs/server.csr -signkey certs/server.key -out certs/server.crt


tests:
	docker run -d -p 3306:3306 --name mysql -e MYSQL_DATABASE=leaks -e MYSQL_ROOT_PASSWORD=root mysql
	go test -v ./... -count=1
	docker stop mysql
	docker rm mysql