package dao

import (
	"runtime"
	"testing"

	"github.com/SoftwareUndagi/golibs/common"

	"github.com/sirupsen/logrus"

	"github.com/SoftwareUndagi/golibs/testhelper"
)

func TestRunSQLFile(t *testing.T) {

	db, logCapture, err := testhelper.MysqlTestAppSetup(t)
	logentry := logrus.WithField("testMethod", "TestRunSQLFile")
	(*logCapture).SetJSONFormatLog("/tmp/dodol.json")
	if err != nil {
		t.Errorf("Error open connection , error: %s", err.Error())
		return
	}
	_, filename, _, _ := runtime.Caller(0)
	directory := common.GetFileDirectory(filename)
	sqlFile := common.AppendFilePath(directory, "sql-file-executor-mysql_test.sql")
	err = RunSQLFile(db, logentry, sqlFile, true)
	if err != nil {
		logentry.WithError(err).WithFields(logrus.Fields{"path": sqlFile}).Errorf("Fail to invoke , error: %s", err.Error())
	}

}
