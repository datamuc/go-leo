package main

import (
	"fmt"
	"github.com/beevik/etree"
	"github.com/gosuri/uitable"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

func main() {
	search := url.QueryEscape(strings.Join(os.Args[1:], " "))
	url_fmt := "https://dict.leo.org/dictQuery/m-vocab/ende/query.xml?" +
		"tolerMode=nof&lp=ende&lang=de&rmWords=off&rmSearch=on&searchLoc=0" +
		"&resultOrder=basic&multiwordShowSingle=on&search=%s"
	resp, err := http.Get(fmt.Sprintf(url_fmt, search))
	if err != nil {
		log.Fatal(err)
	}
	content, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		log.Fatal(err)
	}

	doc := etree.NewDocument()
	if err := doc.ReadFromBytes(content); err != nil {
		panic(err)
	}
	root := doc.SelectElement("xml")

	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("more")
	} else {
		cmd = exec.Command("less", "-FX")
	}
	r, stdin := io.Pipe()
	cmd.Stdin = r
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	c := make(chan struct{})
	go func() {
		defer close(c)
		cmd.Run()
	}()

	MinMax := func(a int, b int) (int, int) {
		if a > b {
			return b, a
		} else {
			return a, b
		}
	}

	Geti := func(s []string, i int) string {
		if i >= len(s) {
			return ""
		} else {
			return s[i]
		}
	}

	table := uitable.New()
	table.MaxColWidth = 35
	table.Wrap = true

	for _, entry := range root.FindElements("//entry") {
		de := entry.FindElement("//side[@lang='de']")
		en := entry.FindElement("//side[@lang='en']")
		dewords := make([]string, 0)
		enwords := make([]string, 0)

		for _, w := range de.FindElements(".//words/word") {
			dewords = append(dewords, w.Text())
		}
		for _, w := range en.FindElements(".//words/word") {
			enwords = append(enwords, w.Text())
		}

		_, max := MinMax(len(dewords), len(enwords))
		for i := 0; i < max; i++ {
			table.AddRow("|"+Geti(dewords, i), Geti(enwords, i))
		}
	}

	fmt.Fprintf(stdin, "%s\n", table)

	stdin.Close()
	<-c
}
