package coremodel

import "time"

//tablenameGroupToMenuAssignment nama table . di constant untuk optimasi
const tablenameGroupToMenuAssignment = "sec_menu_assignment"

//ApplicationGroupToMenuAssignment table: sec_menu_assignment
type ApplicationGroupToMenuAssignment struct {
	//ID primary key data, column: id
	ID int64 `gorm:"column:id;AUTO_INCREMENT;primary_key" json:"id"`
	//MenuID id dari menu, column: menu_id
	MenuID int64 `gorm:"column:menu_id" json:"menuId"`
	//GroupID id group, column: group_id
	GroupID int64 `gorm:"column:group_id" json:"groupId"`
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
	CreatorIPAddress string `gorm:"column:creator_name" json:"creatorIpAddress"`
	//UpdatedAt last update at column : updatedAt
	UpdatedAt *time.Time `gorm:"column:updatedAt" json:"updatedAt"`
	//ModifiedBy audit trail. latest user that modify data
	ModifiedBy *string `gorm:"column:modified_by" json:"modifiedBy"`
	//ModifiedIPAddress IP address from where user modify data(latest update)
	ModifiedIPAddress string `gorm:"column:modified_by_ip" json:"modifiedIpAddress"`
}

//TableName table name for struct GroupToMenuAssignment
func (u ApplicationGroupToMenuAssignment) TableName() string {
	return tablenameGroupToMenuAssignment
}
