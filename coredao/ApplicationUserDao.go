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
		usr.ModifiedIPAddress = ipAddress
		usr.UpdatedAt = &skr
		if rslt := db.Save(&usr); rslt.Error != nil {
			logEntry.WithError(rslt.Error).WithField("userID", userID).WithField("increment", increment).Errorf("Fail to increment count, error: %s", err.Error())
			return rslt.Error
		}
	}
	return
}
