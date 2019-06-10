package coremodel

import (
	"time"

	"github.com/SoftwareUndagi/golibs/common"
	"github.com/jinzhu/gorm"
)

//tablenameUserToGroupAssignment nama table . di constant untuk optimasi
const tablenameUserToGroupAssignment = "sec_group_assignment"

//ApplicationUserToGroupAssignment table: sec_group_assignment
type ApplicationUserToGroupAssignment struct {
	//ID id data(surrogate key), column: pk
	ID int32 `gorm:"column:pk;AUTO_INCREMENT;primary_key" json:"id"`
	//GroupID id dari group, column: group_id
	GroupID int32 `gorm:"column:group_id" json:"groupId"`
	//UserID id dari user, column: user_id
	UserID int32 `gorm:"column:user_id" json:"userId"`
	//User refer with column user_id
	User ApplicationUser `gorm:"foreignkey:UserID" json:"user"`
	//Group refer with column group_id
	Group ApplicationGroup `gorm:"foreignkey:GroupID" json:"group"`
	//Creator column: creator_name
	CreatorName string `gorm:"column:creator_name" json:"creator"`
	//CreatorIpAddress column: creator_ip_address
	CreatorIPAddress string `gorm:"column:creator_ip_address" json:"creatorIpAddress"`
	//CreatedAt column : createdAt time when data was created
	CreatedAt *time.Time `gorm:"column:createdAt" json:"createdAt"`
}

//TableName table name for struct UserToGroupAssignment
func (u ApplicationUserToGroupAssignment) TableName() string {
	return tablenameUserToGroupAssignment
}

//BeforeCreate before create task. to assign IP address and username on data
func (u ApplicationUserToGroupAssignment) BeforeCreate(scope *gorm.Scope) (err error) {
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
