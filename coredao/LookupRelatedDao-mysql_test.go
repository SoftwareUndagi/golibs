package coredao

import (
	"runtime"
	"testing"

	"github.com/SoftwareUndagi/golibs/coremodel"

	"github.com/SoftwareUndagi/golibs/common"
	"github.com/SoftwareUndagi/golibs/dao"
	"github.com/SoftwareUndagi/golibs/testhelper"
	"github.com/sirupsen/logrus"
)

func TestLookup(t *testing.T) {

	db, logCapture, err := testhelper.MysqlTestAppSetup(t)
	dao.RegisterCoreModel()
	logentry := logrus.WithField("testMethod", "TestLookup")
	logPath := "/tmp/go-test-TestLookup.json" //common.CreateFilePathOnTempDirectory("go-test-TestLookup.json")
	(*logCapture).SetJSONFormatLog(logPath)
	if err != nil {
		t.Errorf("Error open connection , error: %s", err.Error())
		t.Fail()
		return
	}
	_, filename, _, _ := runtime.Caller(0)
	directory := common.GetFileDirectory(filename)
	sqlFile := common.AppendFilePath(directory, "LookupRelatedDao-mysql_test.sql")
	db.LogMode(true)
	err = dao.RunSQLFile(db, logentry, sqlFile, true)
	if err != nil {
		t.Error(err.Error())
		logentry.WithError(err).WithFields(logrus.Fields{"path": sqlFile}).Errorf("Fail to invoke , error: %s", err.Error())
		t.Fail()
		return
	}
	lovHeaders, errFind := FindLookupHeaders([]string{"SIMPLE_LOOKUP", "TEST::MasterUser"}, db, logentry)
	if errFind != nil {
		t.Error(errFind.Error())
		logentry.WithError(errFind).Errorf("Fail to find lookup header, error: %s", errFind.Error())
		t.Fail()
		return
	}
	lh := len(lovHeaders)
	if lh != 2 {

		t.Errorf("not match, suppose tobe 2, actual: %d", lh)
		logentry.WithField("method", "FindLookupHeaders").Errorf("Fail query for header. result suppose to be 2 , actual found: %d", len(lovHeaders))
		t.Fail()
		return
	}
	logentry.WithFields(logrus.Fields{"lookupHeaders": lovHeaders}).Infof("Read all lokkup complete")
	var lookupBySQL coremodel.LookupHeader
	for _, lkh := range lovHeaders {
		if lkh.ID == "TEST::MasterUser" {
			lookupBySQL = lkh
			break
		}
	}
	errVersi := QueryForLookupWithSQLVersion([]*coremodel.LookupHeader{&lookupBySQL}, db, logentry)
	if errVersi != nil {
		t.Error(errVersi.Error())
		logentry.WithError(errVersi).Errorf("Fail to read varsi, error: %s", errVersi.Error())
		t.Fail()
		return
	}
	logentry.WithField("raw", lookupBySQL).Infof("Read version  success. versi: %s", (lookupBySQL).Version)
	flatDetails, indexedDetails, errDtlSQL := FindSQLDrivenLookupDetails([]coremodel.LookupHeader{lookupBySQL}, db, logentry)
	if errDtlSQL != nil {
		logentry.WithError(errDtlSQL).Errorf("Error while read lookup detail by sql, error: %s", errDtlSQL.Error())
		t.Errorf(errDtlSQL.Error())
		t.Fail()
		return
	}
	logentry.WithFields(logrus.Fields{"flatDetail": flatDetails, "indexedDetails": indexedDetails}).Infof("Read LOV detail by sql complete")
}

func TestGormQuery(t *testing.T) {
	db, logCapture, err := testhelper.MysqlTestAppSetup(t)
	dao.RegisterCoreModel()
	logentry := logrus.WithField("testMethod", "TestGormQuery")
	logPath := "/tmp/go-test-TestGormQuery.json"
	(*logCapture).SetJSONFormatLog(logPath)
	if err != nil {
		t.Errorf("Error open connection , error: %s", err.Error())
		t.Fail()
		return
	}
	_, filename, _, _ := runtime.Caller(0)
	directory := common.GetFileDirectory(filename)
	sqlFile := common.AppendFilePath(directory, "LookupRelatedDao-mysql_test.sql")
	db.LogMode(true)
	err = dao.RunSQLFile(db, logentry, sqlFile, true)
	if err != nil {
		t.Error(err.Error())
		logentry.WithError(err).WithFields(logrus.Fields{"path": sqlFile}).Errorf("Fail to invoke , error: %s", err.Error())
		t.Fail()
		return
	}

	var lookupUppers []coremodel.LookupHeader
	s := db.Where("id in (?)  ", []string{"TEST::MasterUser", "SIMPLE_LOOKUP"})
	rlst := s.Find(&lookupUppers)
	if rlst.Error != nil {
		logentry.WithError(rlst.Error).Errorf("Failt query, error : %s ", rlst.Error.Error())
		t.Fail()
	}

	logentry.WithField("result", lookupUppers).Info("Selesai")
	var lkpDetails []coremodel.LookupDetail
	db.Find(&lkpDetails)
	logentry.WithField("lkpDetails", lkpDetails).Info("Selesai")

}

func Test_Run_query_with_named_parameter(t *testing.T) {
	db, logCapture, err := testhelper.MysqlTestAppSetup(t)
	dao.RegisterCoreModel()
	logentry := logrus.WithField("testMethod", "Test_Run_query_with_named_parameter")
	logPath := "/tmp/Test_Run_query_with_named_parameter.json"
	(*logCapture).SetJSONFormatLog(logPath)
	if err != nil {
		t.Errorf("Error open connection , error: %s", err.Error())
		t.Fail()
		return
	}
	var rslt []coremodel.LookupHeader
	swapDb := db.NamedParamWere("is_cacheable = :cacheableFlag and ( remark like :where1 or remark like :where2)", map[string]interface{}{"cacheableFlag": "Y", "where1": "%user)", "where2": "%content%"})
	if dbRslt := swapDb.Find(&rslt); dbRslt.Error != nil {
		t.Fail()
		logentry.WithError(dbRslt.Error).Errorf("Fail to invoke named query, error: %s", dbRslt.Error.Error())
	}
	logentry.WithField("result", rslt).Infof("Selesai")

}

func TestRunQueryForStruct(t *testing.T) {
	db, logCapture, err := testhelper.MysqlTestAppSetup(t)
	dao.RegisterCoreModel()
	logentry := logrus.WithField("testMethod", "TestGormQuery")
	logPath := "/tmp/TestRunQueryForStruct.json"
	(*logCapture).SetJSONFormatLog(logPath)
	if err != nil {
		t.Errorf("Error open connection , error: %s", err.Error())
		t.Fail()
		return
	}

	/*
		_, filename, _, _ := runtime.Caller(0)
		directory := common.GetFileDirectory(filename)
		sqlFile := common.AppendFilePath(directory, "LookupRelatedDao-mysql_test.sql")
		db.LogMode(true)
		err = dao.RunSQLFile(db, logentry, sqlFile, true)
		if err != nil {
			t.Error(err.Error())
			logentry.WithError(err).WithFields(logrus.Fields{"path": sqlFile}).Errorf("Fail to invoke , error: %s", err.Error())
			t.Fail()
			return
		}
	*/
	SQLForStruct := `
	set @rownum=0;
	select 
		1 id , 
		m.user_name detail_code , 
		'TEST::MasterUser' lov_id , 
		m.real_name  label , 
		m.email   val_1 ,  
		m.locale_code  val_2  , 
		null  i18n_key  , 
		@rownum:=ifnull(@rownum , 0 )+1 seq_no
	from 
		sec_user m 
	where
		m.user_name in ('dev-test-helper.gede.sutarsa@gmail.com' , 'dev-test-helper.yuli.dana@gmail.com') 
	order by  m.user_name`
	var swapRslt []coremodel.LookupDetail
	rsltQuery := db.Raw(SQLForStruct).Find(&swapRslt)
	if rsltQuery.Error != nil {
		t.Fail()
		logentry.WithError(rsltQuery.Error).Errorf("Gagal query , error : %s ", rsltQuery.Error.Error())
		return
	}
	logentry.WithFields(logrus.Fields{"affected": rsltQuery.RowsAffected, "result": &swapRslt}).Info("Query done")

	t.Log("Selesai")

}
