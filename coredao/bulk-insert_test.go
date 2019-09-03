package coredao

import (
	"testing"

	"github.com/SoftwareUndagi/golibs/dao"
	"github.com/SoftwareUndagi/golibs/testhelper"
	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"
)

type dummyBulkData struct {
	//ID id dari data
	ID int64 `gorm:"column:id;AUTO_INCREMENT"`
	//Name name of
	Name string `gorm:"column:name"`
	//Email email of user
	Email string `gorm:"column:email"`
}

func Test_bulk_insert(t *testing.T) {
	db, logCapture, err := testhelper.MysqlTestAppSetup(t)
	if err != nil {
		t.Fail()
		t.Errorf("Gagal buka koneksi")
		return
	}
	dao.RegisterCoreModel()
	logentry := logrus.WithField("testMethod", "Test_bulk_insert")
	logPath := "/tmp/go-test-Test_bulk_insert.json"
	(*logCapture).SetJSONFormatLog(logPath)
	db.LogMode(true)
	var smplData []interface{}
	smplData = append(smplData, &dummyBulkData{Name: "gede", Email: "gede.sutarsa@gmail.com"})
	smplData = append(smplData, &dummyBulkData{Name: "yuli", Email: "yuli.dana@gmail.com"})
	scpBaru := db.NewScope(smplData[0])
	gorm.CreateBulkCallback(scpBaru, smplData)
	if db.Error != nil {
		logentry.WithError(db.Error).Errorf("Gagal")
		t.Fail()
		return
	}
	logentry.WithField("result", smplData).Info("selesai")
}
