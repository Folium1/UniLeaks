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
ifeq ("$(wildcard $(main))","")
	rm main
	go build cmd/main.go
else
	go build cmd/main.go
endif
	./main

certs:
	mkdir -p certs

certs/server.crt: certs
	openssl req -new -newkey rsa:2048 -nodes -keyout certs/server.key -out certs/server.csr

certs/server.key: certs/server.crt

.PHONY: generate-certs
generate-certs: certs/server.key
	openssl x509 -req -days 365 -in certs/server.csr -signkey certs/server.key -out certs/server.crt


test:
	go test -v ./testing -count=1