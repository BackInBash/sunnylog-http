package main

// Fetch todays production from an SMA Sunny Boy via the HTTP interface
// Print it on stdout for now

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
)

var InfluxURL = ""
var InfluxToken = ""

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func loginToken(baseURL string, user string, pwd string) string {
	loginUrl := baseURL + "/dyn/login.json"
	values := map[string]string{"right": user, "pass": pwd}
	jsonValue, _ := json.Marshal(values)
	resp, err := http.Post(loginUrl, "application/json;charset=UTF-8", bytes.NewBuffer(jsonValue))
	check(err)

	fmt.Println("login Status:", resp.Status)
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("response Body:", string(body))

	var loginResult map[string]interface{}
	err = json.Unmarshal(body, &loginResult)
	check(err)

	sidMap := loginResult["result"].(map[string]interface{})
	return sidMap["sid"].(string)
}

func logout(baseURL string, token string) {
	logoutUrl := baseURL + "/dyn/logout.json?sid=" + token
	resp, err := http.Post(logoutUrl, "application/json", bytes.NewBuffer([]byte("")))
	check(err)

	fmt.Println("logout Status:", resp.Status)
}

func gettimestamp() int64 {
	now := time.Now()
	midnightUtc := fmt.Sprintf("%d-%02d-%02dT00:00:00-00:00",
		now.Year(), now.Month(), now.Day())
	form := "2006-01-02T15:04:05-07:00"
	t2, e := time.Parse(form, midnightUtc)
	check(e)
	return t2.Unix()
}

func getlog(baseURL string, timeFrom int64, timeTo int64, token string) {
	logUrl := baseURL + "/dyn/getLogger.json?sid=" + token
	type LogRequest struct {
		DestDev []int `json:"destDev"`
		Key     int   `json:"key"`
		TStart  int64 `json:"tStart"`
		TEnd    int64 `json:"tEnd"`
	}

	//Log keys are hardcoded in SMA scripts.js
	totWhOut5min := 28672
	values := &LogRequest{
		DestDev: make([]int, 0),
		Key:     totWhOut5min,
		TStart:  timeFrom,
		TEnd:    timeTo}
	jsonValue, e := json.Marshal(values)
	check(e)
	resp, err := http.Post(logUrl, "application/json", bytes.NewBuffer(jsonValue))
	check(err)
	fmt.Println("logger Status:", resp.Status)
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("response Body:", string(body))
	var logResult map[string]interface{}
	err = json.Unmarshal(body, &logResult)
	check(err)
	var inverterLog = logResult["result"].(map[string]interface{})
	for inverter, logValues := range inverterLog {
		savelog(inverter, logValues.([]interface{}))
	}
}

func savelog(inverter string, logValues []interface{}) {
	fmt.Println(inverter)

	client := influxdb2.NewClientWithOptions(InfluxURL, InfluxToken, influxdb2.DefaultOptions().SetBatchSize(20))

	writeAPI := client.WriteAPIBlocking("sunnylog", "sunnylog")

	tags := map[string]string{"inverter": inverter}
	for _, value := range logValues {
		dataPoint := value.(map[string]interface{})
		fields := map[string]interface{}{
			"watt_hours": dataPoint["v"],
		}
		pt := influxdb2.NewPoint("production",
			tags,
			fields,
			time.Unix(int64(dataPoint["t"].(float64)), 0))

		err := writeAPI.WritePoint(context.Background(), pt)
		if err != nil {
			panic(err)
		}
	}
	// Ensures background processes finishes
	client.Close()
}

func main() {

	var baseURL, password = "", ""

	if os.Args == nil {
		panic("No CLI Args specified!")
	}

	for index, arg := range os.Args {
		if arg == "--solarUrl" {
			baseURL = os.Args[index+1]
		}
		if arg == "--solarPassword" {
			password = os.Args[index+1]
		}
		if arg == "--influxUrl" {
			InfluxURL = os.Args[index+1]
		}
		if arg == "--influxToken" {
			InfluxToken = os.Args[index+1]
		}
	}

	// Check for Base Url
	if baseURL == "" {
		log.Fatalf("No SunnyBoy Base URL specified!")
	}
	// Check for API Key
	if InfluxToken == "" {
		log.Fatalf("No API Key specified!")
	}
	// Check for FloatingIP
	if password == "" {
		log.Fatalf("No SunnyBoy Password specified!")
	}
	if baseURL == "" || password == "" || InfluxURL == "" {
		print("No CLI Parameter...\n --solarUrl      |  SunnyBoy Web URL\n --solarPassword |  SunnyBoy Password\n --influxUrl     |  InfluxDB API\n")
		return
	}

	token := loginToken(baseURL, "usr", password)
	timeFrom := gettimestamp()
	timeTo := time.Now().Unix()
	getlog(baseURL, timeFrom, timeTo, token)
	logout(baseURL, token)
}
