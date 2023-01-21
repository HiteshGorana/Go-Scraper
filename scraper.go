package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"sync"
	"time"
)

var wg sync.WaitGroup

type R06 struct {
	R01 string  `json:"r01"`
	R02 float32 `json:"r02"`
	R03 int     `json:"r03"`
	R04 int     `json:"r04"`
	R05 int     `json:"r05"`
	R06 int     `json:"r06"`
	R07 int     `json:"r07"`
	R08 float32 `json:"r08"`
	R09 float32 `json:"r09"`
	R10 int     `json:"r10"`
	R11 int     `json:"r11"`
	R12 int     `json:"r12"`
	R13 int     `json:"r13"`
	R14 int     `json:"r14"`
	R15 string  `json:"r15"`
	R16 string  `json:"r16"`
	R17 string  `json:"r17"`
	R18 int     `json:"r18"`
	R20 string  `json:"r20"`
	R21 float32 `json:"r21"`
	R24 string  `json:"r24"`
	R25 string  `json:"r25"`
	R26 string  `json:"r26"`
	R27 string  `json:"r27"`
	R28 string  `json:"r28"`
	R30 string  `json:"r30"`
	R31 string  `json:"r31"`
}
type Course struct {
	R01 string   `json:"r01"`
	R02 string   `json:"r02"`
	R03 string   `json:"r03"`
	R04 string   `json:"r04"`
	R05 []string `json:"r05"`
	R06 []R06    `json:"r06"`
}

type PayloadRequest struct {
	P01 []int  `json:"p01"`
	P02 string `json:"p02"`
	P03 string `json:"p03"`
	P04 string `json:"p04"`
	P05 int    `json:"p05"`
	P06 int    `json:"p06"`
	P07 bool   `json:"p07"`
}

func request(date string, i int) {
	defer wg.Done()
	var data Course
	url := "https://camarillosprings.ezlinksgolf.com/api/search/search"
	method := "POST"
	var PayloadData = PayloadRequest{
		P01: []int{5885, 23372},
		P02: date,
		P03: "5:00 AM",
		P04: "7:00 PM",
		P05: 0,
		P06: 2,
		P07: false,
	}
	payloadBytes, err := json.Marshal(PayloadData)
	payload := bytes.NewReader(payloadBytes)

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Add("Accept", "application/json, text/plain, */*")
	req.Header.Add("Accept-Language", "en-US,en;q=0.9")
	req.Header.Add("Connection", "keep-alive")
	req.Header.Add("Content-Type", "application/json; charset=UTF-8")
	req.Header.Add("Origin", "https://camarillosprings.ezlinksgolf.com")
	req.Header.Add("Referer", "https://camarillosprings.ezlinksgolf.com/index.html")
	req.Header.Add("Sec-Fetch-Dest", "empty")
	req.Header.Add("Sec-Fetch-Mode", "cors")
	req.Header.Add("Sec-Fetch-Site", "same-origin")
	req.Header.Add("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/104.0.0.0 Safari/537.36")
	req.Header.Add("sec-ch-ua", "\"Chromium\";v=\"104\", \" Not A;Brand\";v=\"99\", \"Google Chrome\";v=\"104\"")
	req.Header.Add("sec-ch-ua-mobile", "?0")
	req.Header.Add("sec-ch-ua-platform", "\"Linux\"")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
	}
	err = json.Unmarshal(body, &data)
	if err == nil {
		fmt.Println("Struct is:", data.R06[0].R16, "Goroutine #", i)
		fmt.Println("Struct is:", data.R06[0].R24)
		JsonData, _ := json.Marshal(data)
		filename := fmt.Sprintf("data_%s.json", date)
		filename = strings.Join(strings.Split(filename, "/"), "-")
		err = ioutil.WriteFile(filename, JsonData, 0644)
	} else {
		fmt.Errorf("error")
	}
}
func main() {
	currentTime := time.Now()
	yyyy, mm, dd := currentTime.Date()
	tomorrow := time.Now()
	for i := 0; i < 10; i++ {
		wg.Add(1)
		tomorrow = time.Date(yyyy, mm, dd+i, 0, 0, 0, 0, currentTime.Location())
		go request(tomorrow.Format("01/02/2006"), i)
		fmt.Println("Goroutine #", i)
	}
	wg.Wait()
	fmt.Println("Completed")
}
