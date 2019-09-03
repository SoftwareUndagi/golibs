package dao

import (
	"testing"

	"github.com/SoftwareUndagi/golibs/testhelper"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/sirupsen/logrus"
)

func TestUsageCountIdInt16Mysql(t *testing.T) {
	db, _, err := testhelper.MysqlTestAppSetup(t)

	if err != nil {
		t.Errorf("Error open connection , error: %s", err.Error())
		return
	}
	db.LogMode(false)
	logEntry := logrus.WithField("method", "TestUsageCountIdInt16Mysql")
	err = testhelper.RunDbInitializationScripts(db, logEntry, TempGroupTableCreateMySQL, TempGroupAssignmentCreateMysql, TempGroupInsertSQL, TempGroupAssignInsertMySQL)
	if err != nil {
		logEntry.Errorf("Init failed , error: %s", err.Error())
		return
	}
	db.LogMode(true)
	defer testhelper.DropMysqlTempTables(db, logEntry, SecGroupTempTableName, SecGroupAssignmentTempTableName)
	DefaultDaoManager.RegisterUsageCountQuery("mysql", "sec_group_temp", `select 'sec_group_temp' modelName ,  group_id   dataID , count(*) usageCount  from `+SecGroupAssignmentTempTableName+` where group_id in ( :ids) `)
	cnt, errCount := GetUsageCountIDInt16(db, logEntry, "sec_group_temp", 1)
	if errCount != nil {
		logEntry.WithError(errCount).Errorf("Fail to run usage count , error: %s", errCount.Error())
		return
	}
	if cnt != 2 {
		t.Errorf("Usage seharusnya 2")
		return
	}
	logEntry.Infof("Selesai, usage: %d", cnt)

}

func TestUsageCountIdInt32Mysql(t *testing.T) {
	db, _, err := testhelper.MysqlTestAppSetup(t)
	if err != nil {
		t.Errorf("Error open connection , error: %s", err.Error())
		return
	}
	db.LogMode(false)
	logEntry := logrus.WithField("method", "TestUsageCountIdInt32Mysql")
	err = testhelper.RunDbInitializationScripts(db, logEntry,
		TempGroupTableCreateMySQL, TempGroupAssignmentCreateMysql, TempUserCreateMysql,
		TempUserGroupAssignmentCreateMySQL, TempUserRoleCreateMysql,
		TempGroupInsertSQL, TempGroupAssignInsertMySQL, TempUserInsertMySQL, TempUserAssignmentInsertMySQL, TempUserRoleInsertMySQL)
	if err != nil {
		logEntry.Errorf("Init failed , error: %s", err.Error())
		return
	}
	db.LogMode(true)
	defer testhelper.DropMysqlTempTables(db, logEntry, SecGroupTempTableName, SecGroupAssignmentTempTableName, TempUserGroupAssignmentTableName, TempUserTempTableName, TempUserRoleTableName)
	DefaultDaoManager.RegisterUsageCountQuery("mysql", "sec_user", `select 'sec_user' modelName ,  user_id   dataID , count(*) usageCount  from `+TempUserGroupAssignmentTableName+` where user_id in ( :ids) `)
	DefaultDaoManager.RegisterUsageCountQuery("mysql", "sec_user", `select 'sec_user' modelName ,  user_id   dataID , count(*) usageCount  from `+TempUserRoleTableName+` where user_id in ( :ids) `)
	cnt, errCount := GetUsageCountIDInt32(db, logEntry, "sec_user", 1)
	if errCount != nil {
		logEntry.WithError(errCount).Errorf("Fail to run usage count , error: %s", errCount.Error())
		return
	}
	if cnt != 5 {
		t.Errorf("Usage seharusnya 5, result : %d", cnt)
		return
	}
	logEntry.Infof("Selesai, usage: %d", cnt)
}

func TestUsageCountIdInt64Mysql(t *testing.T) {
	db, _, err := testhelper.MysqlTestAppSetup(t)
	if err != nil {
		t.Errorf("Error open connection , error: %s", err.Error())
		return
	}
	db.LogMode(false)
	logEntry := logrus.WithField("method", "TestUsageCountIdInt64Mysql")
	err = testhelper.RunDbInitializationScripts(db, logEntry,
		TempGroupTableCreateMySQL, TempGroupAssignmentCreateMysql, TempUserCreateMysql,
		TempUserGroupAssignmentCreateMySQL, TempUserRoleCreateMysql,
		TempGroupInsertSQL, TempGroupAssignInsertMySQL, TempUserInsertMySQL, TempUserAssignmentInsertMySQL, TempUserRoleInsertMySQL)
	if err != nil {
		logEntry.Errorf("Init failed , error: %s", err.Error())
		return
	}
	db.LogMode(true)
	defer testhelper.DropMysqlTempTables(db, logEntry, SecGroupTempTableName, SecGroupAssignmentTempTableName, TempUserGroupAssignmentTableName, TempUserTempTableName, TempUserRoleTableName)
	DefaultDaoManager.RegisterUsageCountQuery("mysql", "sec_user", `select 'sec_user' modelName ,  user_id   dataID , count(*) usageCount  from `+TempUserGroupAssignmentTableName+` where user_id in ( :ids) `)
	DefaultDaoManager.RegisterUsageCountQuery("mysql", "sec_user", `select 'sec_user' modelName ,  user_id   dataID , count(*) usageCount  from `+TempUserRoleTableName+` where user_id in ( :ids) `)
	cnt, errCount := GetUsageCountIDInt64(db, logEntry, "sec_user", 1)
	if errCount != nil {
		logEntry.WithError(errCount).Errorf("Fail to run usage count , error: %s", errCount.Error())
		return
	}
	if cnt != 5 {
		t.Errorf("Usage seharusnya 5, result : %d", cnt)
		return
	}
	logEntry.Infof("Selesai, usage: %d", cnt)
}
