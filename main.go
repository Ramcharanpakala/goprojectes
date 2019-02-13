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
// 1.1        18Jan2019    RAM        Type declaration of All logs(7) Creation
// 1.2        21Jan2019    RAM        changes  of initial setup routes
// 1.3        11 feb 2019  RAM        Activity log add
// Copyright (c) 2018, Zetatek Technologies Pvt Ltd.
// Developed by CheckSum InfoSoft Pvt Ltd.
//------------------------------------------------------------------------------

package main

import (
	"database/sql"
	"encoding/json"
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

// Struct to hold Logs_Data

type Loop_Data struct {
	Dtsp      float64 `json:"temp_sp"`
	Dtpv      float64 `json:"temp_pv"`
	Dhsp      float64 `json:"hum_sp"`
	Dhpv      float64 `json:"hum_pv"`
	Dpsp      float64 `json:"press_sp"`
	Dppv      float64 `json:"press_pv"`
	Ddatatime string  `json:"date_time_date"`
}

// Struct to hold Io_card_Info

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

	router.POST("/Logs_Event", processEvent_Log)
	router.POST("/Logs_Event_Type", processEvent_typeLog)
	router.POST("/Logs_Test", processTest_Log)
	router.POST("/Logs_Test_Type", processTest_typeLog)
	router.POST("/Logs_Maintenance", processMaintenance_Log)
	router.POST("/Loop_Data", processLoopDataInsert)
	router.PUT("/Loop_Data/:date_time_date", processLoopDataCreateOrUpdate)
	router.POST("/set_io_card_info", processIocardinfo)
}

func processEvent_Log(c *gin.Context) {

	//fmt.Println("Hello")

	var finalResult int = 1

	var log Logs_Event
	c.BindJSON(&log)
	//fmt.Println(log)

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

	// Activity log

	totaldata := map[string]interface{}{
		"log_id":                 log.Lid,
		"program_name":           log.Pname,
		"program_date_time":      log.Pdatetime,
		"ZTK_Logs_Event_Type_id": log.Etypeid,
		"ZTK_Users_id":           log.Eid,
		"created_by":             log.Createdby,
		"created":                log.Ecreated,
		"modified_by":            log.Modifiedby,
		"modified ":              log.Emodified,
	}
	datat, _ := json.Marshal(totaldata)
	newvalue := string(datat)
	stmt, err = db.Prepare("insert into ZTK_Activity_Log (`ZTK_Table_Id`, `action_type`, `new_value`, `ZTK_Users_Id`) values(?,?,?,?);")

	if err != nil {

		fmt.Print("Error: Creating Prepared Statement")

		fmt.Print(err.Error())

		finalResult = 0
	}

	_, err = stmt.Exec(5, "INSERT", newvalue, 1)

	if err != nil {

		fmt.Print("Error: Executing Prepared Statement")

		fmt.Print(err.Error())

		finalResult = 0
	}

	defer stmt.Close()

	if finalResult != 0 {

		c.JSON(http.StatusOK, gin.H{
			"Status = 1 ": fmt.Sprintf(" %s - Table id recorded.", 5),
			"Status = 2 ": fmt.Sprintf(" %s - action type recorded.", "INSERT"),
			"Status = 3 ": fmt.Sprintf(" %s - new value recorded.", newvalue),
			"Status = 4 ": fmt.Sprintf(" %s - User id recorded.", 1),
		})

	} else {

		c.JSON(http.StatusOK, gin.H{
			"Status = -1 ": fmt.Sprintf(" %s - Error of Table id.", 5),
			"Status = -2 ": fmt.Sprintf(" %s - Error of action type Log.", "INSERT"),
			"Status = -3 ": fmt.Sprintf(" %s - Error of new value Log.", newvalue),
			"Status = -4 ": fmt.Sprintf(" %s - Error ofUser idLog.", 1),
		})
	}
}

func processEvent_typeLog(c *gin.Context) {

	//fmt.Println("Hello")

	var finalResult int = 1

	var log Logs_Event_Type
	c.BindJSON(&log)
	//fmt.Println(log)

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

	// Activity log
	totaldata := map[string]interface{}{
		"events_type": log.Levents,
		"created":     log.Lcreated,
		"modified":    log.Lmodified,
		"created_by":  log.Lcreated1,
		"modified_by": log.Lmodified2,
	}
	datat, _ := json.Marshal(totaldata)
	newvalue := string(datat)
	stmt, err = db.Prepare("insert into ZTK_Activity_Log (`ZTK_Table_Id`, `action_type`, `new_value`, `ZTK_Users_Id`) values(?,?,?,?);")

	if err != nil {

		fmt.Print("Error: Creating Prepared Statement")

		fmt.Print(err.Error())

		finalResult = 0
	}

	_, err = stmt.Exec(5, "INSERT", newvalue, 1)

	if err != nil {

		fmt.Print("Error: Executing Prepared Statement")

		fmt.Print(err.Error())

		finalResult = 0
	}

	defer stmt.Close()

	if finalResult != 0 {

		c.JSON(http.StatusOK, gin.H{
			"Status = 1 ": fmt.Sprintf(" %s - Table id recorded.", 5),
			"Status = 2 ": fmt.Sprintf(" %s - action type recorded.", "INSERT"),
			"Status = 3 ": fmt.Sprintf(" %s - new value recorded.", newvalue),
			"Status = 4 ": fmt.Sprintf(" %s - User id recorded.", 1),
		})

	} else {

		c.JSON(http.StatusOK, gin.H{
			"Status = -1 ": fmt.Sprintf(" %s - Error of Table id.", 5),
			"Status = -2 ": fmt.Sprintf(" %s - Error of action type Log.", "INSERT"),
			"Status = -3 ": fmt.Sprintf(" %s - Error of new value Log.", newvalue),
			"Status = -4 ": fmt.Sprintf(" %s - Error ofUser idLog.", 1),
		})
	}
}

func processTest_Log(c *gin.Context) {

	//fmt.Println("Hello")

	var finalResult int = 1

	var log Logs_Test
	c.BindJSON(&log)
	//fmt.Println(log)

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

	// Activity log

	totaldata := map[string]interface{}{
		"log_id":                log.Tid,
		"log_name":              log.Tname,
		"log_date_time":         log.Tdatetime,
		"ZTK_Logs_Test_Type_id": log.Ttypeid,
		"ZTK_Users_id":          log.Tuserid,
		"created_by":            log.Tcreatedby,
		"created":               log.Tcreated,
		"modified_by":           log.Tmodifiedby,
		"modified ":             log.Tmodified,
	}
	datat, _ := json.Marshal(totaldata)
	newvalue := string(datat)
	stmt, err = db.Prepare("insert into ZTK_Activity_Log (`ZTK_Table_Id`, `action_type`, `new_value`, `ZTK_Users_Id`) values(?,?,?,?);")

	if err != nil {

		fmt.Print("Error: Creating Prepared Statement")

		fmt.Print(err.Error())

		finalResult = 0
	}

	_, err = stmt.Exec(5, "INSERT", newvalue, 1)

	if err != nil {

		fmt.Print("Error: Executing Prepared Statement")

		fmt.Print(err.Error())

		finalResult = 0
	}

	defer stmt.Close()

	if finalResult != 0 {

		c.JSON(http.StatusOK, gin.H{
			"Status = 1 ": fmt.Sprintf(" %s - Table id recorded.", 5),
			"Status = 2 ": fmt.Sprintf(" %s - action type recorded.", "INSERT"),
			"Status = 3 ": fmt.Sprintf(" %s - new value recorded.", newvalue),
			"Status = 4 ": fmt.Sprintf(" %s - User id recorded.", 1),
		})

	} else {

		c.JSON(http.StatusOK, gin.H{
			"Status = -1 ": fmt.Sprintf(" %s - Error of Table id.", 5),
			"Status = -2 ": fmt.Sprintf(" %s - Error of action type Log.", "INSERT"),
			"Status = -3 ": fmt.Sprintf(" %s - Error of new value Log.", newvalue),
			"Status = -4 ": fmt.Sprintf(" %s - Error ofUser idLog.", 1),
		})
	}
}

func processTest_typeLog(c *gin.Context) {

	//fmt.Println("Hello")

	var finalResult int = 1

	var log Logs_Test_Type
	c.BindJSON(&log)
	//fmt.Println(log)

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
	// Activity log

	totaldata := map[string]interface{}{
		"test_type":   log.Ltesttype,
		"created":     log.Tcreated1,
		"modified":    log.Tmodified2,
		"created_by":  log.Tcreatedby1,
		"modified_by": log.Tmodifiedby2,
	}
	datat, _ := json.Marshal(totaldata)
	newvalue := string(datat)
	stmt, err = db.Prepare("insert into ZTK_Activity_Log (`ZTK_Table_Id`, `action_type`, `new_value`, `ZTK_Users_Id`) values(?,?,?,?);")

	if err != nil {

		fmt.Print("Error: Creating Prepared Statement")

		fmt.Print(err.Error())

		finalResult = 0
	}

	_, err = stmt.Exec(5, "INSERT", newvalue, 1)

	if err != nil {

		fmt.Print("Error: Executing Prepared Statement")

		fmt.Print(err.Error())

		finalResult = 0
	}

	defer stmt.Close()

	if finalResult != 0 {

		c.JSON(http.StatusOK, gin.H{
			"Status = 1 ": fmt.Sprintf(" %s - Table id recorded.", 5),
			"Status = 2 ": fmt.Sprintf(" %s - action type recorded.", "INSERT"),
			"Status = 3 ": fmt.Sprintf(" %s - new value recorded.", newvalue),
			"Status = 4 ": fmt.Sprintf(" %s - User id recorded.", 1),
		})

	} else {

		c.JSON(http.StatusOK, gin.H{
			"Status = -1 ": fmt.Sprintf(" %s - Error of Table id.", 5),
			"Status = -2 ": fmt.Sprintf(" %s - Error of action type Log.", "INSERT"),
			"Status = -3 ": fmt.Sprintf(" %s - Error of new value Log.", newvalue),
			"Status = -4 ": fmt.Sprintf(" %s - Error ofUser idLog.", 1),
		})
	}

}

func processMaintenance_Log(c *gin.Context) {

	//fmt.Println("Hello")

	var finalResult int = 1

	var log Logs_Maintenance
	c.BindJSON(&log)
	//fmt.Println(log)

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
	// Activity log

	totaldata := map[string]interface{}{
		"component_name":      log.Mname,
		"runtime_hr":          log.Mruntime,
		"counter":             log.Mcounter,
		"days_till_service":   log.Mservice,
		"maintenance_pending": log.Mpending,
		"maintenance_status":  log.Mstatus,
		"created":             log.Mcreated,
		"modified":            log.Mmodified,
		"created_by ":         log.Mcreatedby,
		"modified_by ":        log.Mmodifiedby,
	}
	datat, _ := json.Marshal(totaldata)
	newvalue := string(datat)
	stmt, err = db.Prepare("insert into ZTK_Activity_Log (`ZTK_Table_Id`, `action_type`, `new_value`, `ZTK_Users_Id`) values(?,?,?,?);")

	if err != nil {

		fmt.Print("Error: Creating Prepared Statement")

		fmt.Print(err.Error())

		finalResult = 0
	}

	_, err = stmt.Exec(5, "INSERT", newvalue, 1)

	if err != nil {

		fmt.Print("Error: Executing Prepared Statement")

		fmt.Print(err.Error())

		finalResult = 0
	}

	defer stmt.Close()

	if finalResult != 0 {

		c.JSON(http.StatusOK, gin.H{
			"Status = 1 ": fmt.Sprintf(" %s - Table id recorded.", 5),
			"Status = 2 ": fmt.Sprintf(" %s - action type recorded.", "INSERT"),
			"Status = 3 ": fmt.Sprintf(" %s - new value recorded.", newvalue),
			"Status = 4 ": fmt.Sprintf(" %s - User id recorded.", 1),
		})

	} else {

		c.JSON(http.StatusOK, gin.H{
			"Status = -1 ": fmt.Sprintf(" %s - Error of Table id.", 5),
			"Status = -2 ": fmt.Sprintf(" %s - Error of action type Log.", "INSERT"),
			"Status = -3 ": fmt.Sprintf(" %s - Error of new value Log.", newvalue),
			"Status = -4 ": fmt.Sprintf(" %s - Error ofUser idLog.", 1),
		})
	}
}

func processLoopDataCreateOrUpdate(c *gin.Context) {

	dateTime := c.Params.ByName("date_time_date")
	var count int
	fmt.Println(dateTime)

	// step 1 check record is exist or not with dataTime value
	stmt, err := db.Prepare("select count(id) from ZTK_Loop_Data where date_time = ?")
	if err != nil {
		fmt.Print(err.Error())
	}
	row := stmt.QueryRow(dateTime)
	row.Scan(&count)
	if err != nil {
		fmt.Print(err.Error())
	}
	fmt.Println(count)
	if count == 0 {
		fmt.Println("record is not exists, need to create")
		processLoopDataInsert(c)
	} else {
		fmt.Println("record is  exists, need to update")
		processLoopDataUpdate(c)
	}

}

func processLoopDataUpdate(c *gin.Context) {

	var finalResult int = 1

	var log Loop_Data
	c.BindJSON(&log)
	fmt.Println(log)

	stmt, err := db.Prepare("UPDATE ZTK_Loop_Data SET temp_sp=?,temp_pv=?,hum_sp=?,hum_pv=?,press_sp=?,press_pv=? WHERE date_time= ? ")

	if err != nil {

		fmt.Print("Error: Creating Prepared Statement")

		fmt.Print(err.Error())

		finalResult = 0
	}

	_, err = stmt.Exec(log.Dtsp, log.Dtpv, log.Dhsp, log.Dhpv, log.Dpsp, log.Dppv, log.Ddatatime)

	if err != nil {

		fmt.Print("Error: Executing Prepared Statement")

		fmt.Print(err.Error())

		finalResult = 0
	}

	defer stmt.Close()

	if finalResult != 0 {

		c.JSON(http.StatusOK, gin.H{

			"Status = 1 ": fmt.Sprintf(" %s - Dtsp  Log recorded.", log.Dtsp),
			"Status = 2 ": fmt.Sprintf(" %s - Dtpv  Log recorded.", log.Dtpv),
			"Status = 3 ": fmt.Sprintf(" %s - Dhsp  Log recorded.", log.Dhsp),
			"Status = 4 ": fmt.Sprintf(" %s - Dhpv Log recorded.", log.Dhpv),
			"Status = 5 ": fmt.Sprintf(" %s - Dpsp  Log recorded.", log.Dpsp),
			"Status = 6 ": fmt.Sprintf(" %s - Dppv  Log recorded.", log.Dppv),
			"Status = 7 ": fmt.Sprintf(" %s - DdatatimeUpdate  Log recorded.", log.Ddatatime),
		})

	} else {

		c.JSON(http.StatusOK, gin.H{

			"Status = -1 ": fmt.Sprintf(" %s - Error of Dtsp Log.", log.Dtsp),
			"Status = -2 ": fmt.Sprintf(" %s - Error of Dtpv Log.", log.Dtpv),
			"Status = -3 ": fmt.Sprintf(" %s - Error of Dhsp Log.", log.Dhsp),
			"Status = -4 ": fmt.Sprintf(" %s - Error of Dhpv Log.", log.Dhpv),
			"Status = -5 ": fmt.Sprintf(" %s - Error of Dpsp Log.", log.Dpsp),
			"Status = -6 ": fmt.Sprintf(" %s - Error of Dppv Log.", log.Dppv),
			"Status = -7 ": fmt.Sprintf(" %s - Error of DdatatimeUpdate Log.", log.Ddatatime),
		})
	}
}

func processLoopDataInsert(c *gin.Context) {

	//fmt.Println("Hello")

	var finalResult int = 1

	var log Loop_Data
	c.BindJSON(&log)
	fmt.Println(log)

	stmt, err := db.Prepare("insert into ZTK_Loop_Data (temp_sp,temp_pv,hum_sp,hum_pv,press_sp,press_pv,date_time ) values(?,?,?,?,?,?,?);")

	if err != nil {

		fmt.Print("Error: Creating Prepared Statement")

		fmt.Print(err.Error())

		finalResult = 0
	}

	_, err = stmt.Exec(log.Dtsp, log.Dtpv, log.Dhsp, log.Dhpv, log.Dpsp, log.Dppv, log.Ddatatime)

	if err != nil {

		fmt.Print("Error: Executing Prepared Statement")

		fmt.Print(err.Error())

		finalResult = 0
	}

	defer stmt.Close()

	if finalResult != 0 {

		c.JSON(http.StatusOK, gin.H{
			"Status = 1 ": fmt.Sprintf(" %s - Dtsp  Log recorded.", log.Dtsp),
			"Status = 2 ": fmt.Sprintf(" %s - Dtpv  Log recorded.", log.Dtpv),
			"Status = 3 ": fmt.Sprintf(" %s - Dhsp  Log recorded.", log.Dhsp),
			"Status = 4 ": fmt.Sprintf(" %s - Dhpv Log recorded.", log.Dhpv),
			"Status = 5 ": fmt.Sprintf(" %s - Dpsp  Log recorded.", log.Dpsp),
			"Status = 6 ": fmt.Sprintf(" %s - Dppv  Log recorded.", log.Dppv),
			"Status = 7 ": fmt.Sprintf(" %s - Ddatatime  Log recorded.", log.Ddatatime),
		})

	} else {

		c.JSON(http.StatusOK, gin.H{
			"Status = -1 ": fmt.Sprintf(" %s - Error of Dtsp Log.", log.Dtsp),
			"Status = -2 ": fmt.Sprintf(" %s - Error of Dtpv Log.", log.Dtpv),
			"Status = -3 ": fmt.Sprintf(" %s - Error of Dhsp Log.", log.Dhsp),
			"Status = -4 ": fmt.Sprintf(" %s - Error of Dhpv Log.", log.Dhpv),
			"Status = -5 ": fmt.Sprintf(" %s - Error of Dpsp Log.", log.Dpsp),
			"Status = -6 ": fmt.Sprintf(" %s - Error of Dppv Log.", log.Dppv),
			"Status = -7 ": fmt.Sprintf(" %s - Error of Ddatatime Log.", log.Ddatatime),
		})
	}
	// Activity log

	totaldata := map[string]interface{}{
		"temp_sp":   log.Dtsp,
		"temp_pv":   log.Dtpv,
		"hum_sp":    log.Dhsp,
		"hum_pv":    log.Dhpv,
		"press_sp":  log.Dpsp,
		"press_pv":  log.Dppv,
		"date_time": log.Ddatatime,
	}
	datat, _ := json.Marshal(totaldata)
	newvalue := string(datat)
	stmt, err = db.Prepare("insert into ZTK_Activity_Log (`ZTK_Table_Id`, `action_type`, `new_value`, `ZTK_Users_Id`) values(?,?,?,?);")

	if err != nil {

		fmt.Print("Error: Creating Prepared Statement")

		fmt.Print(err.Error())

		finalResult = 0
	}

	_, err = stmt.Exec(5, "INSERT", newvalue, 1)

	if err != nil {

		fmt.Print("Error: Executing Prepared Statement")

		fmt.Print(err.Error())

		finalResult = 0
	}

	defer stmt.Close()

	if finalResult != 0 {

		c.JSON(http.StatusOK, gin.H{
			"Status = 1 ": fmt.Sprintf(" %s - Table id recorded.", 5),
			"Status = 2 ": fmt.Sprintf(" %s - action type recorded.", "INSERT"),
			"Status = 3 ": fmt.Sprintf(" %s - new value recorded.", newvalue),
			"Status = 4 ": fmt.Sprintf(" %s - User id recorded.", 1),
		})

	} else {

		c.JSON(http.StatusOK, gin.H{
			"Status = -1 ": fmt.Sprintf(" %s - Error of Table id.", 5),
			"Status = -2 ": fmt.Sprintf(" %s - Error of action type Log.", "INSERT"),
			"Status = -3 ": fmt.Sprintf(" %s - Error of new value Log.", newvalue),
			"Status = -4 ": fmt.Sprintf(" %s - Error ofUser idLog.", 1),
		})
	}
}

func processIocardinfo(c *gin.Context) {

	//fmt.Println("Hello")

	var finalResult int = 1

	var log Io_card_Info
	c.BindJSON(&log)
	fmt.Println(log)

	stmt, err := db.Prepare("insert into ZTK_IO_Card_Info (card_address,card_type,card_version,card_serial_number,secret_key,customer_id,mfg_date,created,modified,created_by,modified_by ) values(?,?,?,?,?,?,?,?,?,?,?);")

	if err != nil {

		fmt.Print("Error: Creating Prepared Statement")

		fmt.Print(err.Error())

		finalResult = 0
	}

	_, err = stmt.Exec(log.Iaddress, log.Itype, log.Iversion, log.Inumber, log.Ikey, log.Iid, log.Idate, log.Icreated, log.Imodified, log.Icreatedby, log.Imodifiedby)

	if err != nil {

		fmt.Print("Error: Executing Prepared Statement")

		fmt.Print(err.Error())

		finalResult = 0
	}

	defer stmt.Close()

	if finalResult != 0 {

		c.JSON(http.StatusOK, gin.H{
			"Status = 1 ":   fmt.Sprintf(" %s - address  Log recorded.", log.Iaddress),
			"Status = 2 ":   fmt.Sprintf(" %s - type  Log recorded.", log.Itype),
			"Status = 3 ":   fmt.Sprintf(" %s - version  Log recorded.", log.Iversion),
			"Status = 4 ":   fmt.Sprintf(" %s - number Log recorded.", log.Inumber),
			"Status = 5 ":   fmt.Sprintf(" %s - key  Log recorded.", log.Ikey),
			"Status = 6 ":   fmt.Sprintf(" %s - id  Log recorded.", log.Iid),
			"Status = 7 ":   fmt.Sprintf(" %s - date Log recorded.", log.Idate),
			"Status = 8 ":   fmt.Sprintf(" %s - created  Log recorded.", log.Icreated),
			"Status = 9 ":   fmt.Sprintf(" %s - Modified  Log recorded.", log.Imodified),
			"Status = 10 ":  fmt.Sprintf(" %s - createdby  Log recorded.", log.Icreatedby),
			"Status = 11  ": fmt.Sprintf(" %s - Modifiedby  Log recorded.", log.Imodifiedby),
		})

	} else {

		c.JSON(http.StatusOK, gin.H{
			"Status = -1 ":  fmt.Sprintf(" %s - Error of address Log.", log.Iaddress),
			"Status = -2 ":  fmt.Sprintf(" %s - Error of type Log.", log.Itype),
			"Status = -3 ":  fmt.Sprintf(" %s - Error of version Log.", log.Iversion),
			"Status = -4 ":  fmt.Sprintf(" %s - Error of number Log.", log.Inumber),
			"Status = -5 ":  fmt.Sprintf(" %s - Error of key Log.", log.Ikey),
			"Status = -6 ":  fmt.Sprintf(" %s - Error of id Log.", log.Iid),
			"Status = -7 ":  fmt.Sprintf(" %s - Error of date Log.", log.Idate),
			"Status = -8 ":  fmt.Sprintf(" %s - Error of Created Log.", log.Icreated),
			"Status = -9":   fmt.Sprintf(" %s - Error of Modified Log.", log.Imodified),
			"Status = -10":  fmt.Sprintf(" %s - Error of createdby Log.", log.Icreatedby),
			"Status = -11 ": fmt.Sprintf(" %s - Error of Modifiedby Log.", log.Imodifiedby),
		})
	}
	// Activity log

	totaldata := map[string]interface{}{
		"card_address":       log.Iaddress,
		"card_type":          log.Itype,
		"card_version":       log.Iversion,
		"card_serial_number": log.Inumber,
		"secret_key":         log.Ikey,
		"customer_id":        log.Iid,
		"mfg_date":           log.Idate,
		"created":            log.Icreated,
		"modified ":          log.Imodified,
		"created_by ":        log.Icreatedby,
		"modified_by ":       log.Imodifiedby,
	}
	datat, _ := json.Marshal(totaldata)
	newvalue := string(datat)
	stmt, err = db.Prepare("insert into ZTK_Activity_Log (`ZTK_Table_Id`, `action_type`, `new_value`, `ZTK_Users_Id`) values(?,?,?,?);")

	if err != nil {

		fmt.Print("Error: Creating Prepared Statement")

		fmt.Print(err.Error())

		finalResult = 0
	}

	_, err = stmt.Exec(5, "INSERT", newvalue, 1)

	if err != nil {

		fmt.Print("Error: Executing Prepared Statement")

		fmt.Print(err.Error())

		finalResult = 0
	}

	defer stmt.Close()

	if finalResult != 0 {

		c.JSON(http.StatusOK, gin.H{
			"Status = 1 ": fmt.Sprintf(" %s - Table id recorded.", 5),
			"Status = 2 ": fmt.Sprintf(" %s - action type recorded.", "INSERT"),
			"Status = 3 ": fmt.Sprintf(" %s - new value recorded.", newvalue),
			"Status = 4 ": fmt.Sprintf(" %s - User id recorded.", 1),
		})

	} else {

		c.JSON(http.StatusOK, gin.H{
			"Status = -1 ": fmt.Sprintf(" %s - Error of Table id.", 5),
			"Status = -2 ": fmt.Sprintf(" %s - Error of action type Log.", "INSERT"),
			"Status = -3 ": fmt.Sprintf(" %s - Error of new value Log.", newvalue),
			"Status = -4 ": fmt.Sprintf(" %s - Error ofUser idLog.", 1),
		})
	}
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

	connectStr += strconv.Itoa(ngcsLogConfig.LocalLogServerPort)

	return connectStr
}
