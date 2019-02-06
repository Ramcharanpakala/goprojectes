package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/tkanos/gonfig"
)

// Struct to hold DBConfig

type DBConfig struct {
	DBServer     string
	DBServerPort int
	DBUserName   string
	DBPassword   string
	DBName       string
}

// Struct to hold NGCSLogConfig

type NGCSLogConfig struct {
	LocalLogServer      string
	LocalLogServerPort  int
	RemoteLogServer     string
	RemoteLogServerPort int
	LogLocally          int
	LogRemotely         int
}
type Loop_Data struct {
	Dtsp      float64 `json:"temp_sp"`
	Dtpv      float64 `json:"temp_pv"`
	Dhsp      float64 `json:"hum_sp"`
	Dhpv      float64 `json:"hum_pv"`
	Dpsp      float64 `json:"press_sp"`
	Dppv      float64 `json:"press_pv"`
	Ddatatime string  `json:"date_time_date"`
}

var ngcsLogConfig NGCSLogConfig

var db *sql.DB

var err error

var router *gin.Engine

func main() {

	var dbConnectStr, ngcsLocalLogConnectStr string

	dbConnectStr = getDBConnectString()

	// Open the configured DB
	db, err = sql.Open("mysql", dbConnectStr)

	if err != nil {

		fmt.Println("Error: Unable to open DB Connection.")

		fmt.Println(err.Error())

		os.Exit(500)
	}

	defer db.Close()

	// Ensure that the connection is avaiable
	err = db.Ping()

	if err != nil {
		fmt.Println("Error: DB Connection is NOT available.")

		fmt.Println(err.Error())

		os.Exit(500)
	}

	// Initialise router, setup routes and wait for requests.
	router = gin.Default()

	initialiseRoutes()

	ngcsLocalLogConnectStr = getNGCSLocalLogServerConnectStr()

	router.Run(ngcsLocalLogConnectStr)

}

func initialiseRoutes() {

	router.GET("/Loop_Data", processLoop_data)

}

func processLoop_data(c *gin.Context) {

	stmt, err := db.Prepare("select temp_sp,temp_pv,hum_sp,hum_pv,press_sp,press_pv,date_time from  ZTK_Loop_Data")

	logs := []Loop_Data{}
	if err != nil {

		fmt.Print(err.Error())
	}

	rows, err := stmt.Query()
	if err != nil {
		fmt.Println(err)
	}
	for rows.Next() {
		var log Loop_Data
		err = rows.Scan(&log.Dtsp, &log.Dtpv, &log.Dhsp, &log.Dhpv, &log.Dpsp, &log.Dppv, &log.Ddatatime)
		if err != nil {
			fmt.Println(err)
		}
		logs = append(logs, log)
	}

	fmt.Println(logs)

	defer stmt.Close()
	c.JSON(http.StatusOK, logs)
}

// Read the contents of the DBConfig, form the dbConnectStr
// and return the same to the caller.
func getDBConnectString() string {

	dbConfiguration := DBConfig{}

	err := gonfig.GetConf("../ngcsCfg/dbconfig.json", &dbConfiguration)

	if err != nil {

		fmt.Println("Error: Unable to open the DBConfig file.")

		fmt.Println(err.Error())

		os.Exit(500)
	}

	// Now form the dbConnectStr and return the same.

	var dbConnectionStr string

	dbConnectionStr += dbConfiguration.DBUserName
	dbConnectionStr += ":"
	dbConnectionStr += dbConfiguration.DBPassword
	dbConnectionStr += "@tcp("
	dbConnectionStr += dbConfiguration.DBServer
	dbConnectionStr += ":"
	dbConnectionStr += strconv.Itoa(dbConfiguration.DBServerPort)
	dbConnectionStr += ")/"
	dbConnectionStr += dbConfiguration.DBName
	dbConnectionStr += "?charset=utf8"

	return dbConnectionStr
}

// Read the contents of the NGCSLogConfig
func readNGCSLogConfig() {

	err := gonfig.GetConf("../ngcsCfg/ngcsLogConfig.json", &ngcsLogConfig)

	if err != nil {

		fmt.Println("Error: Unable to open the NGCS Log Config file.")

		fmt.Println(err.Error())

		os.Exit(500)
	}
}

// Create and return the connect string for the NGCS Local Log Server
func getNGCSLocalLogServerConnectStr() string {

	var connectStr string

	readNGCSLogConfig()

	connectStr = ngcsLogConfig.LocalLogServer

	connectStr += ":"

	connectStr += strconv.Itoa(ngcsLogConfig.LocalLogServerPort)

	return connectStr
}
