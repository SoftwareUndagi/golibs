package dao

import (
	"testing"

	"github.com/sirupsen/logrus"

	"github.com/SoftwareUndagi/golibs/testhelper"
)

func Test_ConvertNamedParameterToPlainParameter(t *testing.T) {
	logCapture := testhelper.CaptureLog(t)
	logPath := "/tmp/go-Test_ConvertNamedParameterToPlainParameter.json"
	logCapture.Release()
	logCapture.SetJSONFormatLog(logPath)
	sqlSample := ` where 1=1
			and type = :employeeType
			and (address like :pattern or office_address like :pattern)
			and age > :ageParam`
	sampleData := map[string]interface{}{"employeeType": "C1", "pattern": "%dodol", "ageParam": 30}
	cnvSql, qParam := ConvertNamedParameterToPlainParameter(sqlSample, sampleData)
	if len(qParam) != 4 || (qParam[0] != "C1" || qParam[1] != "%dodol" || qParam[2] != "%dodol" || qParam[3] != 30) {
		logrus.WithFields(logrus.Fields{"sql": cnvSql, "qParam": qParam}).Errorf("Hasil tidak sesuai, len seharusnya 4(actual %d)", len(qParam))
		t.Fail()
	}
	logrus.WithFields(logrus.Fields{"sql": cnvSql, "qParam": qParam}).Info("selesai")
}
