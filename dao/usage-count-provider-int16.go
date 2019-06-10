package dao

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"
)

//UsageCountInt16IDQueryResult usage count for data with id = int16
type UsageCountInt16IDQueryResult struct {
	//ModelName model name / table name
	ModelName string `json:"modelName"`
	//DataID id of data to count for usage
	DataID int16 `json:"dataId"`
	//usageCount count of data usage for specified ID
	UsageCount int64 `json:"usageCount"`
}

//GetUsageCountIDInt16 count usage of single data
func GetUsageCountIDInt16(db *gorm.DB, logEntry *logrus.Entry, tableName string, id int16) (usageCount int64, err error) {
	_, mappedUsageCount, errSwap := GetUsageCountsIDInt16(db, logEntry, tableName, id)
	if errSwap != nil {
		err = errSwap
		return
	}
	if rslt, ok := mappedUsageCount[id]; ok {
		usageCount = rslt.UsageCount
	}
	return
}

//GetUsageCountsIDInt16 count usage count for data with id = int16
func GetUsageCountsIDInt16(db *gorm.DB, logEntry *logrus.Entry, tableName string, ids ...int16) (usageCounts []*UsageCountInt16IDQueryResult, mappedUsageCount map[int16]*UsageCountInt16IDQueryResult, err error) {
	idLen := len(ids)
	if idLen == 0 {
		logEntry.WithField("tableName", tableName).Warnf("Parameter for usage count is empty. no query run for this task")
		return
	}
	queries := DefaultDaoManager.GetUsageCountQueries(db.Dialect().GetName(), tableName)
	if len(queries) == 0 {
		logEntry.WithField("tableName", tableName).Warnf("No usage count query was defined for this table. no query was invoked for usage counts")
		return
	}
	usageCounts = make([]*UsageCountInt16IDQueryResult, idLen, idLen)
	mappedUsageCount = make(map[int16]*UsageCountInt16IDQueryResult)
	idString := make([]string, idLen, idLen)
	for index, id := range ids {
		usg := UsageCountInt16IDQueryResult{ModelName: tableName, DataID: id, UsageCount: 0}

		usageCounts[index] = &usg
		mappedUsageCount[id] = &usg
		idString[index] = fmt.Sprintf("%d", id)
	}
	finalIn := strings.Join(idString, ",")
	finalQ := strings.Join(queries, " union all ")
	finalQ = strings.Join(strings.Split(finalQ, ":ids"), fmt.Sprintf("%s", finalIn))
	var rslt *sql.Rows
	if rslt, err = db.Raw(finalQ).Rows(); err != nil {
		logEntry.WithFields(logrus.Fields{"tableName": tableName, "id": ids}).WithError(err).Errorf(`Fail to to invoke usage count, error: %s`, err.Error())
		return
	}
	defer rslt.Close()

	for rslt.Next() {
		var modelName string
		var dataID int16
		var usageCount int64
		rslt.Scan(&modelName, &dataID, &usageCount)
		usgCount := mappedUsageCount[dataID]
		usgCount.UsageCount += usageCount
	}
	return
}
