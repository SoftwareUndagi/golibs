package dao

import (
	"bufio"
	"bytes"
	"io"
	"os"

	"github.com/sirupsen/logrus"

	"github.com/jinzhu/gorm"
)

//RunSQLFile run sql file. scan for ; on end of line
func RunSQLFile(db *gorm.DB, loggerEntry *logrus.Entry, SQLFilePath string, ignoreOnError bool) (err error) {

	file, errOpenFile := os.Open(SQLFilePath)
	if errOpenFile != nil {
		loggerEntry.WithError(errOpenFile).Errorf("Failed to open JSON file : %s ", errOpenFile.Error())
		err = errOpenFile
		return
	}
	defer file.Close()
	reader := bufio.NewReader(file)
	var eofFound bool
	var currentScript []byte
	qSeparator := []byte{';'}
	addToCurrentSQLScript := func(readedData []byte) {
		if len(currentScript) > 0 {
			currentScript = append(currentScript, ' ')
		}
		currentScript = append(currentScript, readedData...)
	}
	SQLinvoker := func(script []byte) (err error) {
		if len(script) == 0 {
			return
		}
		sqlStr := string(script)
		if rslt := db.Exec(sqlStr); rslt.Error != nil {
			loggerEntry.WithError(rslt.Error).WithField("sql", sqlStr).Errorf("Sql invoke failed, error : %s", rslt.Error.Error())
			return rslt.Error
		}
		return nil
	}

	lineNo := 0
	for {
		rowData, errRead := reader.ReadBytes('\n')
		if errRead != nil {
			if errRead == io.EOF {
				eofFound = true
			} else {
				loggerEntry.WithError(err).Errorf("Error membaca buffer, row filled index [%d]. error: %s", lineNo, err.Error())
				return
			}
		}
		addToCurrentSQLScript(rowData)
		if bytes.Contains(rowData, qSeparator) {
			err = SQLinvoker(currentScript)
			if err != nil {
				if !ignoreOnError {
					return
				}
				loggerEntry.Warnf("Script was error, but ignored")
			}
			currentScript = currentScript[:0]
		}

		lineNo++
		if eofFound {
			if len(currentScript) > 0 {
				err = SQLinvoker(currentScript)
				if err != nil {
					return
				}
			}
			break
		}

	}

	return
}
