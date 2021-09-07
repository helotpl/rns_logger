package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"os"
	"strings"
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

func durationGraph(duration time.Duration) string {
	m := duration.Truncate(time.Minute)
	minutes := int(math.Round(m.Minutes()))
	tenSeconds := int(math.Round(((duration - m).Round(time.Second*10).Seconds() / 10)))
	return strings.Repeat("#", minutes) + strings.Repeat(".", tenSeconds)
}

func main() {
	noTrack := "Radio Nowy Świat - Pion i poziom!"

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
	prevt := time.Now()
	current := getCurrent()
	previous := current
	if current != noTrack {
		fmt.Fprintln(file, time.Now().Format(time.Stamp), current)
	} else if !onlyTrackNames {
		fmt.Fprint(file, time.Now().Format(time.Stamp))
	}
	for {
		select {
		case t := <-ticker.C:
			current = getCurrent()
			if current != previous {
				if previous == noTrack && !onlyTrackNames {
					d := t.Sub(prevt)
					fmt.Fprintln(file, "", d.Round(time.Second), durationGraph(d))
				}
				previous = current
				prevt = t
				if current != noTrack {
					fmt.Fprintln(file, t.Format(time.Stamp), current)
				} else if !onlyTrackNames {
					fmt.Fprint(file, t.Format(time.Stamp))
				}
			}
		case <-done:
			if previous == noTrack && !onlyTrackNames {
				d := time.Now().Sub(prevt)
				fmt.Fprintln(file, "", d.Round(time.Second), durationGraph(d))
			}
			return
		}
	}
}
