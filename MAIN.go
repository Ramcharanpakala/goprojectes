//------------------------------------------------------------------------------
// Filename:    main.go
// Desc:        Contains the logic for the NGCSLocalLogServer. Processes the
//              following requests listed below and stores them into the
//              appropriate DB tables.
//                  * diagnostic log
//                  * program execution log
// Rev History:
//
// Ver#       Date         Author     Desc
//------------------------------------------------------------------------------
// 1.0        27Sep2018    GCB        Initial Creation
// 1.1        18Jan2019    RAM        Type declaration of All logs(5) Creation
// 1.2        21Jan2019    RAM        changes  of initial setup routes
// Copyright (c) 2018, Zetatek Technologies Pvt Ltd.
// Developed by CheckSum InfoSoft Pvt Ltd.
//------------------------------------------------------------------------------

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

// Type declaration of All logs
// Struct to hold Logs_Event

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

// Struct to hold Logs_Event_Type

type Logs_Event_Type struct {
	Levents    string `json:"events_type"`
	Lcreated   int    `json:"created_by"`
	Lmodified  int    `json:"modified_by"`
	Lcreated1  string `json:"create_date"`
	Lmodified2 string `json:"modified_date"`
}

// Struct to hold Logs_Test

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

// Struct to hold Logs_Test_Type

type Logs_Test_Type struct {
	Ltesttype    string `json:"test_type"`
	Tcreated1    string `json:"create_date"`
	Tmodified2   string `json:"modified_date"`
	Tcreatedby1  int    `json:"created_by"`
	Tmodifiedby2 int    `json:"modified_by"`
}

// Struct to hold Logs_Maintenance

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

	router.GET("/Logs_Event", processEvent_Log)
	router.GET("/Logs_Event_Type", processEvent_typeLog)
	router.GET("/Logs_Test", processTest_Log)
	router.GET("/Logs_Test_Type", processTest_typeLog)
	router.GET("/Logs_Maintenance", processMaintenance_Log)
}

func processEvent_Log(c *gin.Context) {

	stmt, err := db.Prepare("select log_id,program_name,program_date_time,ZTK_Logs_Event_Type_id,ZTK_Users_id,created_by,created,modified_by,modified from  ZTK_Logs_Event")

	logs := []Logs_Event{}
	if err != nil {

		fmt.Print(err.Error())
	}

	rows, err := stmt.Query()
	if err != nil {
		fmt.Println(err)
	}
	for rows.Next() {
		var log Logs_Event
		err = rows.Scan(&log.Lid, &log.Pname, &log.Pdatetime, &log.Etypeid, &log.Eid, &log.Createdby, &log.Ecreated, &log.Modifiedby, &log.Emodified)
		if err != nil {
			fmt.Println(err)
		}
		logs = append(logs, log)
	}

	fmt.Println(logs)

	defer stmt.Close()
	c.JSON(http.StatusOK, logs)
}

func processEvent_typeLog(c *gin.Context) {

	stmt, err := db.Prepare("select events_type,created,modified,created_by,modified_by from  ZTK_Logs_Event_Type")

	logs := []Logs_Event_Type{}
	if err != nil {

		fmt.Print(err.Error())
	}

	rows, err := stmt.Query()
	if err != nil {
		fmt.Println(err)
	}
	for rows.Next() {
		var log Logs_Event_Type
		err = rows.Scan(&log.Levents, &log.Lcreated, &log.Lmodified, &log.Lcreated1, &log.Lmodified2)
		if err != nil {
			fmt.Println(err)
		}
		logs = append(logs, log)
	}

	fmt.Println(logs)

	defer stmt.Close()
	c.JSON(http.StatusOK, logs)
}

func processTest_Log(c *gin.Context) {

	stmt, err := db.Prepare("select log_id,log_name,log_date_time,ZTK_Logs_Test_Type_id,ZTK_Users_id,created_by,created,modified_by,modified from  ZTK_Logs_Test")

	logs := []Logs_Test{}
	if err != nil {

		fmt.Print(err.Error())
	}

	rows, err := stmt.Query()
	if err != nil {
		fmt.Println(err)
	}
	for rows.Next() {
		var log Logs_Test
		err = rows.Scan(&log.Tid, &log.Tname, &log.Tdatetime, &log.Ttypeid, &log.Tuserid, &log.Tcreatedby, &log.Tcreated, &log.Tmodifiedby, &log.Tmodified)
		if err != nil {
			fmt.Println(err)
		}
		logs = append(logs, log)
	}

	fmt.Println(logs)

	defer stmt.Close()
	c.JSON(http.StatusOK, logs)
}

func processTest_typeLog(c *gin.Context) {

	stmt, err := db.Prepare("select test_type,created,modified,created_by,modified_by from  ZTK_Logs_Test_Type")

	logs := []Logs_Test_Type{}
	if err != nil {

		fmt.Print(err.Error())
	}

	rows, err := stmt.Query()
	if err != nil {
		fmt.Println(err)
	}
	for rows.Next() {
		var log Logs_Test_Type
		err = rows.Scan(&log.Ltesttype, &log.Tcreated1, &log.Tmodified2, &log.Tcreatedby1, &log.Tmodifiedby2)
		if err != nil {
			fmt.Println(err)
		}
		logs = append(logs, log)
	}

	fmt.Println(logs)

	defer stmt.Close()
	c.JSON(http.StatusOK, logs)
}

func processMaintenance_Log(c *gin.Context) {

	stmt, err := db.Prepare("select component_name,runtime_hr,counter,days_till_service,maintenance_pending,maintenance_status,created,modified,created_by, modified_by from  ZTK_Logs_Maintenance")

	logs := []Logs_Maintenance{}
	if err != nil {

		fmt.Print(err.Error())
	}

	rows, err := stmt.Query()
	if err != nil {
		fmt.Println(err)
	}
	for rows.Next() {
		var log Logs_Maintenance
		err = rows.Scan(&log.Mname, &log.Mruntime, &log.Mcounter, &log.Mservice, &log.Mpending, &log.Mstatus, &log.Mcreated, &log.Mmodified, &log.Mcreatedby, &log.Mmodifiedby)
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
