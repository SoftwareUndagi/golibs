package dao

import (
	"testing"

	"github.com/SoftwareUndagi/golibs/testhelper"
	"github.com/sirupsen/logrus"
)

func TestUsageCountIdString32Mysql(t *testing.T) {
	db, _, err := testhelper.MysqlTestAppSetup(t)
	if err != nil {
		t.Errorf("Error open connection , error: %s", err.Error())
		return
	}
	db.LogMode(true)
	logEntry := logrus.WithField("method", "TestUsageCountIdInt64Mysql")
	err = testhelper.RunDbInitializationScripts(db, logEntry, TempRoleCreateMysql,
		TempUserCreateMysql, TempUserRoleCreateMysql, TempUserInsertMySQL, TempRoleInsertMySQL, TempUserRoleInsertMySQL)
	if err != nil {
		logEntry.Errorf("Init failed , error: %s", err.Error())
		return
	}
	//db.LogMode(true)
	defer testhelper.DropMysqlTempTables(db, logEntry, TempRoleTableName, TempUserTempTableName, TempUserRoleTableName)
	DefaultDaoManager.RegisterUsageCountQuery("mysql", "sec_role", `select 'sec_user' modelName ,  role_code   dataID , count(*) usageCount  from `+TempUserRoleTableName+` where user_id in ( :ids) `)
	cnt, errCount := GetUsageCountIDString(db, logEntry, "sec_role", "ROLE001")
	if errCount != nil {
		logEntry.WithError(errCount).Errorf("Fail to run usage count , error: %s", errCount.Error())
		return
	}
	if cnt != 2 {
		t.Errorf("Usage seharusnya 2, result : %d", cnt)
		return
	}
	logEntry.Infof("Selesai, usage: %d", cnt)
}
