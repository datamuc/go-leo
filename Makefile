leo: leo.go
	go build -ldflags -s -o leo

leo.exe: leo.go
	GOOS=windows GOARCH=386 go build -ldflags -s -o leo.exe


.PHONY: all

all: leo leo.exe
