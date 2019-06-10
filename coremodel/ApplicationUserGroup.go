package coremodel

import (
	"github.com/jinzhu/gorm"
)

const applicationUserGroupTableName = "sec_group_assignment"

//ApplicationUserGroup wrapper for table sec_group_assignment
type ApplicationUserGroup struct {
	//ID for column pk .
	ID int32 `gorm:"column:pk;AUTO_INCREMENT;primary_key" json:"-"`
	//UUID business user visible
	UUID string `gorm:"column:uuid" json:"uid"`

	//GroupID group_id
	GroupID int32 `gorm:"column:group_id" json:"groupId"`
	//UserID column user_id
	UserID int32 `gorm:"column:user_id" json:"userId"`
	//User reference to user
	User ApplicationUser `gorm:"foreignKey:UserID" json:"user"`
	//Group group of user
	Group ApplicationGroup `gorm:"foreignKey:GroupID" json:"group"`
}

//TableName access to table name for sec_group_assignment
func (p *ApplicationUserGroup) TableName(db *gorm.DB) string {
	return applicationUserGroupTableName
}
