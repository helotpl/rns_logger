package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

func getCurrent() string {
	resp, err := http.Get("https://nowyswiat.online/dev/current.txt")
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
	}
	return string(body)
}

func main() {
	duration := flag.Int("dur", 3600, "duration of logging")
	interval := flag.Int("i", 2, "interval between downloads")
	flag.Parse()
	if flag.NArg() != 1 {
		fmt.Println("Please supply non option parameter = output filename")
	}
	fileName := flag.Args()[0]
	file, err := os.Create(fileName)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	ticker := time.NewTicker(time.Duration(*interval)*time.Second)
	defer ticker.Stop()
	done := make(chan bool)
	go func() {
		time.Sleep(time.Duration(*duration) * time.Second)
		done <- true
	}()
	previous := ""
	for {
		select {
		case t := <-ticker.C:
			current := getCurrent()
			if current != "Radio Nowy Świat - Pion i poziom!" && current != previous {
				previous = current
				fmt.Fprintln(file, t.Format(time.Stamp), current)
			}
		case <-done:
			return
		}
	}
}