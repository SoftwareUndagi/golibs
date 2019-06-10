package coremodel

import (
	"time"

	"github.com/SoftwareUndagi/golibs/common"
	"github.com/jinzhu/gorm"
)

//tablenameApplicationGroup nama table . di constant untuk optimasi
const tablenameApplicationGroup = "sec_group"

//ApplicationGroup table: sec_group
type ApplicationGroup struct {
	//ID id dari group, column: id
	ID int32 `gorm:"column:id;AUTO_INCREMENT;primary_key" json:"-"`
	//UUID UUID of security group
	UUID string `gorm:"column:uuid" json:"uid"`
	//Code kode group, column: group_code
	Code string `gorm:"column:code" json:"code"`
	//Name nama group, column: group_name
	Name string `gorm:"column:name" json:"name"`
	//Remark catatan dari group, column: group_remark
	Remark string `gorm:"column:remark" json:"remark"`
	//UsageCounter berapa data yang sudah merefer ini , column: usage_count
	UsageCounter int16 `gorm:"column:usage_count" json:"usageCounter"`
	//CreatedAt column : createdAt time when data was created
	CreatedAt *time.Time `gorm:"column:createdAt" json:"createdAt"`
	//CreatorName username (audit trail), who create data
	CreatorName string `gorm:"column:creator_name" json:"creatorName"`
	//CreatorIPAddress IP address(audit trail), from which IP address data created
	CreatorIPAddress string `gorm:"column:creator_ip_address" json:"creatorIpAddress"`

	//UpdatedAt last update at column : updatedAt
	UpdatedAt *time.Time `gorm:"column:updatedAt;nullable" json:"updatedAt"`
	//ModifiedBy audit trail. latest user that modify data
	ModifiedBy *string `gorm:"column:modified_by" json:"modifiedBy"`
	//ModifiedIPAddress IP address from where user modify data(latest update)
	ModifiedIPAddress *string `gorm:"column:modified_by_ip" json:"modifiedIpAddress"`
}

//TableName table name for struct ApplicationGroup
func (u *ApplicationGroup) TableName(db *gorm.DB) string {
	return tablenameApplicationGroup
}

//BeforeCreate before create task. to assign IP address and username on data
func (u *ApplicationGroup) BeforeCreate(scope *gorm.Scope) (err error) {
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
func (u *ApplicationGroup) BeforeUpdate(scope *gorm.Scope) (err error) {
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
