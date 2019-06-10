package coredao

import (
	"fmt"

	"github.com/SoftwareUndagi/golibs/coremodel"
	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"
)

//ApplicationRoleDao consolidation function for auth related data
type ApplicationRoleDao struct {
}

var cacheTableAuth = coremodel.ApplicationRole{}
var commonApplicationRoleDao = ApplicationRoleDao{}

//IncrementUsageCount increment usage_count field on table sec_authority
func (p *ApplicationRoleDao) IncrementUsageCount(authCodes []string, username string, ipAddress string, db *gorm.DB, logEntry *logrus.Entry) (err error) {
	if len(authCodes) == 0 {
		return
	}
	updSQL := fmt.Sprintf(`update %s 
	set 
		usage_count =  ifnull(usage_count,0) + 1  , 
		modified_by = ? ,
		modified_by_ip = ? , 
		updatedAt = current_timestamp
	where 
		code in ( ?) `, cacheTableAuth.TableName(db))
	if dbRslt := db.Exec(updSQL, username, ipAddress, authCodes); dbRslt.Error != nil {
		logEntry.WithError(dbRslt.Error).WithField("authCodes", authCodes).Errorf("Fail to increment authorities , error: %s", dbRslt.Error.Error())
		return dbRslt.Error
	}
	return
}

//DecrementUsageCount decrement usage count on user sec_authority
func (p *ApplicationRoleDao) DecrementUsageCount(authCodes []string, username string, ipAddress string, db *gorm.DB, logEntry *logrus.Entry) (err error) {
	if len(authCodes) == 0 {
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
		and ifnull(usage_count,0) > 0 `, cacheTableAuth.TableName(db))
	if dbRslt := db.Exec(updSQL, username, ipAddress, authCodes); dbRslt.Error != nil {
		logEntry.WithError(dbRslt.Error).WithField("authCodes", authCodes).Errorf("Fail to decrement authorities , error: %s", dbRslt.Error.Error())
		return dbRslt.Error
	}
	return
}
