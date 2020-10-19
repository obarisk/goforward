.PHONY: build

build: goforward goforward.exe

goforward: main.go
	go build -o goforward .

goforward.exe: main.go
	GOOS=windows GOARCH=amd64 go builf -o goforward.exe .
