package main

// Fetch todays production from an SMA Sunny Boy via the HTTP interface
// Print it on stdout for now

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"
    "github.com/influxdata/influxdb/client/v2"
)

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func loginToken(baseUrl string, user string, pwd string) string {
	loginUrl := baseUrl + "/dyn/login.json"
	values := map[string]string{"right": user, "pass": pwd}
	jsonValue, _ := json.Marshal(values)
	resp, err := http.Post(loginUrl, "application/json", bytes.NewBuffer(jsonValue))
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

func logout(baseUrl string, token string) {
	logoutUrl := baseUrl + "/dyn/logout.json?sid=" + token
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

func getlog(baseUrl string, timeFrom int64, timeTo int64, token string) {
	logUrl := baseUrl + "/dyn/getLogger.json?sid=" + token
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
    c, err := client.NewHTTPClient(client.HTTPConfig{
        Addr: "http://localhost:8086",
    })
    if err != nil {
        fmt.Println("Error creating InfluxDB Client: ", err.Error())
    }
    defer c.Close()

    bp, _ := client.NewBatchPoints(client.BatchPointsConfig{
        Database:  "sunnylog",
        Precision: "s",
    })

    tags := map[string]string{"inverter": inverter}
    for _, value := range logValues {
        dataPoint := value.(map[string]interface{})
        fields := map[string]interface{}{
            "watt_hours": dataPoint["v"],
        }
        pt, err := client.NewPoint("production",
            tags,
            fields,
            time.Unix(int64(dataPoint["t"].(float64)), 0))
        if err != nil {
            fmt.Println("Error: ", err.Error())
        }
        bp.AddPoint(pt)
    }
    // Write the batch
    c.Write(bp)
}

func main() {
	baseUrl := os.Getenv("SMA_BASEURL")
	password := os.Getenv("SMA_PASSWORD")
	token := loginToken(baseUrl, "usr", password)
	timeFrom := gettimestamp()
	timeTo := time.Now().Unix()
	getlog(baseUrl, timeFrom, timeTo, token)
	logout(baseUrl, token)
}
