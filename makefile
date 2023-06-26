.SILENT:


build:
ifeq ($(wildcard go.mod),"")
		rm -f go.mod
endif
ifeq ($(wildcard go.sum),"")
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

test:
	go test -v ./testing -count=1
