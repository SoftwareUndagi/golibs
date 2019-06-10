package dao

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"
)

//UsageCountStringIDQueryResult struct container usage count(id = string)
type UsageCountStringIDQueryResult struct {
	//ModelName model name / table name
	ModelName string `json:"modelName"`
	//DataID id of data to count for usage
	DataID string `json:"dataId"`
	//usageCount count of data usage for specified ID
	UsageCount int64 `json:"usageCount"`
}

//GetUsageCountIDString  count usage count. for single item
func GetUsageCountIDString(db *gorm.DB, logEntry *logrus.Entry, tableName string, id string) (usageCount int64, err error) {
	_, mappedUsageCount, errSwap := GetUsageCountsIDString(db, logEntry, tableName, id)
	if errSwap != nil {
		err = errSwap
		return
	}
	if rslt, ok := mappedUsageCount[id]; ok {
		usageCount = rslt.UsageCount
	}
	return

}

//GetUsageCountsIDString get usage for data with id with type string . data ids is multiple item
func GetUsageCountsIDString(db *gorm.DB, logEntry *logrus.Entry, tableName string, ids ...string) (usageCounts []UsageCountStringIDQueryResult, mappedUsageCount map[string]UsageCountStringIDQueryResult, err error) {
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
	usageCounts = make([]UsageCountStringIDQueryResult, idLen, idLen)
	// indexedUsage := make(map[string]UsageCountStringIDQueryResult)
	mappedUsageCount = make(map[string]UsageCountStringIDQueryResult)
	sanitizedIds := make([]string, idLen, idLen)
	for index, id := range ids {
		usg := UsageCountStringIDQueryResult{ModelName: tableName, DataID: id, UsageCount: 0}
		usageCounts[index] = usg
		mappedUsageCount[id] = usg
		if strings.Contains(id, "'") {
			sanitizedIds[index] = "'" + strings.Join(strings.Split(id, "'"), "''") + "'"
		} else {
			sanitizedIds[index] = "'" + id + "'"
		}
	}
	finalIn := strings.Join(sanitizedIds, ",")
	finalQ := strings.Join(queries, " union all ")
	finalQ = strings.Join(strings.Split(finalQ, ":ids"), fmt.Sprintf("%s", finalIn))
	var rslt *sql.Rows
	if rslt, err = db.Raw(finalQ).Rows(); err != nil {
		logEntry.WithFields(logrus.Fields{"tableName": tableName, "id": ids}).WithError(err).Errorf(`Fail to to invoke usage count, error: %s`, err.Error())
		return
	}
	defer rslt.Close()
	var modelName, dataID string
	var currentUsageCount int64
	for rslt.Next() {
		rslt.Scan(&modelName, &dataID, &currentUsageCount)
		usgCount := mappedUsageCount[dataID]
		usgCount.UsageCount += currentUsageCount
	}
	return
}
