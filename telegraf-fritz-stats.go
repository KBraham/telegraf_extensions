package main

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/xml"
	"flag"
	"fmt"
	"log"
	"strconv"
	"strings"

	"golang.org/x/text/encoding/unicode"
	"golang.org/x/text/transform"
)

var (
	baseURL           = ""
	noSessionID       = "0000000000000000"
	loginURL          = "/login_sid.lua"
	homeautoswitchURL = "/webservices/homeautoswitch.lua?sid=%s&switchcmd=%s"
)

type XMLSessionInfo struct {
	SID       string `xml:"SID"`
	Challenge string `xml:"Challenge"`
	BlockTime string `xml:"BlockTime"`
	Rights    string `xml:"Rights"`
}

// Example usage /usr/local/bin/telegraf-fritz-stats -h "http://fritz.box" -p="" --dev "087610161840"
func main() {
	flag.Var(&myDevices, "dev", "Descy")

	passPtr := flag.String("p", "", "Fritz.box password for authentication")
	flag.StringVar(&baseURL, "h", "http://fritz.box", "Fritz.box host url. Default http://fritz.box")

	flag.Parse()
	if *passPtr == "" {
		log.Fatal("Please specify password using -p=")
	}

	sid := loginFritz(*passPtr)
	//fmt.Printf("\n%s\n", sid)
	// fmt.Printf("\n%s\n", homeautoswitch(sid, "getswitchlist", ""))

	for _, ain := range myDevices {
		// fmt.Printf("name\n%s\n", homeautoswitch(sid, "getswitchname", ain))
		// fmt.Printf("state\n%s\n", homeautoswitch(sid, "getswitchstate", ain))
		// fmt.Printf("present\n%s\n", homeautoswitch(sid, "getswitchpresent", ain))
		powerS, err := homeautoswitch(sid, "getswitchpower", ain)
		if err != nil {
			log.Fatalf("Couldn't download powert :: %q", err)
		}

		energyS, err := homeautoswitch(sid, "getswitchenergy", ain)
		if err != nil {
			log.Fatalf("Couldn't download energy :: %q", err)
		}

		powerS = strings.TrimSpace(powerS)
		energyS = strings.TrimSpace(energyS)

		power, err := strconv.ParseInt(powerS, 0, 64)
		if err != nil {
			log.Fatal(err)
		}
		energy, _ := strconv.ParseInt(energyS, 0, 64)
		// temperature := homeautoswitch(sid, "getswitchtemperature", ain)

		fmt.Printf("fritz,ain=%s power=%d,energy=%d\n", ain, power, energy)
	}

	// /net/home_auto_query.lua
	// Does not provide expected results
	//consumptionS, err :=
	//url := baseURL + "/net/home_auto_query.lua?sid=" + sid + "&command=EnergyStats_10&id="+ain+"&xhr=1"
	//if ain != "" {
	//	url += "&ain=" + ain
	//}
	//
	//consumptionS, err := fetchURL(url)
	//fmt.Printf("Consump %s\n", consumptionS)
}

func loginFritz(pass string) string {
	// Test if login is needed
	data, err := FetchURL(baseURL + loginURL)
	if err != nil {
		log.Fatal(err)
	}

	var session XMLSessionInfo
	if err := xml.Unmarshal([]byte(data), &session); err != nil {
		log.Fatal(err)
	}

	if session.SID == noSessionID {
		// Compute response
		challengeResponse := calcResponse(session.Challenge, pass)

		// Submit reponse and fetch sid
		data, err := FetchURL(baseURL + loginURL + "?response=" + challengeResponse)
		if err != nil {
			log.Fatal(err)
		}

		if err := xml.Unmarshal([]byte(data), &session); err != nil {
			log.Fatal(err)
		}
	}

	return session.SID
}

func homeautoswitch(sid string, cmd string, ain string) (string, error) {
	url := baseURL + fmt.Sprintf(homeautoswitchURL, sid, cmd)
	if ain != "" {
		url += "&ain=" + ain
	}

	return FetchURL(url)
}

func calcResponse(token string, pass string) string {
	return token + "-" + utf16leMd5(token+"-"+pass)
}

func utf16leMd5(s string) string {
	enc := unicode.UTF16(unicode.LittleEndian, unicode.IgnoreBOM).NewEncoder()
	hasher := md5.New()
	t := transform.NewWriter(hasher, enc)
	t.Write([]byte(s))
	return hex.EncodeToString(hasher.Sum(nil))
}

type arrayDevices []string

func (i *arrayDevices) String() string {
	return "my string representation"
}

func (i *arrayDevices) Set(value string) error {
	*i = append(*i, value)
	return nil
}

var myDevices arrayDevices
