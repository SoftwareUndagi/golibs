package coredao

import (
	"fmt"
	"runtime"
	"testing"

	"github.com/SoftwareUndagi/golibs/common"
	"github.com/SoftwareUndagi/golibs/dao"
	"github.com/SoftwareUndagi/golibs/testhelper"
	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"
)

type bulkInsertStudy struct {
	ID                   int64  `gorm:"column:id;primary_key;AUTO_INCREMENT"`
	Name                 string `gorm:"column:name"`
	Email                string `gorm:"column:email"`
	DefaultColumn        string `gorm:"DEFAULT;column:def_col"`
	TagUserName          string `gorm:"column:tag_username"`
	AfterCreateSampleTag string `gorm:"-"`
}

type logrusDbLogger struct {
	logEntry *logrus.Entry
}

func (p *logrusDbLogger) Print(v ...interface{}) {
	if len(v) > 0 {
		p.logEntry.Infof("SQL: %v", v[0])
	}

}

func (p *bulkInsertStudy) TableName(db *gorm.DB) string {
	return "bulk_insert_study"
}
func (p *bulkInsertStudy) BeforeCreate(scope *gorm.Scope) (err error) {
	usrTag, ok := scope.Get("username")
	if ok {
		p.TagUserName = usrTag.(string)
	} else {
		p.TagUserName = "unknown"
	}
	return nil
}
func (p *bulkInsertStudy) AfterCreate(scope *gorm.Scope) (err error) {
	p.AfterCreateSampleTag = fmt.Sprintf("%d-%s", p.ID, p.Name)
	return nil
}

func Test_BulkCreate(t *testing.T) {
	db, logCapture, err := testhelper.MysqlTestAppSetup(t)
	if err != nil {
		t.Fail()
		t.Errorf("Gagal buka koneksi")
		return
	}
	logentry := logrus.WithField("testMethod", "Test_BulkCreate")
	_, filename, _, _ := runtime.Caller(0)
	directory := common.GetFileDirectory(filename)
	sqlFile := common.AppendFilePath(directory, "BulkCreate_test-mysql.sql")
	db.LogMode(true)
	err = dao.RunSQLFile(db, logentry, sqlFile, true)
	if err != nil {
		t.Error(err.Error())
		logentry.WithError(err).WithFields(logrus.Fields{"path": sqlFile}).Errorf("Fail to invoke , error: %s", err.Error())
		t.Fail()
		return
	}

	logPath := "/tmp/go-test-Test_BulkCreate.json"
	(*logCapture).SetJSONFormatLog(logPath)
	smpl1 := bulkInsertStudy{Name: "Gede", Email: "gede.sutarsa@gmail.com"}
	smpl2 := bulkInsertStudy{Name: "Yuli", Email: "yuli.dana@gmail.com"}
	smpl3 := bulkInsertStudy{Name: "Indra", Email: "indra@gmail.com"}
	smpl4 := bulkInsertStudy{Name: "Tasya", Email: "tasya@gmail.com"}
	smpl5 := bulkInsertStudy{Name: "nadya", Email: "nadya@gmail.com"}
	smpl6 := bulkInsertStudy{Name: "rani", Email: "rani@gmail.com"}
	wrapper := []interface{}{&smpl1, &smpl2, &smpl3, &smpl4, &smpl5, &smpl6}

	db.InstantSet("username", "gede.sutarsa")
	dbLogger := logrusDbLogger{logEntry: logentry}
	db.SetLogger(&dbLogger)
	dbRslt := db.BulkCreate(wrapper)
	if dbRslt.Error != nil {
		logentry.WithError(dbRslt.Error).Errorf("Fail to bulk insert, error : %s", dbRslt.Error.Error())
		t.Fail()
		return
	}
	logentry.WithFields(logrus.Fields{"affectedRows": dbRslt.RowsAffected, "data": wrapper}).Infof("Selesai")
}
