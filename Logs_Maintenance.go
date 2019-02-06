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

	router.POST("/Logs_Maintenance", processMaintenance_Log)
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
