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
type Logs_Test_Type struct {
	Ltesttype    string `json:"test_type"`
	Tcreated1    string `json:"create_date"`
	Tmodified2   string `json:"modified_date"`
	Tcreatedby1  int    `json:"created_by"`
	Tmodifiedby2 int    `json:"modified_by"`
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

	router.POST("/Logs_Test_Type", processTest_typeLog)
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
