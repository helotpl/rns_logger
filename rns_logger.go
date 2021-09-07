package main

import (
	"flag"
	"fmt"
	"io/ioutil"
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
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
	}
	return string(body)
}

func main() {
	duration := flag.Int("dur", 3600, "duration of logging")
	interval := flag.Int("i", 2, "interval between downloads")
	names := flag.Bool("names", false, "record only track names")
	flag.Parse()
	onlyTrackNames := *names
	if flag.NArg() != 1 {
		fmt.Println("Please supply non option parameter = output filename")
	}
	fileName := flag.Args()[0]
	file, err := os.Create(fileName)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	ticker := time.NewTicker(time.Duration(*interval) * time.Second)
	defer ticker.Stop()
	done := make(chan bool)
	go func() {
		time.Sleep(time.Duration(*duration) * time.Second)
		done <- true
	}()
	previous := ""
	prevt := time.Now()
	if err != nil {
		fmt.Println(err)
	}
	for {
		select {
		case t := <-ticker.C:
			current := getCurrent()
			if current != previous {
				if previous == "Radio Nowy Świat - Pion i poziom!" && !onlyTrackNames {
					duration := t.Sub(prevt)
					fmt.Fprintln(file, "", duration.Round(time.Second))
				}
				previous = current
				prevt = t
				if current != "Radio Nowy Świat - Pion i poziom!" {
					fmt.Fprintln(file, t.Format(time.Stamp), current)
				} else {
					fmt.Fprint(file, t.Format(time.Stamp))
				}
			}
		case <-done:
			if previous == "Radio Nowy Świat - Pion i poziom!" && !onlyTrackNames {
				duration := time.Now().Sub(prevt)
				fmt.Fprintln(file, "", duration.Round(time.Second))
			}
			return
		}
	}
}
