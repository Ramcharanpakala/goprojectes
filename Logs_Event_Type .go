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
//
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
type Logs_Event_Type struct {
	Levents    string `json:"events_type"`
	Lcreated   int    `json:"created_by"`
	Lmodified  int    `json:"modified_by"`
	Lcreated1  string `json:"create_date"`
	Lmodified2 string `json:"modified_date"`
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

	router.POST("/Logs_Event_Type_1", processEvent_typeLog)
}

func processEvent_typeLog(c *gin.Context) {

	//fmt.Println("Hello")

	var finalResult int = 1

	var log Logs_Event_Type
	c.BindJSON(&log)
	fmt.Println(log)

	stmt, err := db.Prepare("insert into ZTK_Logs_Event_Type (events_type,created_by,modified_by,created,modified ) values(?,?,?,?,?);")

	if err != nil {

		fmt.Print("Error: Creating Prepared Statement")

		fmt.Print(err.Error())

		finalResult = 0
	}

	_, err = stmt.Exec(log.Levents, log.Lcreated, log.Lmodified, log.Lcreated1, log.Lmodified2)

	if err != nil {

		fmt.Print("Error: Executing Prepared Statement")

		fmt.Print(err.Error())

		finalResult = 0
	}

	defer stmt.Close()

	if finalResult != 0 {

		c.JSON(http.StatusOK, gin.H{
			"Status = 1 ": fmt.Sprintf(" %s - Event_type  Log recorded.", log.Levents),
			"Status = 2 ": fmt.Sprintf(" %s - Created_type  Log recorded.", log.Lcreated),
			"Status = 3 ": fmt.Sprintf(" %s - Modified_type  Log recorded.", log.Lmodified),
			"Status = 4 ": fmt.Sprintf(" %s - Created1 Log recorded.", log.Lcreated1),
			"Status = 5 ": fmt.Sprintf(" %s - Modified2  Log recorded.", log.Lmodified2),
		})

	} else {

		c.JSON(http.StatusOK, gin.H{
			"Status = -1 ": fmt.Sprintf(" %s - Error of Event_type Log.", log.Levents),
			"Status = -2 ": fmt.Sprintf(" %s - Error of Created_type Log.", log.Lcreated),
			"Status = -3 ": fmt.Sprintf(" %s - Error of Modified_type Log.", log.Lmodified),
			"Status = -4 ": fmt.Sprintf(" %s - Error of Created1 Log.", log.Lcreated1),
			"Status = -5 ": fmt.Sprintf(" %s - Error of Modified2 Log.", log.Lmodified2),
		})
	}
}

// Read the contents of the DBConfig, form the dbConnectStr
// and return the same to the caller.
func getDBConnectString() string {

	dbConfiguration := DBConfig{}

	err := gonfig.GetConf("../ngcsLogger/dbconfig.json", &dbConfiguration)

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

	err := gonfig.GetConf("../ngcsLogger/ngcsLogConfig.json", &ngcsLogConfig)

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
