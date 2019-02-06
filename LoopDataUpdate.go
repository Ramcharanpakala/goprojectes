package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/tkanos/gonfig"
)

type Loop_DataUpdate struct {
	Dtsp      float64 `json:"temp_sp"`
	Dtpv      float64 `json:"temp_pv"`
	Dhsp      float64 `json:"hum_sp"`
	Dhpv      float64 `json:"hum_pv"`
	Dpsp      float64 `json:"press_sp"`
	Dppv      float64 `json:"press_pv"`
	Ddatatime string  `json:"date_time_date"`
}

func main() {

	for i := 1; i <= 2; i++ {
		fmt.Println(i)

		// calculating loop time
		start := time.Now()
		time.Sleep(time.Second * 2)

		elapsed := time.Since(start)
		fmt.Printf("loop time %s", elapsed)

		logloopdataupdate_loadjson()
	}

}

func logloopdataupdate_loadjson() {

	log := Loop_DataUpdate{}

	err := gonfig.GetConf("LoopDataUpdate.json", &log)
	if err != nil {
		panic(err)
	}
	log_out, error := json.Marshal(log)
	if error != nil {
		panic(error)
	}
	//	fmt.Println(string(log_out))

	url := "http://127.0.0.1:8181/Loop_Data/:date_time_date"
	fmt.Println("URL:>", url)

	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(log_out))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	fmt.Println("response Status:", resp.Status)
	fmt.Println("response Headers:", resp.Header)
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("response Body:", string(body))
}
