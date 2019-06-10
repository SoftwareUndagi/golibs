package coredao

import (
	"fmt"

	"github.com/SoftwareUndagi/golibs/coremodel"

	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"
)

//ApplicationUserRoleDao dao for user auth
type ApplicationUserRoleDao struct {
}

var sampleAppAuthStruct = coremodel.ApplicationUserRole{}

//QueryUserAuthoritiesByUserID query for user auth by user id
//userID id = user id to limit user
func (p *ApplicationUserRoleDao) QueryUserAuthoritiesByUserID(userID int32, db *gorm.DB, logEntry *logrus.Entry) (authCodes []string, err error) {
	logEntry = logEntry.WithField("userId", userID)
	selectSmt := fmt.Sprintf(`select authority_code authCode  from %s where user_id = ? `, sampleAppAuthStruct.TableName(db))
	dbRslt := db.Exec(selectSmt, userID)
	if dbRslt.Error != nil {
		logEntry.WithError(dbRslt.Error).Errorf("Fail ")
		return nil, err
	}
	rows, errRows := dbRslt.Rows()
	if errRows != nil {
		return nil, errRows
	}
	defer rows.Close()
	for rows.Next() {
		var authCode string
		rows.Scan(&authCode)
		authCodes = append(authCodes, authCode)
	}
	return
}

//QueryUserAuthoritiesByUserUUID query for auth code by user uuid
func (p *ApplicationUserRoleDao) QueryUserAuthoritiesByUserUUID(userUUID string, db *gorm.DB, logEntry *logrus.Entry) (authCodes []string, err error) {
	selectSmt := fmt.Sprintf(`select auth.authority_code authCode  
	from 
		%s auth inner join 
		%s usr on auth.user_id = usr.id 
	where 
		usr.uid =   ? `, sampleAppAuthStruct.TableName(db), appUserCache.TableName(db))
	dbRslt := db.Exec(selectSmt, userUUID)
	if dbRslt.Error != nil {
		return nil, err
	}
	rows, errRows := dbRslt.Rows()
	if errRows != nil {
		return nil, errRows
	}
	defer rows.Close()
	for rows.Next() {
		var authCode string
		rows.Scan(&authCode)
		authCodes = append(authCodes, authCode)
	}
	return
}

//SaveUserAuthorities save user authorities. this process include increment usage count on user and authorities
//userID = id of user to update
func (p *ApplicationUserRoleDao) SaveUserAuthorities(userID int32, username string, ipAddress string, authCodes []string, db *gorm.DB, logEntry *logrus.Entry) (err error) {
	authCurrents, errFindCoded := p.QueryUserAuthoritiesByUserID(userID, db, logEntry)
	if errFindCoded != nil {
		err = errFindCoded
		return
	}
	var removedCodes []string
	var addedCodes []string
	idxCurrentDb := make(map[string]bool)
	idxNewRole := make(map[string]bool)
	for _, addAuth := range authCodes {
		idxNewRole[addAuth] = true
	}
	for _, authCur := range authCurrents {
		idxCurrentDb[authCur] = true
		if _, ok := idxNewRole[authCur]; !ok {
			removedCodes = append(removedCodes, authCur)
		}
	}
	for _, addAuth := range authCodes {
		idxNewRole[addAuth] = true
		if _, ok2 := idxCurrentDb[addAuth]; !ok2 {
			addedCodes = append(addedCodes, addAuth)
		}
	}
	usageDiff := len(addedCodes) - len(removedCodes)
	if usageDiff > 0 {
		x := ApplicationUserDao{}
		errInc := x.IncrementUsageCount(userID, int16(usageDiff), username, ipAddress, db, logEntry)
		if errInc != nil {
			err = errInc
			return
		}
	}
	if len(removedCodes) > 0 {
		commonApplicationRoleDao.DecrementUsageCount(removedCodes, username, ipAddress, db, logEntry)
	}
	if len(addedCodes) > 0 {
		commonApplicationRoleDao.IncrementUsageCount(addedCodes, username, ipAddress, db, logEntry)
	}
	// cari dulu auth skr
	return
}

//PlugID assign id to data application user role
func (p *ApplicationUserRoleDao) PlugID(userRole coremodel.ApplicationUserRole, username string, ipAddress string, db *gorm.DB, logEntry *logrus.Entry) (userRoleFromDB coremodel.ApplicationUserRole, err error) {
	var dbRslt *gorm.DB
	dbRslt = db.Where(&coremodel.ApplicationUserRole{UUID: userRole.UUID}).First(&userRoleFromDB)
	if dbRslt.Error != nil {
		err = dbRslt.Error
		logEntry.WithField("userRole", userRole).WithError(err).Errorf("Fail to query for user role. error: %s", err.Error())
		return
	}
	if !dbRslt.RecordNotFound() {
		userRole.ID = userRoleFromDB.ID
		userRole.UserID = userRoleFromDB.UserID
	}
	return
}

//PlugIDs assign ids to Application users. this bulk version of PlugID
func (p *ApplicationUserRoleDao) PlugIDs(userRoles []coremodel.ApplicationUserRole, username string, ipAddress string, db *gorm.DB, logEntry *logrus.Entry) (userRolesFromDB []coremodel.ApplicationUserRole, err error) {
	var assignUUIDs []string
	idxByUUIDs := make(map[string]coremodel.ApplicationUserRole)
	for _, m := range userRoles {
		assignUUIDs = append(assignUUIDs, m.UUID)
		idxByUUIDs[m.UUID] = m
	}

	if dbRslt := db.Where("uid in ( ?) ", assignUUIDs).Find(&userRolesFromDB); dbRslt.Error != nil {
		err = dbRslt.Error
		logEntry.WithField("users", userRoles).WithField("userUuids", assignUUIDs).WithError(err).Errorf("Fail to query for user  by uuid, error: %s", err.Error())
		return
	}
	for _, m := range userRolesFromDB {
		if mnu, ok := idxByUUIDs[m.UUID]; ok {
			mnu.ID = m.ID
		}
	}
	return
}
