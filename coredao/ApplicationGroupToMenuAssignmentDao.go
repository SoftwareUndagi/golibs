package coredao

import (
	"time"

	"github.com/SoftwareUndagi/golibs/coremodel"
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"
)

//ApplicationGroupToMenuAssignmentDao dao for ApplicationGroupToMenuAssignment
type ApplicationGroupToMenuAssignmentDao struct{}

//cachedApplicationGroupToMenuAssignment cache of struct. to simplify get table name or else
var cachedApplicationGroupToMenuAssignment = coremodel.ApplicationGroupToMenuAssignment{}

//sharedApplicationGroupToMenuAssignmentDao dao cache. for internal package
var sharedApplicationGroupToMenuAssignmentDao = ApplicationGroupToMenuAssignmentDao{}

//GetGroupMenus get menus of specified groups
func (p *ApplicationGroupToMenuAssignmentDao) GetGroupMenus(groupID int32, db *gorm.DB, logEntry *logrus.Entry) (menuAssignments []coremodel.ApplicationGroupToMenuAssignment, err error) {
	logEntry = logEntry.WithField("groupID", groupID)
	rslt := db.Where(&coremodel.ApplicationGroupToMenuAssignment{GroupID: groupID}).Find(&menuAssignments)
	if rslt.Error != nil {
		logEntry.WithError(rslt.Error).Errorf("Fail to query for group menus, error : %s", rslt.Error.Error())
		err = rslt.Error
		return
	}

	return
}

//SaveGroupMenus save group menus
func (p *ApplicationGroupToMenuAssignmentDao) SaveGroupMenus(groupID int32, menuAssignments []coremodel.ApplicationGroupToMenuAssignment, username string, userIPAddress string, db *gorm.DB, logEntry *logrus.Entry) (err error) {
	idxUpdateDataByUUID := make(map[int32]coremodel.ApplicationGroupToMenuAssignment)
	dbMenuAssign, errQuery := p.GetGroupMenus(groupID, db, logEntry) // query for existing assignment
	var rmMenuIDs []int32
	var addedMenuIds []int32
	if errQuery != nil {
		err = errQuery
		return
	}
	idxExistingDataByMenuID := make(map[int32]coremodel.ApplicationGroupToMenuAssignment)
	for _, p := range menuAssignments {
		idxUpdateDataByUUID[p.MenuID] = p
	}
	//index db data, mar deleted item
	for _, exDb := range dbMenuAssign {
		idxExistingDataByMenuID[exDb.MenuID] = exDb
		_, ok := idxUpdateDataByUUID[exDb.MenuID]
		if !ok {
			rmMenuIDs = append(rmMenuIDs, exDb.MenuID)
		}
	}
	upd := time.Now()
	for _, cUpd := range menuAssignments {
		dbRslt, ok := idxExistingDataByMenuID[cUpd.MenuID]
		if ok {
			dbRslt.AllowCreateFlag = cUpd.AllowCreateFlag
			dbRslt.AllowEditFlag = cUpd.AllowEditFlag
			dbRslt.AllowEraseFlag = cUpd.AllowEraseFlag
			dbRslt.ModifiedBy = &username
			dbRslt.ModifiedIPAddress = &userIPAddress
			dbRslt.UpdatedAt = &upd
			if updRslt := db.Save(&dbRslt); updRslt.Error != nil {
				logEntry.WithError(updRslt.Error).WithField("roleAssignmentId", dbRslt.ID).Errorf("Fail to update record on role assignment, error: %s", updRslt.Error.Error())
				err = updRslt.Error
				return
			}
			continue
		}
		uuidRslt, errUUID := uuid.NewUUID()
		if errUUID != nil {
			logEntry.WithError(errUUID).Errorf("Fail to generate UUID, error: %s", errUUID.Error())
			err = errUUID
			return
		}
		cUpd.CreatedAt = &upd
		cUpd.CreatorIPAddress = userIPAddress
		cUpd.CreatorName = username
		cUpd.UUID = uuidRslt.String()
		addedMenuIds = append(addedMenuIds, cUpd.MenuID)
		db.Create(&cUpd)
	}
	errUsafeAppGroup := sharedApplicationGroupDao.IncrementWithSpecifiedNumber(groupID, int16(len(addedMenuIds)-len(rmMenuIDs)), username, userIPAddress, db, logEntry)
	if errUsafeAppGroup != nil {
		err = errUsafeAppGroup
		return
	}
	if len(rmMenuIDs) > 0 {
		err = sharedApplicationMenuDao.DecrementUsageCount(rmMenuIDs, username, userIPAddress, db, logEntry)
		if err != nil {
			return
		}
	}
	if len(addedMenuIds) > 0 {
		err = sharedApplicationMenuDao.IncrementUsageCount(addedMenuIds, username, userIPAddress, db, logEntry)
		if err != nil {
			return
		}
	}
	return
}

//AddMenus add menus to group.
func (p *ApplicationGroupToMenuAssignmentDao) AddMenus(groupID int32, menus []coremodel.ApplicationGroupToMenuAssignment, username string, ipAddress string, db *gorm.DB, logEntry *logrus.Entry) (menuAssignments []coremodel.ApplicationGroupToMenuAssignment, err error) {
	if len(menus) == 0 {
		return
	}
	skr := time.Now()
	var menuIds []int32
	for idx, mnu := range menus {
		uuid, errUUID := uuid.NewUUID()
		if errUUID != nil {
			logEntry.WithError(errUUID).WithField("index", idx).Errorf("Error save menu on index %d, error :%s", idx, errUUID.Error())
			err = errUUID
			return
		}
		mnu.UUID = uuid.String()
		mnu.CreatedAt = &skr
		mnu.CreatorName = username
		mnu.CreatorIPAddress = ipAddress
		if rslt := db.Create(&mnu); rslt.Error != nil {
			err = rslt.Error
			logEntry.WithError(rslt.Error).WithField("menu", mnu).WithField("index", idx).Errorf("Fail to save data for index %d, error: %s", idx, rslt.Error.Error())
			return
		}
		menuAssignments = append(menuAssignments, mnu)
		menuIds = append(menuIds, mnu.MenuID)
	}
	err = sharedApplicationGroupDao.IncrementWithSpecifiedNumber(groupID, int16(len(menuAssignments)), username, ipAddress, db, logEntry)
	if err != nil {
		return
	}
	return menuAssignments, sharedApplicationMenuDao.IncrementUsageCount(menuIds, username, ipAddress, db, logEntry)
}

//RemoveMenus remove menus from group
func (p *ApplicationGroupToMenuAssignmentDao) RemoveMenus(menus []coremodel.ApplicationGroupToMenuAssignment, username string, ipAddress string, db *gorm.DB, logEntry *logrus.Entry) (err error) {
	if len(menus) == 0 {
		logEntry.Warnf("Parameter menus is empty. no menu removed from database")
		return
	}
	logEntry = logEntry.WithField("menus", menus)
	menuIDMap := make(map[int32]bool)
	groupIDMap := make(map[int32]bool)
	var menuIDs []int32
	var groupIDs []int32
	var assignIDs []int32
	for _, mnu := range menus {
		mnuID := mnu.MenuID
		grpID := mnu.GroupID
		if _, ok1 := menuIDMap[mnuID]; !ok1 {
			menuIDMap[mnuID] = true
			menuIDs = append(menuIDs, mnuID)
		}
		if _, ok2 := groupIDMap[grpID]; !ok2 {
			groupIDMap[grpID] = true
			groupIDs = append(groupIDs, grpID)
		}
		assignIDs = append(assignIDs, mnu.ID)
	}
	if dbrslt := db.Where("id in (?)").Delete(cachedApplicationGroupToMenuAssignment); dbrslt.Error != nil {
		err = dbrslt.Error
		logEntry.WithError(err).Errorf("Fail to delete menu assignments, error : %s", err.Error())
		return
	}
	//menuIDs := reflect.ValueOf(menuIDMap).MapKeys()
	if err = sharedApplicationMenuDao.DecrementUsageCount(menuIDs, username, ipAddress, db, logEntry); err != nil {
		return
	}
	if err = sharedApplicationGroupDao.DecrementUsageCount(groupIDs, username, ipAddress, db, logEntry); err != nil {
		return
	}
	return
}

//UpdateMenus update menus data. this did not do validation. just loop and run update statment
func (p *ApplicationGroupToMenuAssignmentDao) UpdateMenus(menus []coremodel.ApplicationGroupToMenuAssignment, username string, ipAddress string, db *gorm.DB, logEntry *logrus.Entry) (err error) {
	for idx, mnu := range menus {
		if rslt := db.Save(&mnu); rslt.Error != nil {
			err = rslt.Error
			logEntry.WithFields(logrus.Fields{"index": idx, "menu": mnu}).WithError(err).Errorf("Fail to update menu assignment, error: %s", err.Error())
			return
		}
	}
	return
}

//PlugID assign ID ( from db ) to struct.
//ID is not serialized ( on json field) from and to client. so when user do update to data, server will not recieve ID of ApplicationMenuAssignment, only UUID
// with this uuid string, id will be queried from server and plug back to data
func (p *ApplicationGroupToMenuAssignmentDao) PlugID(menu coremodel.ApplicationGroupToMenuAssignment, username string, ipAddress string, db *gorm.DB, logEntry *logrus.Entry) (menuFromDB coremodel.ApplicationGroupToMenuAssignment, err error) {
	var dbRslt *gorm.DB
	dbRslt = db.Where(&coremodel.ApplicationGroupToMenuAssignment{UUID: menu.UUID}).First(&menuFromDB)
	if dbRslt.Error != nil {
		err = dbRslt.Error
		logEntry.WithField("menuAssignment", menu).WithError(err).Errorf("Fail to get menu assignment. error: %s", err.Error())
		return
	}
	if !dbRslt.RecordNotFound() {
		menu.ID = menuFromDB.ID
	}
	return
}

//PlugIDs assign ids to menus. this bulk version of PlugID
func (p *ApplicationGroupToMenuAssignmentDao) PlugIDs(groupMenuassignments []coremodel.ApplicationGroupToMenuAssignment, username string, ipAddress string, db *gorm.DB, logEntry *logrus.Entry) (menuFromDB []coremodel.ApplicationGroupToMenuAssignment, err error) {
	var assignUUIDs []string
	idxByUUIDs := make(map[string]coremodel.ApplicationGroupToMenuAssignment)
	for _, m := range groupMenuassignments {
		assignUUIDs = append(assignUUIDs, m.UUID)
		idxByUUIDs[m.UUID] = m
	}

	if dbRslt := db.Where("uid in ( ?) ", assignUUIDs).Find(&menuFromDB); dbRslt.Error != nil {
		err = dbRslt.Error
		logEntry.WithField("menus", groupMenuassignments).WithField("menuUUIDs", assignUUIDs).WithError(err).Errorf("Fail to query for menu assignment by uuid, error: %s", err.Error())
		return
	}
	for _, m := range menuFromDB {
		if mnu, ok := idxByUUIDs[m.UUID]; ok {
			mnu.ID = m.ID
		}
	}
	return
}
