package coredao

import (
	"github.com/SoftwareUndagi/golibs/coremodel"
	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"
)

//ApplicationUserGroupDao dao for ApplicationUserGroupDao
type ApplicationUserGroupDao struct {
}

var cachedApplicationUserGroup = coremodel.ApplicationUserGroup{}
var sharedApplicationUserGroupDao = ApplicationUserGroupDao{}

//PlugID assign user id to user group struct
//
func (p *ApplicationUserGroupDao) PlugID(userGroup coremodel.ApplicationUserGroup, username string, ipAddress string, db *gorm.DB, logEntry *logrus.Entry) (userGroupFromDB coremodel.ApplicationUserGroup, err error) {
	var dbRslt *gorm.DB
	dbRslt = db.Where(&coremodel.ApplicationUserGroup{UUID: userGroup.UUID}).First(&userGroupFromDB)
	if dbRslt.Error != nil {
		err = dbRslt.Error
		logEntry.WithField("userGroup", userGroup).WithError(err).Errorf("Fail to query for user group. error: %s", err.Error())
		return
	}
	if !dbRslt.RecordNotFound() {
		userGroup.ID = userGroupFromDB.ID
		userGroup.GroupID = userGroupFromDB.GroupID
		userGroup.UserID = userGroupFromDB.UserID
	}
	return
}

//PlugIDs assign ids to Application user groups. this bulk version of PlugID
func (p *ApplicationUserGroupDao) PlugIDs(userGroups []coremodel.ApplicationUserGroup, username string, ipAddress string, db *gorm.DB, logEntry *logrus.Entry) (userGroupsFromDB []coremodel.ApplicationUserGroup, err error) {
	var assignUUIDs []string
	idxByUUIDs := make(map[string]coremodel.ApplicationUserGroup)
	for _, m := range userGroups {
		assignUUIDs = append(assignUUIDs, m.UUID)
		idxByUUIDs[m.UUID] = m
	}

	if dbRslt := db.Where("uid in ( ?) ", assignUUIDs).Find(&userGroupsFromDB); dbRslt.Error != nil {
		err = dbRslt.Error
		logEntry.WithField("userGroups", userGroups).WithField("userGroupUuids", assignUUIDs).WithError(err).Errorf("Fail to query for user group  by uuid, error: %s", err.Error())
		return
	}
	for _, m := range userGroupsFromDB {
		if mnu, ok := idxByUUIDs[m.UUID]; ok {
			mnu.ID = m.ID
		}
	}
	return
}
