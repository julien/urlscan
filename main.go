package main

import (
	"bufio"
	"fmt"
	"net/http"
	"os"
	"regexp"
	"time"
)

var validURL = regexp.MustCompile(`(ftp|http|https):\/\/(\w+:{0,1}\w*@)?(\S+)(:[0-9]+)?(\/|\/([\w#!:.?+=&%@!\-\/]))?`)

type job struct {
	url string // the URL to check
}

func (j job) Execute() error {

	var c = &http.Client{
		Timeout: time.Second * 10,
	}

	resp, err := c.Get(j.url)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d\n", resp.StatusCode)
	}

	return nil
}

type worker struct {
	id int
}

func (w worker) process(j job) {
	if err := j.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %s\n", err)
	} else {
		fmt.Printf("%s is OK\n", j.url)
	}
}

func main() {
	jobCh := make(chan job, 10)

	for i := 0; i < 10; i++ {
		w := worker{i}
		go func(w worker) {
			for j := range jobCh {
				w.process(j)
			}
		}(w)
	}

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		if !validURL.MatchString(scanner.Text()) {
			fmt.Fprintf(os.Stderr, "That doesn't look like a URL...ignoring\n")
			continue
		}

		j := job{url: scanner.Text()}
		go func() {
			jobCh <- j
		}()

	}
}
