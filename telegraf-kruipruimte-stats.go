package main

import (
	"encoding/json"
	"fmt"
	"log"
	"flag"
)

var (
	baseURL       = ""
	serialIndoor  = "28ae6f2c0600009c"
	serialOutdoor = "28831b2b06000094"
)

type SensorContainer struct {
	Sensors []JSONSensor `json:"sensors"`
}

type JSONSensor struct {
	Serial      string  `json:"serial"`
	Status      string  `json:"status"`
	Temperature float32 `json:"celsius"`
}

func main() {
	flag.StringVar(&baseURL, "h", "", "Sensor host url.")
	flag.Parse()

	if baseURL == "" {
		log.Fatal("Please pass sensor url using -h")
	}

	data, err := FetchURL(baseURL)
	if err != nil {
		log.Fatalf("Error dialing: %v", err)
	}

	var sensors SensorContainer
	if xmlerr := json.Unmarshal([]byte(data), &sensors); xmlerr != nil {
		panic(err)
	}

	for _, sensor := range sensors.Sensors {
		if sensor.Serial == serialIndoor {
			sensor.Serial = "indoor"
		} else if sensor.Serial == serialOutdoor {
			sensor.Serial = "outdoor"
		}
		fmt.Printf("kruipruimte,server=bs0,serial=%s temperature=%f\n", sensor.Serial, sensor.Temperature)
	}
}
