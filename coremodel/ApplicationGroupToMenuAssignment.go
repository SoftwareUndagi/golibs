package coremodel

import (
	"time"

	"github.com/SoftwareUndagi/golibs/common"
	"github.com/jinzhu/gorm"
)

//tablenameGroupToMenuAssignment nama table . di constant untuk optimasi
const tablenameGroupToMenuAssignment = "sec_menu_assignment"

//ApplicationGroupToMenuAssignment table: sec_menu_assignment
type ApplicationGroupToMenuAssignment struct {
	//ID primary key data, column: id
	ID int32 `gorm:"column:id;AUTO_INCREMENT;primary_key" json:"-"`
	//UUID uuid of data, user visible id. column=  uuid
	UUID string `gorm:"column:uuid" json:"uid"`
	//MenuID id dari menu, column: menu_id
	MenuID int32 `gorm:"column:menu_id" json:"-"`
	//Menu reference to ApplicationMenu(over column: menu_id)
	Menu ApplicationMenu `gorm:"foreignkey:MenuID" json:"menu"`
	//GroupID id group, column: group_id
	GroupID int32 `gorm:"column:group_id" json:"-"`
	//Group reference to user group
	Group ApplicationGroup `gorm:"foreignkey:GroupID" json:"group"`
	//AllowCreateFlag flag : ijinkan new data, column: is_allow_create
	AllowCreateFlag string `gorm:"column:is_allow_create" json:"allowCreateFlag"`
	//AllowEditFlag flag : ijinkan edit data, column: is_allow_edit
	AllowEditFlag string `gorm:"column:is_allow_edit" json:"allowEditFlag"`
	//AllowEraseFlag flag : ijinkan hapus data, column: is_allow_erase
	AllowEraseFlag string `gorm:"column:is_allow_erase" json:"allowEraseFlag"`
	//CreatedAt column : createdAt time when data was created
	CreatedAt *time.Time `gorm:"column:createdAt" json:"createdAt"`
	//CreatorName username (audit trail), who create data
	CreatorName string `gorm:"column:creator_name" json:"creatorName"`
	//CreatorIPAddress IP address(audit trail), from which IP address data created
	CreatorIPAddress string `gorm:"column:creator_ip_address" json:"creatorIpAddress"`
	//UpdatedAt last update at column : updatedAt
	UpdatedAt *time.Time `gorm:"column:updatedAt" json:"updatedAt"`
	//ModifiedBy audit trail. latest user that modify data
	ModifiedBy *string `gorm:"column:modified_by" json:"modifiedBy"`
	//ModifiedIPAddress IP address from where user modify data(latest update)
	ModifiedIPAddress *string `gorm:"column:modified_by_ip" json:"modifiedIpAddress"`
}

//TableName table name for struct GroupToMenuAssignment
func (u ApplicationGroupToMenuAssignment) TableName() string {
	return tablenameGroupToMenuAssignment
}

//BeforeCreate before create task. to assign IP address and username on data
func (u ApplicationGroupToMenuAssignment) BeforeCreate(scope *gorm.Scope) (err error) {
	if len(u.CreatorName) == 0 {
		if uname, okName := scope.Get(common.GormVariableUsername); okName {
			u.CreatorName = uname.(string)
		}
	}
	if len(u.CreatorIPAddress) == 0 {
		if ipAddr, okIP := scope.Get(common.GormVariableIPAddress); okIP {
			u.CreatorIPAddress = ipAddr.(string)
		}
	}
	return nil
}

//BeforeUpdate task before update
func (u ApplicationGroupToMenuAssignment) BeforeUpdate(scope *gorm.Scope) (err error) {
	if u.ModifiedBy == nil || len(*u.ModifiedBy) == 0 {
		if uname, okName := scope.Get(common.GormVariableUsername); okName {
			strUname := uname.(string)
			u.ModifiedBy = &strUname
		}
	}
	if u.ModifiedIPAddress == nil || len(*u.ModifiedIPAddress) == 0 {
		if ipAddr, okIP := scope.Get(common.GormVariableIPAddress); okIP {
			strIP := ipAddr.(string)
			u.ModifiedIPAddress = &strIP
		}
	}
	return nil
}
