package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"foe-client/cmd/properties"
	"github.com/robfig/cron/v3"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

func main() {
	fmt.Println("Starting Foe")
	propertiesPath := parseArgs()
	if propertiesPath == "" {
		propertiesPath = "~/.foe/app.properties"
	}
	properties.LoadProperties(propertiesPath)

	c := cron.New()
	c.AddFunc(properties.GetProperties().RunInterval, retrieveDataAndPhoneHome)
	c.Start()

	go forever()
	select {}
}

func forever() {
	for {
		fmt.Printf("%v+\n", time.Now())
		time.Sleep(time.Minute)
	}
}

func parseArgs() string {
	if len(os.Args[0:]) > 1 {
		arg := os.Args[1]
		args := strings.Split(arg, "=")
		if args[0] == "--config" {
			return args[1]
		} else {
			printUsage()
			os.Exit(1)
		}
	}
	return ""
}

func printUsage() {
	fmt.Println("Incorrect usage, supported arguments are:")
	fmt.Println("--config=/path/to/config/file (defaults to ~/.foe/app.properties)")
}

type IpInfo struct {
	Ip       string `json:"ip"`
	Hostname string `json:"hostname"`
	City     string `json:"city"`
	Region   string `json:"region"`
	Country  string `json:"country"`
	Loc      string `json:"loc"`
	Org      string `json:"org"`
	Postal   string `json:"postal"`
	Timezone string `json:"timezone"`
	Readme   string `json:"readme"`
}

func retrieveDataAndPhoneHome() {
	log.Println("Retrieve and phone home")
	ipInfo := retrieveData()
	phoneHome(ipInfo)
}

func retrieveData() IpInfo {
	res, err := http.Get("https://ipinfo.io/json")
	if err != nil {
		log.Print(err)
	}

	defer res.Body.Close()
	bodyBytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Print(err)
	}

	var ipInfo IpInfo
	json.Unmarshal(bodyBytes, &ipInfo)
	return ipInfo
}

func phoneHome(ipInfo IpInfo) {
	jsonData, err := json.Marshal(ipInfo)
	if err != nil {
		log.Println(err)
	}

	res, err := http.Post(properties.GetProperties().RemoteUrl, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Print(err)
	}
	log.Println(res.StatusCode)
}
