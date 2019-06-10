package coremodel

import (
	"time"

	"github.com/SoftwareUndagi/golibs/common"
	"github.com/jinzhu/gorm"
)

//LookupHeader table m_lookup_header
type LookupHeader struct {
	//ID code/ id of lookup column : lov_id
	ID string `gorm:"column:id;primary_key" json:"id" `
	//Reamark remark for lookup. for maintence purpose. column lov_remark
	Remark string `gorm:"column:remark" json:"remark"`
	//FlagCacheable Y or N. flag if lookup is cacheable on client
	FlagCacheable string `gorm:"column:is_cacheable" json:"flagCachable"`
	//FlagUseCustomSQL Y and N flag. is lookup using custom sql or from lookup detail
	FlagUseCustomSQL string `gorm:"column:is_use_custom_sql" json:"flagUseCustomSql"`
	//Version version of lookup. to force client reload if cache is expired
	Version string `gorm:"column:version" json:"version"`
	//SQLForData sql for check LOV version. lookup version based on query result
	SQLForData *string `gorm:"column:sql_data" json:"sqlForData"`
	//SQLForDataFiltered sql for data filtered(id is passed with prev data)
	SQLForDataFiltered string `gorm:"column:sql_data_filtered" json:"sqlForDataFiltered"`
	//SQLForVersion sql for version of lookup
	SQLForVersion string `gorm:"column:sql_version" json:"sqlForVersion"`
	//Details list of LookupDetail under this lookup header
	Details *[]LookupDetail `gorm:"-" json:"details"`
	//CodeActualDataType actual data type for code on lookup detail . posible value is : string , number. this value is use to generate in statement on sql
	CodeActualDataType string `gorm:"column:code_actual_data_type" json:"codeActualDataType"`

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

//TableName table name for m_lookup_header
func (p *LookupHeader) TableName(db *gorm.DB) (name string) {
	return "ct_lookup_header"
}

//AppendLookupDetail append lookup detail to header
func (p *LookupHeader) AppendLookupDetail(detail LookupDetail) {
	if p.Details == nil {
		dtl := []LookupDetail{}
		p.Details = &dtl
	}
	swap := append(*p.Details, detail)
	p.Details = &swap
}

//BeforeUpdate hook pre create/ update
func (p *LookupHeader) BeforeUpdate(scope *gorm.Scope) (err error) {
	if len(p.FlagUseCustomSQL) == 0 {
		if len(*p.SQLForData) > 0 {
			p.FlagUseCustomSQL = "Y"
		} else {
			p.FlagUseCustomSQL = "N"
		}

	}
	if len(p.FlagCacheable) == 0 {
		p.FlagCacheable = "Y"
	}
	if len(p.Version) == 0 {
		p.Version = "001"
	}
	if p.ModifiedBy == nil || len(*p.ModifiedBy) == 0 {
		if uname, okName := scope.Get(common.GormVariableUsername); okName {
			strUname := uname.(string)
			p.ModifiedBy = &strUname
		}
	}
	if p.ModifiedIPAddress == nil || len(*p.ModifiedIPAddress) == 0 {
		if ipAddr, okIP := scope.Get(common.GormVariableIPAddress); okIP {
			strIP := ipAddr.(string)
			p.ModifiedIPAddress = &strIP
		}
	}
	return
}
