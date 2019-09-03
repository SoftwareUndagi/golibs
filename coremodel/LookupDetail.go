package coremodel

import (
	"time"

	"github.com/SoftwareUndagi/golibs/common"
	"github.com/jinzhu/gorm"
)

//LookupDetail Simple table lookup.for table : m_lookup_details
type LookupDetail struct {
	//ID column id
	ID int32 `gorm:"column:id;AUTO_INCREMENT;primary_key" json:"id"`
	//DetailCode column: detail_code kode detail
	DetailCode string `gorm:"column:detail_code" json:"detailCode"`
	//LovID column: lov_id
	LovID string `gorm:"column:lov_id" json:"lovId"`
	//Label label for lookup column: lov_label
	Label string `gorm:"column:label" json:"label"`
	//Value1 label for value 1. arbitary data 1
	Value1 *string `gorm:"column:val_1" json:"value11"`
	//Value2 label for value 2. arbitary data 2
	Value2 *string `gorm:"column:val_2" json:"value12"`
	//I18nKey key internalization for lookup
	I18nKey *string `gorm:"column:i18n_key" json:"i18nKey"`
	//SequenceNo sort no for lookup
	SequenceNo *int16 `gorm:"column:seq_no" json:"sequenceNo"`

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

//TableName table name for m_lookup_details
func (u *LookupDetail) TableName(db *gorm.DB) (name string) {
	return "ct_lookup_details"
}

//BeforeCreate before create task. to assign IP address and username on data
func (u *LookupDetail) BeforeCreate(scope *gorm.Scope) (err error) {
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
func (u *LookupDetail) BeforeUpdate(scope *gorm.Scope) (err error) {
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
