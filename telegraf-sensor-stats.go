package main

import (
	"flag"
	"fmt"
	"log"
	"strings"

	"github.com/dustin/go-coap"
)

func doCoapRequest(host string, url string) *coap.Message {
	req := coap.Message{
		Type:      coap.Confirmable,
		Code:      coap.GET,
		MessageID: 12345,
	}

	req.AddOption(coap.URIQuery, "s=1")
	req.SetPathString(url)

	c, err := coap.Dial("udp", host)
	if err != nil {
		log.Fatalf("Error dialing: %v", err)
	}

	rv, err := c.Send(req)
	if err != nil {
		log.Fatalf("Error sending request: %v", err)
	}

	return rv
}

func main() {
	var host string
	flag.StringVar(&host, "h", "", "COAP Host url")
	flag.Parse()

	if host == "" {
		log.Fatal("Please specify host using -h=")
	}

	rv := doCoapRequest(host, "/co2")
	if rv != nil {
		fmt.Printf("sensor_box,box=desk,channel=co2 val=%s\n", rv.Payload)
	}

	rv = doCoapRequest(host, "/ldr")
	if rv != nil {
		fmt.Printf("sensor_box,box=desk,channel=ldr val=%s\n", rv.Payload)
	}

	rv = doCoapRequest(host, "/temp")
	if rv != nil {
		vals := strings.Split(string(rv.Payload), ",")
		if len(vals) > 1 {
			fmt.Printf("sensor_box,box=desk,channel=temp temp=%s,hum=%s\n", vals[0], vals[1])
		}
	}
}
