package main

import (
	"encoding/xml"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

var (
	plexURL          = ""
	clientIdentifier = ""
	token            = ""
	plexStatusAPI    = "/status/sessions"
)

type MediaContainer struct {
	Videos []Video `xml:"Video"`
}

type Video struct {
	User XMLUser `xml:"User"`
}

type XMLUser struct {
}

func main() {
	flag.StringVar(&plexURL, "h", "http://plex.lan:32400", "Plex server url")
	flag.StringVar(&clientIdentifier, "ci", "", "Plex client identifier")
	flag.StringVar(&token, "t", "", "Plex auth token")
	flag.Parse()

	if plexURL == "" {
		log.Fatal("Please specify host using -h=")
	}

	if token == "" {
		log.Fatal("Please specify token using -t=")
	}

	client := &http.Client{}

	req, err := http.NewRequest("GET", plexURL+plexStatusAPI, nil)
	if err != nil {
		log.Fatalf("Error dialing: %v", err)
	}

	req.Header.Add("X-Plex-Client-Identifier", clientIdentifier)
	req.Header.Add("X-Plex-Product", "Telegraf")
	req.Header.Add("X-Plex-Version", "1")
	req.Header.Add("X-Plex-Token", token)
	response, err := client.Do(req)

	if err != nil {
		return
	}
	defer response.Body.Close()

	data, _ := ioutil.ReadAll(response.Body)
	var media MediaContainer
	if xmlerr := xml.Unmarshal(data, &media); xmlerr != nil {
		log.Fatal("Unable to parse xml data", xmlerr)
	}

	fmt.Printf("plex,server=plex.lan stream=%d\n", len(media.Videos))
}
