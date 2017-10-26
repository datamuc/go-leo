leo: import leo.go
	go build -ldflags -s -o leo

import: 
	export GOPATH=$(CURDIR)
	go get github.com/beevik/etree

leo.exe: leo.go
	GOOS=windows GOARCH=386 go build -ldflags -s -o leo.exe

clean:	
	rm -f $(CURDIR)/leo
	rm -f $(CURDIR)/leo.exe
	rm -rf $(CURDIR)/pkg

.PHONY: import all clean

all: leo leo.exe
