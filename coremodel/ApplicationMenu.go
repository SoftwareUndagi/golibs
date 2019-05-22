package coremodel

import "time"

//tablenameApplicationMenu nama table . di constant untuk optimasi
const tablenameApplicationMenu = "sec_menu"

//ApplicationMenu table: sec_menu
type ApplicationMenu struct {
	//ID id dari menu, column: id
	ID int32 `gorm:"column:id;AUTO_INCREMENT;primary_key" json:"id"`
	//Code kode menu, column: menu_code
	Code string `gorm:"column:menu_code" json:"code"`
	//ParentID id induk dari menu, column: parent_id
	ParentID int64 `gorm:"column:parent_id" json:"parentId"`
	//Label label dari menu, column: menu_label
	Label string `gorm:"column:menu_label" json:"label"`
	//MenuTreeCode tree code dari menu, column: menu_tree_code
	MenuTreeCode string `gorm:"column:menu_tree_code" json:"menuTreeCode"`
	//OrderNumber urutan data, column: order_no
	OrderNumber int16 `gorm:"column:order_no" json:"orderNumber"`
	//I18nKey key internalization, column: i18n_key
	I18nKey string `gorm:"column:i18n_key" json:"i18nKey"`
	//RoutePath path dari handler, column: route_path
	RoutePath string `gorm:"column:route_path" json:"routePath"`
	//AdditionalParameter ,, column: additional_param
	AdditionalParameter string `gorm:"column:additional_param" json:"additionalParameter"`
	//StatusCode status data, column: data_status
	StatusCode string `gorm:"column:data_status" json:"statusCode"`
	//TreeLevelPosition level menu. pada level berapa data berada, column: tree_level_position
	TreeLevelPosition int16 `gorm:"column:tree_level_position" json:"treeLevelPosition"`
	//UsageCounter berapa data yang sudah merefer ini , column: usage_count
	UsageCounter int16 `gorm:"column:usage_count" json:"usageCounter"`
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

//TableName table name for struct ApplicationMenu
func (u ApplicationMenu) TableName() string {
	return tablenameApplicationMenu
}
