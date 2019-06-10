package coremodel

import (
	"time"

	"github.com/SoftwareUndagi/golibs/common"
	"github.com/jinzhu/gorm"
)

//aApplicationUserRoleTableName constant table name sec_user_role
const aApplicationUserRoleTableName = "sec_user_role"

//ApplicationUserRole mapper for table : sec_user_authority
type ApplicationUserRole struct {
	//ID column ID
	ID int32 `gorm:"column:id;primary_key;AUTO_INCREMENT" json:"-"`
	//UUID UUID of security group
	UUID string `gorm:"column:uuid" json:"uid"`
	//UserID column: user_id id of user
	UserID int32 `gorm:"column:user_id" json:"-"`

	//User refer with column: user_id
	User ApplicationUser `gorm:"foreignKey:UserID" json:"user"`

	//AuthCode column: authority_code, user auth code
	AuthCode string `gorm:"column:authority_code" json:"authCode"`
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

//TableName table of struct
func (u ApplicationUserRole) TableName(db *gorm.DB) string {
	return aApplicationUserRoleTableName
}

//sampleAppAuthStruct helper. for reuse, when need to get table name.in case used in multiple scheme
var sampleAppAuthStruct = ApplicationUserRole{}

//BeforeCreate before create task. to assign IP address and username on data
func (u ApplicationUserRole) BeforeCreate(scope *gorm.Scope) (err error) {
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
func (u ApplicationUserRole) BeforeUpdate(scope *gorm.Scope) (err error) {
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
