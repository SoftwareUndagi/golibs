package dao

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"
)

//UsageCountInt32IDQueryResult usage count for data with id = int32
type UsageCountInt32IDQueryResult struct {
	//ModelName model name / table name
	ModelName string `json:"modelName"`
	//DataID id of data to count for usage
	DataID int32 `json:"dataId"`
	//usageCount count of data usage for specified ID
	UsageCount int64 `json:"usageCount"`
}

//GetUsageCountIDInt32 count usage of single data
func GetUsageCountIDInt32(db *gorm.DB, logEntry *logrus.Entry, tableName string, id int32) (usageCount int64, err error) {
	_, mappedUsageCount, errSwap := GetUsageCountsIDInt32(db, logEntry, tableName, id)
	if errSwap != nil {
		err = errSwap
		return
	}
	if rslt, ok := mappedUsageCount[id]; ok {
		usageCount = rslt.UsageCount
	}
	return
}

//GetUsageCountsIDInt32 count usage count for data with id = int32
func GetUsageCountsIDInt32(db *gorm.DB, logEntry *logrus.Entry, tableName string, ids ...int32) (usageCounts []*UsageCountInt32IDQueryResult, mappedUsageCount map[int32]*UsageCountInt32IDQueryResult, err error) {
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
	usageCounts = make([]*UsageCountInt32IDQueryResult, idLen, idLen)
	mappedUsageCount = make(map[int32]*UsageCountInt32IDQueryResult)
	idString := make([]string, idLen, idLen)
	for index, id := range ids {
		usg := UsageCountInt32IDQueryResult{ModelName: tableName, DataID: id, UsageCount: 0}
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
	var modelName string
	var dataID int32
	var currentUsageCount int64
	for rslt.Next() {
		rslt.Scan(&modelName, &dataID, &currentUsageCount)
		if usgCount, ok := mappedUsageCount[dataID]; ok {
			usgCount.UsageCount += currentUsageCount
			continue
		}
		logEntry.WithField("dataId", dataID).Warnf("Data id %d tidak di temukan, silakan di cek kembali dengan query anda", dataID)

	}
	return
}
