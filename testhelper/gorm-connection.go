package testhelper

import (
	"fmt"
	"os"
	"testing"

	"github.com/sirupsen/logrus"

	"github.com/SoftwareUndagi/golibs/common"

	"github.com/jinzhu/gorm"
)

//OpenMysqlConnection open mysql database connection. using host + port
func OpenMysqlConnection(username string, password string, hostname string, dbName string, portString string) (db *gorm.DB, err error) {
	/*
		conQuery := os.Getenv("dbUser") + ":" + os.Getenv("dbPassword")
		dbPort := os.Getenv("dbPort")
		dbHostName := os.Getenv("dbHostName")
		socketPath := os.Getenv("dbSocketPath")
		dbSchema := os.Getenv("dbName")
	*/
	if len(portString) == 0 {
		portString = "3306"
	}
	conQuery := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", username, password, hostname, portString, dbName)
	return gorm.Open("mysql", conQuery)
}

//RunDbInitializationScripts run initialization scripts. create table or else
func RunDbInitializationScripts(db *gorm.DB, logEntry *logrus.Entry, setupSQL ...string) (err error) {
	if len(setupSQL) == 0 {
		return
	}
	for index, sql := range setupSQL {
		if rslt := db.Exec(sql); rslt.Error != nil {
			err = rslt.Error
			logEntry.WithField("index", index).WithField("sql", sql).WithError(err).Errorf("Fail to run setup script. error: %s", err.Error())
			return
		}
	}
	return
}

//DropMysqlTempTables statement drop tables
func DropMysqlTempTables(db *gorm.DB, logEntry *logrus.Entry, tableNames ...string) (err error) {
	if len(tableNames) == 0 {
		return
	}
	for _, t := range tableNames {
		dropSQL := fmt.Sprintf("drop temporary table %s", t)
		db.Exec(dropSQL)
	}
	return
}

//OpenPostgreConnection open porgres sql datbase connection. if portString is 0 length, default posgtres port will be used
func OpenPostgreConnection(username string, password string, hostname string, dbName string, portString string) (db *gorm.DB, err error) {
	if len(portString) == 0 {
		portString = "5432"
	}
	conQuery := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s", hostname, portString, username, dbName, password)
	db, err = gorm.Open("postgres", conQuery)
	return
}

//OpenMysqlConnectionWithSocketPath open mysql connection with unix socket path
func OpenMysqlConnectionWithSocketPath(username string, password string, dbName string, socketPath string) (db *gorm.DB, err error) {
	conQuery := fmt.Sprintf("%s:%s@unix(%s)/%s?charset=utf8&parseTime=True&loc=Local", username, password, socketPath, dbName)
	return gorm.Open("mysql", conQuery)
}

//PosgresqlTestAppSetup setup postgres for test env. parameter read from os parameter
func PosgresqlTestAppSetup(t *testing.T) (db *gorm.DB, err error) {
	username := os.Getenv("postgres.dbUser")
	password := os.Getenv("postgres.dbPassword")
	hostname := os.Getenv("postgres.dbHostName")
	dbName := os.Getenv("postgres.dbName")
	loggerUsername := os.Getenv("username")
	ipAddress := os.Getenv("ipaddress")
	if len(loggerUsername) == 0 {
		loggerUsername = "coolMeUser"
	}
	if len(ipAddress) == 0 {
		ipAddress = "127.0.0.1"
	}
	portString := os.Getenv("postgres.dbPort")
	var errorMsgs []string
	if len(username) == 0 {
		errorMsgs = append(errorMsgs, "Parameter key : postgres.dbUser not set on os.getEnv, please set the variable first")
	}
	if len(password) == 0 {
		errorMsgs = append(errorMsgs, "Parameter key : postgres.dbHostName not set on os.getEnv, please set the variable first")
	}
	if len(hostname) == 0 {
		errorMsgs = append(errorMsgs, "Parameter key : postgres.hostname not set on os.getEnv, please set the variable first")
	}
	if len(dbName) == 0 {
		errorMsgs = append(errorMsgs, "Parameter key : postgres.dbName not set on os.getEnv, please set the variable first")
	}
	db, err = OpenPostgreConnection(username, password, hostname, dbName, portString)
	if err != nil {
		t.Errorf("Error opening Postgres connection , error: %s", err.Error())
		return
	}
	db.LogMode(true)
	db.InstantSet(common.GormVariableUsername, loggerUsername)
	db.InstantSet(common.GormVariableIPAddress, ipAddress)
	CaptureLog(t).Release()
	return
}

//MysqlTestAppSetup run setup testing mysql
func MysqlTestAppSetup(t *testing.T) (db *gorm.DB, err error) {
	username := os.Getenv("mysql.dbUser")
	password := os.Getenv("mysql.dbPassword")
	hostname := os.Getenv("mysql.dbHostName")
	dbName := os.Getenv("mysql.dbName")
	loggerUsername := os.Getenv("username")
	ipAddress := os.Getenv("ipaddress")
	if len(loggerUsername) == 0 {
		loggerUsername = "coolMeUser"
	}
	if len(ipAddress) == 0 {
		ipAddress = "127.0.0.1"
	}
	portString := os.Getenv("mysql.dbPort")
	var errorMsgs []string
	if len(username) == 0 {
		errorMsgs = append(errorMsgs, "Parameter key : mysql.dbUser not set on os.getEnv, please set the variable first")
	}
	if len(password) == 0 {
		errorMsgs = append(errorMsgs, "Parameter key : mysql.dbHostName not set on os.getEnv, please set the variable first")
	}
	if len(hostname) == 0 {
		errorMsgs = append(errorMsgs, "Parameter key : mysql.hostname not set on os.getEnv, please set the variable first")
	}
	if len(dbName) == 0 {
		errorMsgs = append(errorMsgs, "Parameter key : mysql.dbName not set on os.getEnv, please set the variable first")
	}
	db, err = OpenMysqlConnection(username, password, hostname, dbName, portString)
	if err != nil {
		t.Errorf("Error opening mysql connection , error: %s", err.Error())
		return
	}
	db.LogMode(true)
	db.InstantSet(common.GormVariableUsername, loggerUsername)
	db.InstantSet(common.GormVariableIPAddress, ipAddress)
	CaptureLog(t).Release()
	return
}
