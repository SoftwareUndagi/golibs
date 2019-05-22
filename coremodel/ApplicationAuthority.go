package coremodel

import (
	"time"

	"github.com/jinzhu/gorm"
)

//applicationAuthorityTableName constant table name sec_authority
const applicationAuthorityTableName = "sec_authority"

//ApplicationAuthority struct for table : sec_authority
type ApplicationAuthority struct {
	//Code column: code. code of authority
	Code string `gorm:"column:code;primary_key" json:"code"`
	//Remark remark for group column: remark
	Remark string `gorm:"column:remark" json:"remark"`
	//CreatedAt column : createdAt time when data was created
	CreatedAt *time.Time `gorm:"column:createdAt" json:"createdAt"`
	//CreatorName username (audit trail), who create data
	CreatorName string `gorm:"column:creator_name" json:"creatorName"`
	//CreatorIPAddress IP address(audit trail), from which IP address data created
	CreatorIPAddress string `gorm:"column:creator_name" json:"creatorIpAddress"`

	//UpdatedAt last update at column : updatedAt
	UpdatedAt *time.Time `gorm:"column:updatedAt" json:"updatedAt"`
	//ModifiedBy audit trail. latest user that modify data
	ModifiedBy *string `gorm:"column:modified_by" json:"modifiedBy"`
	//ModifiedIPAddress IP address from where user modify data(latest update)
	ModifiedIPAddress string `gorm:"column:modified_by_ip" json:"modifiedIpAddress"`
}

//TableName table name for struct ApplicationGroup
func (u ApplicationAuthority) TableName(db *gorm.DB) string {
	return applicationAuthorityTableName
}
