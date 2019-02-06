package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

type Logs_Test_Type struct {
	Ltesttype    string `json:"test_type"`
	Tcreated1    string `json:"create_date"`
	Tmodified2   string `json:"modified_date"`
	Tcreatedby1  int    `json:"created_by"`
	Tmodifiedby2 int    `json:"modified_by"`
}

func main() {

	response, err := http.Get("http://127.0.0.1:8080/Logs_Test_Type")
	if err != nil {
		fmt.Printf("%s", err)
		os.Exit(1)
	} else {
		defer response.Body.Close()
		contents, err := ioutil.ReadAll(response.Body)
		if err != nil {
			fmt.Printf("%s", err)
			os.Exit(1)
		}
		fmt.Printf("%s\n", string(contents))
	}
}
