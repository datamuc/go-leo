package main

import (
	// "runtime"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"strings"

	"github.com/beevik/etree"
)

func main() {
	search := url.QueryEscape(strings.Join(os.Args[1:], " "))
	resp, err := http.Get(fmt.Sprintf(
            "https://dict.leo.org/dictQuery/m-vocab/ende/query.xml?" +
            "tolerMode=nof&lp=ende&lang=de&rmWords=off&rmSearch=on&searchLoc=0"+
            "&resultOrder=basic&multiwordShowSingle=on&search=%s", search))
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

	/*
	   stdin := os.Stdout
	   if runtime.GOOS != "windows" {
	*/
	cmd := exec.Command("less", "-FX")
	r, stdin := io.Pipe()
	cmd.Stdin = r
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	c := make(chan struct{})
	go func() {
		defer close(c)
		cmd.Run()
	}()
	/*
	   }
	*/

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

		_, max := MinMax(len(dewords), len(enwords))
		for i := 0; i < max; i++ {
			fmt.Fprintf(stdin, "%-35s %-35s\n", Geti(dewords, i), Geti(enwords, i))
		}
	}

	/*
	   if runtime.GOOS != "windows" {
	*/
	stdin.Close()
	<-c
	/*
	   }
	*/
}
