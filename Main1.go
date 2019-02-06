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
// 1.1        18Jan2019    RAM        Type declaration of All logs(5)
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
	LogEventType Logs_Event_Type `json:"log_event_type`
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
*/
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

	router.POST("/Logs_All", processAll_Logs)
}

func processAll_Logs(c *gin.Context) {

	//fmt.Println("Hello")

	var finalResult int = 1

	var log Logs_Event
	c.BindJSON(&log)
	fmt.Println(log)


	var logEventType := log.LogEventType

	stmt, err := db.Prepare("insert into ZTK_Logs_Event_Type (events_type,created_by,modified_by,created,modified ) values(?,?,?,?,?);")

	if err != nil {

		fmt.Print("Error: Creating Prepared Statement")

		fmt.Print(err.Error())

		finalResult = 0
	}

	_, err = stmt.Exec(logEventType.Levents, logEventType.Lcreated, logEventType.Lmodified, logEventType.Lcreated1, logEventType.Lmodified2)

	if err != nil {

		fmt.Print("Error: Executing Prepared Statement")

		fmt.Print(err.Error())

		finalResult = 0
	}

	defer stmt.Close()

/*	if finalResult != 0 {

		c.JSON(http.StatusOK, gin.H{
			"Status = 1 ": fmt.Sprintf(" %s - Event_type  Log recorded.", logEventType.Levents),
			"Status = 2 ": fmt.Sprintf(" %s - Created_type  Log recorded.", logEventType.Lcreated),
			"Status = 3 ": fmt.Sprintf(" %s - Modified_type  Log recorded.", logEventType.Lmodified),
			"Status = 4 ": fmt.Sprintf(" %s - Created1 Log recorded.", logEventType.Lcreated1),
			"Status = 5 ": fmt.Sprintf(" %s - Modified2  Log recorded.", logEventType.Lmodified2),
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
*/

	stmt, err := db.Prepare("insert into ZTK_Logs_Event (log_id,program_name,program_date_time,ZTK_Logs_Event_Type_id,ZTK_Users_id,created_by,created,modified_by,modified ) values(?,?,?,?,?,?,?,?,?);")

	if err != nil {

		fmt.Print("Error: Creating Prepared Statement")

		fmt.Print(err.Error())

		finalResult = 0
	}

	_, err = stmt.Exec(log.Lid, log.Pname, log.Pdatetime, log.Etypeid, log.Eid, log.Createdby, log.Ecreated, log.Modifiedby, log.Emodified)

	if err != nil {

		fmt.Print("Error: Executing Prepared Statement")

		fmt.Print(err.Error())

		finalResult = 0
	}

	defer stmt.Close()

	if finalResult != 0 {

		c.JSON(http.StatusOK, gin.H{
			"Status = 1 ": fmt.Sprintf(" %s - Id  Log recorded.", log.Lid),
			"Status = 2 ": fmt.Sprintf(" %s - name  Log recorded.", log.Pname),
			"Status = 3 ": fmt.Sprintf(" %s - Datetime  Log recorded.", log.Pdatetime),
			"Status = 4 ": fmt.Sprintf(" %s - Etype Log recorded.", log.Etypeid),
			"Status = 5 ": fmt.Sprintf(" %s - EiD  Log recorded.", log.Eid),
			"Status = 6":  fmt.Sprintf(" %s - Createdby  Log recorded.", log.Createdby),
			"Status = 7":  fmt.Sprintf(" %s - created  Log recorded.", log.Ecreated),
			"Status = 8":  fmt.Sprintf(" %s - Modifiedby  Log recorded.", log.Modifiedby),
			"Status = 9 ": fmt.Sprintf(" %s - Modified  Log recorded.", log.Emodified),
		})

	} else {

		c.JSON(http.StatusOK, gin.H{
			"Status = -1 ": fmt.Sprintf(" %s - Error of Id Log.", log.Lid),
			"Status = -2 ": fmt.Sprintf(" %s - Error of name Log.", log.Pname),
			"Status = -3 ": fmt.Sprintf(" %s - Error of Datetime Log.", log.Pdatetime),
			"Status = -4 ": fmt.Sprintf(" %s - Error of Etype Log.", log.Etypeid),
			"Status = -5 ": fmt.Sprintf(" %s - Error of Eid Log.", log.Eid),
			"Status = -6 ": fmt.Sprintf(" %s - Error of Createdby Log.", log.Createdby),
			"Status = -7 ": fmt.Sprintf(" %s - Error of Created Log.", log.Ecreated),
			"Status = -8 ": fmt.Sprintf(" %s - Error of Modifiedby Log.", log.Modifiedby),
			"Status = -9 ": fmt.Sprintf(" %s - Error of Modified Log.", log.Emodified),
		})
	

}
/*
func processTest_Log(c *gin.Context) {

	//fmt.Println("Hello")

	var finalResult int = 1

	var log Logs_Test
	c.BindJSON(&log)
	fmt.Println(log)

	stmt, err := db.Prepare("insert into ZTK_Logs_Test (log_id,log_name,log_date_time,ZTK_Logs_Test_Type_id,ZTK_Users_id,created_by,created,modified_by,modified ) values(?,?,?,?,?,?,?,?,?);")

	if err != nil {

		fmt.Print("Error: Creating Prepared Statement")

		fmt.Print(err.Error())

		finalResult = 0
	}

	_, err = stmt.Exec(log.Tid, log.Tname, log.Tdatetime, log.Ttypeid, log.Tuserid, log.Tcreatedby, log.Tcreated, log.Tmodifiedby, log.Tmodified)

	if err != nil {

		fmt.Print("Error: Executing Prepared Statement")

		fmt.Print(err.Error())

		finalResult = 0
	}

	defer stmt.Close()

	if finalResult != 0 {

		c.JSON(http.StatusOK, gin.H{
			"Status = 1 ": fmt.Sprintf(" %s - id  Log recorded.", log.Tid),
			"Status = 2 ": fmt.Sprintf(" %s - name  Log recorded.", log.Tname),
			"Status = 3 ": fmt.Sprintf(" %s - datatime  Log recorded.", log.Tdatetime),
			"Status = 4 ": fmt.Sprintf(" %s - typeid Log recorded.", log.Ttypeid),
			"Status = 5 ": fmt.Sprintf(" %s - userid  Log recorded.", log.Tuserid),
			"Status = 6 ": fmt.Sprintf(" %s - createdby  Log recorded.", log.Tcreatedby),
			"Status = 7 ": fmt.Sprintf(" %s - created  Log recorded.", log.Tcreated),
			"Status = 8 ": fmt.Sprintf(" %s - Modifiedby  Log recorded.", log.Tmodifiedby),
			"Status = 9 ": fmt.Sprintf(" %s - Modified  Log recorded.", log.Tmodified),
		})

	} else {

		c.JSON(http.StatusOK, gin.H{
			"Status = -1 ": fmt.Sprintf(" %s - Error of id Log.", log.Tid),
			"Status = -2 ": fmt.Sprintf(" %s - Error of name Log.", log.Tname),
			"Status = -3 ": fmt.Sprintf(" %s - Error of datetime Log.", log.Tdatetime),
			"Status = -4 ": fmt.Sprintf(" %s - Error of typeid Log.", log.Ttypeid),
			"Status = -5 ": fmt.Sprintf(" %s - Error of userid Log.", log.Tuserid),
			"Status = -6 ": fmt.Sprintf(" %s - Error of createdby Log.", log.Tcreatedby),
			"Status = -7 ": fmt.Sprintf(" %s - Error of Created Log.", log.Tcreated),
			"Status = -8 ": fmt.Sprintf(" %s - Error of Modifiedby Log.", log.Tmodifiedby),
			"Status = -9 ": fmt.Sprintf(" %s - Error of Modified Log.", log.Tmodified),
		})
	}
}

func processTest_typeLog(c *gin.Context) {

	//fmt.Println("Hello")

	var finalResult int = 1

	var log Logs_Test_Type
	c.BindJSON(&log)
	fmt.Println(log)

	stmt, err := db.Prepare("insert into ZTK_Logs_Test_Type (test_type,created,modified,created_by,modified_by ) values(?,?,?,?,?);")

	if err != nil {

		fmt.Print("Error: Creating Prepared Statement")

		fmt.Print(err.Error())

		finalResult = 0
	}

	_, err = stmt.Exec(log.Ltesttype, log.Tcreated1, log.Tmodified2, log.Tcreatedby1, log.Tmodifiedby2)

	if err != nil {

		fmt.Print("Error: Executing Prepared Statement")

		fmt.Print(err.Error())

		finalResult = 0
	}

	defer stmt.Close()

	if finalResult != 0 {

		c.JSON(http.StatusOK, gin.H{
			"Status = 1 ": fmt.Sprintf(" %s - Test_type  Log recorded.", log.Ltesttype),
			"Status = 2 ": fmt.Sprintf(" %s - Created1_type  Log recorded.", log.Tcreated1),
			"Status = 3 ": fmt.Sprintf(" %s - Modified2_type  Log recorded.", log.Tmodified2),
			"Status = 4 ": fmt.Sprintf(" %s - Createdby1 Log recorded.", log.Tcreatedby1),
			"Status = 5 ": fmt.Sprintf(" %s - Modifiedby2  Log recorded.", log.Tmodifiedby2),
		})

	} else {

		c.JSON(http.StatusOK, gin.H{
			"Status = -1 ": fmt.Sprintf(" %s - Error of Test_type Log.", log.Ltesttype),
			"Status = -2 ": fmt.Sprintf(" %s - Error of Created1_type Log.", log.Tcreated1),
			"Status = -3 ": fmt.Sprintf(" %s - Error of Modified2_type Log.", log.Tmodified2),
			"Status = -4 ": fmt.Sprintf(" %s - Error of Createdby1 Log.", log.Tcreatedby1),
			"Status = -5 ": fmt.Sprintf(" %s - Error of Modifiedby2 Log.", log.Tmodifiedby2),
		})
	}
}

func processMaintenance_Log(c *gin.Context) {

	//fmt.Println("Hello")

	var finalResult int = 1

	var log Logs_Maintenance
	c.BindJSON(&log)
	fmt.Println(log)

	stmt, err := db.Prepare("insert into ZTK_Logs_Maintenance (component_name,runtime_hr,counter,days_till_service,maintenance_pending,maintenance_status,created,modified,created_by,modified_by ) values(?,?,?,?,?,?,?,?,?,?);")

	if err != nil {

		fmt.Print("Error: Creating Prepared Statement")

		fmt.Print(err.Error())

		finalResult = 0
	}

	_, err = stmt.Exec(log.Mname, log.Mruntime, log.Mcounter, log.Mservice, log.Mpending, log.Mstatus, log.Mcreated, log.Mmodified, log.Mcreatedby, log.Mmodifiedby)

	if err != nil {

		fmt.Print("Error: Executing Prepared Statement")

		fmt.Print(err.Error())

		finalResult = 0
	}

	defer stmt.Close()

	if finalResult != 0 {

		c.JSON(http.StatusOK, gin.H{
			"Status = 1 ":  fmt.Sprintf(" %s - name  Log recorded.", log.Mname),
			"Status = 2 ":  fmt.Sprintf(" %s - runtime  Log recorded.", log.Mruntime),
			"Status = 3 ":  fmt.Sprintf(" %s - counter  Log recorded.", log.Mcounter),
			"Status = 4 ":  fmt.Sprintf(" %s - service Log recorded.", log.Mservice),
			"Status = 5 ":  fmt.Sprintf(" %s - pending  Log recorded.", log.Mpending),
			"Status = 6":   fmt.Sprintf(" %s - status  Log recorded.", log.Mstatus),
			"Status = 7":   fmt.Sprintf(" %s - created  Log recorded.", log.Mcreated),
			"Status = 8":   fmt.Sprintf(" %s - Modified  Log recorded.", log.Mmodified),
			"Status = 9 ":  fmt.Sprintf(" %s - createdby  Log recorded.", log.Mcreatedby),
			"Status = 10 ": fmt.Sprintf(" %s - Modifiedby  Log recorded.", log.Mmodifiedby),
		})

	} else {

		c.JSON(http.StatusOK, gin.H{
			"Status = -1 ":  fmt.Sprintf(" %s - Error of name Log.", log.Mname),
			"Status = -2 ":  fmt.Sprintf(" %s - Error of runtime Log.", log.Mruntime),
			"Status = -3 ":  fmt.Sprintf(" %s - Error of counter Log.", log.Mcounter),
			"Status = -4 ":  fmt.Sprintf(" %s - Error of service Log.", log.Mservice),
			"Status = -5 ":  fmt.Sprintf(" %s - Error of pending Log.", log.Mpending),
			"Status = -6 ":  fmt.Sprintf(" %s - Error of status Log.", log.Mstatus),
			"Status = -7 ":  fmt.Sprintf(" %s - Error of Created Log.", log.Mcreated),
			"Status = -8 ":  fmt.Sprintf(" %s - Error of Modified Log.", log.Mmodified),
			"Status = -9 ":  fmt.Sprintf(" %s - Error of createdby Log.", log.Mcreatedby),
			"Status = -10 ": fmt.Sprintf(" %s - Error of Modifiedby Log.", log.Mmodifiedby),
		}) 
	}
} */

// Read the contents of the DBConfig, form the dbConnectStr
// and return the same to the caller.
func getDBConnectString() string {
}
	dbConfiguration := DBConfig{}

	err := gonfig.GetConf("../config/dbconfig.json", &dbConfiguration)

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


// Read the contents of the NGCSLogConfig
func readNGCSLogConfig() {

	err := gonfig.GetConf("../config/ngcsLogConfig.json", &ngcsLogConfig)

	if err != nil {

		fmt.Println("Error: Unable to open the NGCS Log Config file.")

		fmt.Println(err.Error())

		os.Exit(500)
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
