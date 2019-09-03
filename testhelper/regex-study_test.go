package testhelper

import (
	"fmt"
	"regexp"
	"strings"
	"testing"

	"github.com/sirupsen/logrus"
)

type mapALias map[string]interface{}

func TestFindDoubdleDotOnString(t *testing.T) {
	logCapture := CaptureLog(t)
	logPath := "/tmp/go-TestFindDoubdleDotOnString.json"
	logCapture.Release()
	logCapture.SetJSONFormatLog(logPath)
	//\040 = spasi
	//072 = :
	var re = regexp.MustCompile(`(?m)(\040\072[a-zA-Z0-9]{1,})`)
	var str = `select 
	1 ID , 'gede.sutarsa' CreatorName , '127.0.0.1' CreatorIPAddress , null CreatedAt , null UpdatedAt, null ModifiedBy, null ModifiedIPAddress , 	m.user_name DetailCode , 
	   'TEST::MasterUser' LovID , 
	   m.real_name  Label , 
	   m.email   Value1 ,  
	   m.locale_code  Value2  , 
	   null  I18nKey  , 
	   @rownum:=ifnull(@rownum , 0 )+1 SequenceNo
   from 
	   sec_user m 
   where
	   ( m.user_name = :user1 or
		m.user_name = :Dodol1 or
		m.user_name = :garut)   
		and age < :garut

   order by  m.user_name `
	matches := re.FindAllString(str, -1)

	ln := len(matches)
	if ln > 0 {
		replaceableString := make([]string, 0)
		replaceableKey := make(map[string]bool)
		var rslt []interface{}
		paramVal := mapALias{
			"user1": "gede.sutarsa",
			"garut": 37}
		for _, mtch := range matches {
			xx := strings.Trim(mtch, " ")
			plain := xx[1:]
			if val, ok := paramVal[plain]; ok {
				rslt = append(rslt, val)
				replaceableKey[plain] = true
				replaceableString = append(replaceableString, mtch)
			}
		}
		var replcRegex string
		//keys := reflect.ValueOf(replaceableKey).MapKeys()

		for key := range replaceableKey {

			if len(replcRegex) > 0 {
				replcRegex = replcRegex + "|"
			}
			//logrus.WithFields(logrus.Fields{"KEY": key}).Info("ss")
			replcRegex = replcRegex + fmt.Sprintf("(\\040\\072%v)", key)
		}
		logrus.WithFields(logrus.Fields{"replcRegex": replcRegex, "qMarks": rslt}).Info("selesai")
		rplRegx := regexp.MustCompile(replcRegex)
		rsltQWithReplacement := rplRegx.ReplaceAllString(str, " ? ")
		logrus.Info(rsltQWithReplacement)

	} else {
		println("0 match key")
	}

}
