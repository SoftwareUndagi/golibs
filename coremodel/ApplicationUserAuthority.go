package coremodel

import (
	"time"

	"github.com/jinzhu/gorm"
)

//applicationUserAuthority constant table name sec_user_authority
const applicationUserAuthority = "sec_user_authority"

//ApplicationUserAuthority mapper for table : sec_user_authority
type ApplicationUserAuthority struct {
	//ID column ID
	ID int32 `gorm:"column:id;primary_key;AUTO_INCREMENT" json:"id"`
	//UserID column: user_id id of user
	UserID int32 `gorm:"column:user_id" json:"userId"`

	//AuthCode column: authority_code, user auth code
	AuthCode string `gorm:"column:authority_code" json:"authCode"`
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

//TableName table of struct
func (p ApplicationUserAuthority) TableName(db *gorm.DB) string {
	return applicationUserAuthority
}

//sampleAppAuthStruct helper. for reuse, when need to get table name.in case used in multiple scheme
var sampleAppAuthStruct = ApplicationUserAuthority{}
