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
	LocalLogServer        string
	LocalLogServerGetPort int
	RemoteLogServer       string
	RemoteLogServerPort   int
	LogLocally            int
	LogRemotely           int
}
type Io_card_Info struct {
	Iaddress    string `json:"card_address"`
	Itype       string `json:"card_type"`
	Iversion    string `json:"card_version"`
	Inumber     string `json:"card_serial_number"`
	Ikey        string `json:"secret_key"`
	Iid         int    `json:"customer_id"`
	Idate       string `json:"mfg_date_date"`
	Icreated    string `json:"created_date"`
	Imodified   string `json:"modified_date"`
	Icreatedby  int    `json:"created_by"`
	Imodifiedby int    `json:"modified_by"`
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

	router.GET("/get_io_card_info", processIocardinfo)
}
func processIocardinfo(c *gin.Context) {

	stmt, err := db.Prepare("select card_address,card_type,card_version,card_serial_number,secret_key,customer_id,mfg_date,created,modified,created_by, modified_by from  ZTK_IO_Card_Info")

	logs := []Io_card_Info{}
	if err != nil {

		fmt.Print(err.Error())
	}

	rows, err := stmt.Query()
	if err != nil {
		fmt.Println(err)
	}
	for rows.Next() {
		var log Io_card_Info
		err = rows.Scan(&log.Iaddress, &log.Itype, &log.Iversion, &log.Inumber, &log.Ikey, &log.Iid, &log.Idate, &log.Icreated, &log.Imodified, &log.Icreatedby, &log.Imodifiedby)
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

	err := gonfig.GetConf("../../config/dbconfig.json", &dbConfiguration)

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

	err := gonfig.GetConf("../../config/ngcsLogConfig.json", &ngcsLogConfig)

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

	connectStr += strconv.Itoa(ngcsLogConfig.LocalLogServerGetPort)

	return connectStr
}
