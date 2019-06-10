package coredao

import (
	"fmt"

	"github.com/SoftwareUndagi/golibs/coremodel"
	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"
)

//ApplicationMenuDao dao for application menu
type ApplicationMenuDao struct{}

var sampleApplicationMenu = coremodel.ApplicationMenu{}
var sharedApplicationMenuDao = ApplicationMenuDao{}

//IncrementUsageCount increase usage count on table sec_menu
func (p *ApplicationMenuDao) IncrementUsageCount(menuIDs []int32, username string, ipAddress string, db *gorm.DB, logEntry *logrus.Entry) (err error) {
	if len(menuIDs) == 0 {
		return
	}
	updSQL := fmt.Sprintf(`update %s 
	set 
		usage_count =  ifnull(usage_count,0) + 1  , 
		modified_by = ? ,
		modified_by_ip = ? , 
		updatedAt = current_timestamp
	where 
		code in ( ?) `, sampleApplicationMenu.TableName(db))
	if dbRslt := db.Exec(updSQL, username, ipAddress, menuIDs); dbRslt.Error != nil {
		logEntry.WithError(dbRslt.Error).WithField("menuIDs", menuIDs).Errorf("Fail to increment Sec group , error: %s", dbRslt.Error.Error())
		return dbRslt.Error
	}
	return
}

//DecrementUsageCount decrement usage count on table sec_group
func (p *ApplicationMenuDao) DecrementUsageCount(menuIDs []int32, username string, ipAddress string, db *gorm.DB, logEntry *logrus.Entry) (err error) {
	if len(menuIDs) == 0 {
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
		and ifnull(usage_count,0) > 0 `, sampleApplicationMenu.TableName(db))
	if dbRslt := db.Exec(updSQL, username, ipAddress, menuIDs); dbRslt.Error != nil {
		logEntry.WithError(dbRslt.Error).WithField("menuIDs", menuIDs).Errorf("Fail to decrement sec group usage , error: %s", dbRslt.Error.Error())
		return dbRslt.Error
	}
	return
}

//PlugID assign ID ( from db ) to struct.
//ID is not serialized ( on json field) from and to client. so when user do update to data, server will not recieve ID of ApplicationGroup, only UUID
// with this uuid string, id will be queried from server and plug back to data
func (p *ApplicationMenuDao) PlugID(menu coremodel.ApplicationMenu, username string, ipAddress string, db *gorm.DB, logEntry *logrus.Entry) (menuFromDB coremodel.ApplicationMenu, err error) {
	var dbRslt *gorm.DB
	dbRslt = db.Where(&coremodel.ApplicationMenu{UUID: menu.UUID}).First(&menuFromDB)
	if dbRslt.Error != nil {
		err = dbRslt.Error
		logEntry.WithField("menu", menu).WithError(err).Errorf("Fail to Group. error: %s", err.Error())
		return
	}
	if !dbRslt.RecordNotFound() {
		menu.ID = menuFromDB.ID
		
	}
	return
}

//PlugIDs assign ids to Application groups. this bulk version of PlugID
func (p *ApplicationMenuDao) PlugIDs(menus []coremodel.ApplicationMenu, username string, ipAddress string, db *gorm.DB, logEntry *logrus.Entry) (menusFromDB []coremodel.ApplicationMenu, err error) {
	var assignUUIDs []string
	idxByUUIDs := make(map[string]coremodel.ApplicationMenu)
	for _, m := range menus {
		assignUUIDs = append(assignUUIDs, m.UUID)
		idxByUUIDs[m.UUID] = m
	}

	if dbRslt := db.Where("uid in ( ?) ", assignUUIDs).Find(&menusFromDB); dbRslt.Error != nil {
		err = dbRslt.Error
		logEntry.WithField("menus", menus).WithField("menuUuids", assignUUIDs).WithError(err).Errorf("Fail to query for menu assignment by uuid, error: %s", err.Error())
		return
	}
	for _, m := range menusFromDB {
		if mnu, ok := idxByUUIDs[m.UUID]; ok {
			mnu.ID = m.ID
		}
	}
	return
}
