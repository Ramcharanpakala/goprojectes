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

	router.POST("/Logs_Test", processTest_Log)
}
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
