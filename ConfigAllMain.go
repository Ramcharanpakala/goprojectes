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

/*
type Logs_Test struct {
	Tid         string `json:"log_id"`
	Tname       string `json:"log_name"`
	Tdatetime   string `json:"log_date_time_date"`
	Ttypeid     int    `json:"ZTK_Logs_Test_Type_id"`
	Tuserid     int    `json:"ZTK_Users_id"`
	Tcreatedby  int    `json:"created_by"`
	Tcreated    string `json:"created_date"`
	Tmodifiedby int    `json:"modified_by"`
	Tmodified   string `json:"modified_date"`
}

type Logs_Test_Type struct {
	Ltesttype    string `json:"test_type"`
	Tcreated1    string `json:"create_date"`
	Tmodified2   string `json:"modified_date"`
	Tcreatedby1  int    `json:"created_by"`
	Tmodifiedby2 int    `json:"modified_by"`
} */

type Logs_Event struct {
	Lid        string `json:"log_id"`
	Pname      string `json:"program_name"`
	Pdatetime  string `json:"program_date_time_date"`
	Etypeid    int    `json:"ZTK_Logs_Event_Type_id"`
	Eid        int    `json:"ZTK_Users_id"`
	Createdby  int    `json:"created_by"`
	Ecreated   string `json:"created_date"`
	Modifiedby int    `json:"modified_by"`
	Emodified  string `json:"modified_date"`
}

type Logs_Event_Type struct {
	Levents    string `json:"events_type"`
	Lcreated   int    `json:"created_by"`
	Lmodified  int    `json:"modified_by"`
	Lcreated1  string `json:"create_date"`
	Lmodified2 string `json:"modified_date"`
}

/*
type Logs_Maintenance struct {
	Mname       string `json:"component_name"`
	Mruntime    int    `json:"runtime_hr"`
	Mcounter    int    `json:"counter"`
	Mservice    int    `json:"days_till_service"`
	Mpending    int    `json:"maintenance_pending"`
	Mstatus     int    `json:"maintenance_status"`
	Mcreated    string `json:"created_date"`
	Mmodified   string `json:"modified_date"`
	Mcreatedby  int    `json:"created_by"`
	Mmodifiedby int    `json:"modified_by"`
}
*/
func main() {

	for i := 1; i <= 10; i++ {
		fmt.Println(i)

		// calculating loop time
		start := time.Now()
		time.Sleep(time.Second * 2)

		elapsed := time.Since(start)
		fmt.Printf("loop time %s", elapsed)

		logconfigallmain_loadjson()
	}

}

func logconfigallmain_loadjson() {

	log := Logs_Event{}

	err := gonfig.GetConf("ConfigAll.json", &log)
	if err != nil {
		panic(err)
	}
	log_out, error := json.Marshal(log)
	if error != nil {
		panic(error)
	}
	//	fmt.Println(string(log_out))

	url := "http://127.0.0.1:8181/Logs_All"
	fmt.Println("URL:>", url)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(log_out))
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
