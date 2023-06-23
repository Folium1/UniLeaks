.SILENT:


build:
	rm go.mod
	rm go.sum
	rm ./main
	go mod init leaks
	go mod tidy
	go build cmd/main.go

run:
ifneq ("$(wildcard $(main))","")
	rm main
	go build cmd/main.go
else
	go build cmd/main.go
endif
	./main

test:
	go test -v ./testing

