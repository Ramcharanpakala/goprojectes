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

	router.POST("/Logs_Event", processEvent_Log)
}
func processEvent_Log(c *gin.Context) {

	//fmt.Println("Hello")

	var finalResult int = 1

	var log Logs_Event
	c.BindJSON(&log)
	fmt.Println(log)

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
