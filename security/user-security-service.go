package security

import (
	"github.com/SoftwareUndagi/golibs/coredao"
	"github.com/SoftwareUndagi/golibs/coremodel"
	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"
)

var appRoleDao = coredao.ApplicationUserRoleDao{}

//ReadUserAuthenticationData read user authentication data
func ReadUserAuthenticationData(userID int32, db *gorm.DB, logEntry *logrus.Entry) (userAuthData SimpleUserData, found bool, err error) {
	var dbUser coremodel.ApplicationUser
	if dbRslt := db.Where(coremodel.ApplicationUser{ID: userID}).First(&dbUser); dbRslt != nil {
		err = dbRslt.Error
		logEntry.WithField("userId", userID).WithError(err).Errorf("Unable to query for user, error: %s", err.Error())
		return
	} else if dbRslt.RecordNotFound() {
		return
	}
	found = true
	userAuthData = GenerateSimpleUserData(dbUser)
	var roles []string
	// var dao
	if roles, err = appRoleDao.QueryUserAuthoritiesByUserID(userID, db, logEntry); err != nil {
		logEntry.WithField("userID", userID).WithError(err).Errorf("Unable to query for user roles: %s ", err.Error())
		return
	}
	userAuthData.UserRoles = roles
	return
}

//ReadUserAuthenticationDataByUserData read role and generate security data
func ReadUserAuthenticationDataByUserData(userData coremodel.ApplicationUser, db *gorm.DB, logEntry *logrus.Entry) (userAuthData SimpleUserData, err error) {
	userAuthData = GenerateSimpleUserData(userData)
	var roles []string
	// var dao
	userID := userData.ID
	if roles, err = appRoleDao.QueryUserAuthoritiesByUserID(userID, db, logEntry); err != nil {
		logEntry.WithField("userID", userID).WithError(err).Errorf("Unable to query for user roles: %s ", err.Error())
		return
	}
	userAuthData.UserRoles = roles
	return
}

//GenerateSimpleUserData generate SimpleUserData from db data
func GenerateSimpleUserData(userData coremodel.ApplicationUser) (userAuthData SimpleUserData) {
	userAuthData.ID = userData.ID
	userAuthData.Email = userData.Email
	userAuthData.Phone1 = userData.Phone1
	userAuthData.Phone2 = userData.Phone2
	userAuthData.RealName = userData.RealName
	userAuthData.SovereignAuthID = userData.SovereignAuthID
	userAuthData.UUID = userData.UUID
	userAuthData.Username = userData.Username
	return
}
