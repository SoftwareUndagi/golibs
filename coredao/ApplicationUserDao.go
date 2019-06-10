package coredao

import (
	"time"

	"github.com/SoftwareUndagi/golibs/coremodel"
	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"
)

//ApplicationUserDao dao for app user
type ApplicationUserDao struct {
}

var appUserCache = coremodel.ApplicationUser{}

//IncrementUsageCount increment usage count on user
func (p *ApplicationUserDao) IncrementUsageCount(userID int32, increment int16, username string, ipAddress string, db *gorm.DB, logEntry *logrus.Entry) (err error) {
	if increment == 0 {
		return
	}
	var usr coremodel.ApplicationUser
	db.Where(&coremodel.ApplicationUser{ID: userID}).First(&usr)
	prevUsageCount := usr.UsageCounter
	usr.UsageCounter += increment
	if usr.UsageCounter < 0 {
		usr.UsageCounter = 0
	}
	skr := time.Now()
	if prevUsageCount != usr.UsageCounter {
		usr.ModifiedBy = &username
		usr.ModifiedIPAddress = &ipAddress
		usr.UpdatedAt = &skr
		if rslt := db.Save(&usr); rslt.Error != nil {
			logEntry.WithError(rslt.Error).WithField("userID", userID).WithField("increment", increment).Errorf("Fail to increment count, error: %s", err.Error())
			return rslt.Error
		}
	}
	return
}

//PlugID assign user id to user struct
//
func (p *ApplicationUserDao) PlugID(user coremodel.ApplicationUser, username string, ipAddress string, db *gorm.DB, logEntry *logrus.Entry) (userFromDB coremodel.ApplicationUser, err error) {
	var dbRslt *gorm.DB
	dbRslt = db.Where(&coremodel.ApplicationUser{UUID: user.UUID}).First(&userFromDB)
	if dbRslt.Error != nil {
		err = dbRslt.Error
		logEntry.WithField("user", user).WithError(err).Errorf("Fail to query for user. error: %s", err.Error())
		return
	}
	if !dbRslt.RecordNotFound() {
		user.ID = userFromDB.ID
	}
	return
}

//PlugIDs assign ids to Application users. this bulk version of PlugID
func (p *ApplicationUserDao) PlugIDs(users []coremodel.ApplicationUser, username string, ipAddress string, db *gorm.DB, logEntry *logrus.Entry) (usersFromDB []coremodel.ApplicationUser, err error) {
	var assignUUIDs []string
	idxByUUIDs := make(map[string]coremodel.ApplicationUser)
	for _, m := range users {
		assignUUIDs = append(assignUUIDs, m.UUID)
		idxByUUIDs[m.UUID] = m
	}

	if dbRslt := db.Where("uid in ( ?) ", assignUUIDs).Find(&usersFromDB); dbRslt.Error != nil {
		err = dbRslt.Error
		logEntry.WithField("users", users).WithField("userUuids", assignUUIDs).WithError(err).Errorf("Fail to query for user  by uuid, error: %s", err.Error())
		return
	}
	for _, m := range usersFromDB {
		if mnu, ok := idxByUUIDs[m.UUID]; ok {
			mnu.ID = m.ID
		}
	}
	return
}
