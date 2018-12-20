all: leo.linux leo.exe leo.freebsd leo.netbsd leo.openbsd leo.darwin

leo.linux: leo.go
	GOOS=linux go build -ldflags -s -o leo.linux

leo.exe: leo.go
	GOOS=windows GOARCH=386 go build -ldflags -s -o leo.exe

leo.freebsd: leo.go
	GOOS=freebsd go build -ldflags -s -o leo.freebsd

leo.linux-arm: leo.go
	GOOS=linux GOARCH=arm go build -ldflags -s -o leo.linux-arm

leo.darwin: leo.go
	GOOS=darwin go build -ldflags -s -o leo.darwin

.PHONY: all

