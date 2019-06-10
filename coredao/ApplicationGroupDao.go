package coredao

import (
	"fmt"

	"github.com/SoftwareUndagi/golibs/coremodel"

	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"
)

//ApplicationGroupDao dao for application group
type ApplicationGroupDao struct {
}

var cacheAppGroupStruct = coremodel.ApplicationGroup{}

//IncrementWithSpecifiedNumber increment usage with specified count
func (p *ApplicationGroupDao) IncrementWithSpecifiedNumber(groupCode int32, increment int16, username string, ipAddress string, db *gorm.DB, logEntry *logrus.Entry) (err error) {
	if increment == 0 {
		return
	}
	var rslt coremodel.ApplicationGroup
	logEntry = logEntry.WithField("groupId", groupCode)
	dbR := db.Where(&coremodel.ApplicationGroup{ID: groupCode}).First(&rslt)
	if dbR.Error != nil {
		logEntry.WithError(dbR.Error).Errorf("Error query group, erorr: %s", dbR.Error.Error())
		return dbR.Error
	}
	if dbR.RecordNotFound() {
		return
	}
	var currCount = rslt.UsageCounter
	rslt.UsageCounter = rslt.UsageCounter + increment
	if rslt.UsageCounter < 0 {
		rslt.UsageCounter = 0
	}
	if currCount != rslt.UsageCounter {
		if rslt2 := db.Save(&rslt); rslt2.Error != nil {
			err = rslt2.Error
			logEntry.Errorf("Fail to update usage count on ApplicationGroup, error: %s", err.Error())
			return
		}
	}
	return
}

//IncrementUsageCount increase usage count on table sec_group
func (p *ApplicationGroupDao) IncrementUsageCount(groupCodes []int32, username string, ipAddress string, db *gorm.DB, logEntry *logrus.Entry) (err error) {
	if len(groupCodes) == 0 {
		return
	}
	updSQL := fmt.Sprintf(`update %s 
	set 
		usage_count =  ifnull(usage_count,0) + 1  , 
		modified_by = ? ,
		modified_by_ip = ? , 
		updatedAt = current_timestamp
	where 
		code in ( ?) `, cacheAppGroupStruct.TableName(db))
	if dbRslt := db.Exec(updSQL, username, ipAddress, groupCodes); dbRslt.Error != nil {
		logEntry.WithError(dbRslt.Error).WithField("groupCodes", groupCodes).Errorf("Fail to increment Sec group , error: %s", dbRslt.Error.Error())
		return dbRslt.Error
	}
	return
}

//DecrementUsageCount decrement usage count on table sec_group
func (p *ApplicationGroupDao) DecrementUsageCount(groupCodes []int32, username string, ipAddress string, db *gorm.DB, logEntry *logrus.Entry) (err error) {
	if len(groupCodes) == 0 {
		return
	}
	updSQL := fmt.Sprintf(`update %s 
	set 
		usage_count =  ifnull(usage_count,0) - 1  , 
		modified_by = ? ,
		modified_by_ip = ? , 
		updatedAt = current_timestamp
	where 
		code in ( ?) 
		and ifnull(usage_count,0) > 0 `, cacheAppGroupStruct.TableName(db))
	if dbRslt := db.Exec(updSQL, username, ipAddress, groupCodes); dbRslt.Error != nil {
		logEntry.WithError(dbRslt.Error).WithField("groupCodes", groupCodes).Errorf("Fail to decrement sec group usage , error: %s", dbRslt.Error.Error())
		return dbRslt.Error
	}
	return
}

//FindByUUID query data by uuid
func (p *ApplicationGroupDao) FindByUUID(UUID string, db *gorm.DB, logEntry *logrus.Entry) (data coremodel.ApplicationGroup, found bool, err error) {
	if rslt := db.Where("uid = ? ", UUID).First(&data); rslt.Error != nil {
		//return data, false, rslt.Error
		err = rslt.Error
	} else {
		found = rslt.RowsAffected > 0
	}
	return
}

//FindByUUIDs find app groups by uuids
func (p *ApplicationGroupDao) FindByUUIDs(UUIDs []string, db *gorm.DB, logEntry *logrus.Entry) (data []coremodel.ApplicationGroup, found bool, err error) {
	if rslt := db.Where("uid in (?)", UUIDs).Find(&data); rslt.Error != nil {
		//return data, false, rslt.Error
		err = rslt.Error
	} else {
		found = rslt.RowsAffected > 0
	}
	return
}

//FindByID find group by id
func (p *ApplicationGroupDao) FindByID(id int32, db *gorm.DB, logEntry *logrus.Entry) (data coremodel.ApplicationGroup, found bool, err error) {
	if rslt := db.Where(&coremodel.ApplicationGroup{ID: id}); rslt.Error != nil {
		//return data, false, rslt.Error
		err = rslt.Error
	} else {
		found = rslt.RowsAffected > 0
	}
	return
}

//FindByIDs find app group by ids
func (p *ApplicationGroupDao) FindByIDs(ids int32, db *gorm.DB, logEntry *logrus.Entry) (data coremodel.ApplicationGroup, found bool, err error) {
	if rslt := db.Where("id in (?)", ids).Find(&data); rslt.Error != nil {
		//return data, false, rslt.Error
		err = rslt.Error
	} else {
		found = rslt.RowsAffected > 0
	}
	return
}

//PlugID assign ID ( from db ) to struct.
//ID is not serialized ( on json field) from and to client. so when user do update to data, server will not recieve ID of ApplicationGroup, only UUID
// with this uuid string, id will be queried from server and plug back to data
func (p *ApplicationGroupDao) PlugID(group coremodel.ApplicationGroup, username string, ipAddress string, db *gorm.DB, logEntry *logrus.Entry) (groupFromDB coremodel.ApplicationGroup, err error) {
	var dbRslt *gorm.DB
	dbRslt = db.Where(&coremodel.ApplicationGroup{UUID: group.UUID}).First(&groupFromDB)
	if dbRslt.Error != nil {
		err = dbRslt.Error
		logEntry.WithField("group", group).WithError(err).Errorf("Fail to Group. error: %s", err.Error())
		return
	}
	if !dbRslt.RecordNotFound() {
		group.ID = groupFromDB.ID
	}
	return
}

//PlugIDs assign ids to Application groups. this bulk version of PlugID
func (p *ApplicationGroupDao) PlugIDs(groups []coremodel.ApplicationGroup, username string, ipAddress string, db *gorm.DB, logEntry *logrus.Entry) (groupFromDB []coremodel.ApplicationGroup, err error) {
	var assignUUIDs []string
	idxByUUIDs := make(map[string]coremodel.ApplicationGroup)
	for _, m := range groups {
		assignUUIDs = append(assignUUIDs, m.UUID)
		idxByUUIDs[m.UUID] = m
	}

	if dbRslt := db.Where("uid in ( ?) ", assignUUIDs).Find(&groupFromDB); dbRslt.Error != nil {
		err = dbRslt.Error
		logEntry.WithField("groups", groups).WithField("groupUuids", assignUUIDs).WithError(err).Errorf("Fail to query for menu assignment by uuid, error: %s", err.Error())
		return
	}
	for _, m := range groupFromDB {
		if mnu, ok := idxByUUIDs[m.UUID]; ok {
			mnu.ID = m.ID
		}
	}
	return
}

//sharedApplicationGroupDao dao applicationGroup. untuk di pakai internal
var sharedApplicationGroupDao = ApplicationGroupDao{}
