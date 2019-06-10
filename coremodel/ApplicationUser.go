package coremodel

import (
	"time"

	"github.com/SoftwareUndagi/golibs/common"
	"github.com/jinzhu/gorm"
)

//tablenameApplicationUser nama table . di constant untuk optimasi
const tablenameApplicationUser = "sec_user"

//ApplicationUser table: sec_user
type ApplicationUser struct {
	//ID primary key, column: pk
	ID int32 `gorm:"column:id;AUTO_INCREMENT;primary_key" json:"-"`
	//SovereignAuthID id of firebase or else
	SovereignAuthID string `gorm:"column:soverign_auth_id" json:"sovereignAuthId"`
	//Uuid UUID, untuk integrasi misal dengan firebase, column: uuid
	UUID string `gorm:"column:uuid" json:"uid"`
	//RealName nama lengkap, column: real_name
	RealName string `gorm:"column:real_name" json:"realName"`
	//Username username untuk login, column: user_name
	Username string `gorm:"column:user_name" json:"username"`
	//AvatarURL url avatar dari user, column: avatar_url
	AvatarURL string `gorm:"column:avatar_url" json:"avatarUrl"`
	//Password password user(kalau pakai metode password standard), column: password
	Password string `gorm:"column:password" json:"-"`
	//DefaultBranchCode kode cabang default(kode unit kerja - di abaikan kalau PNS), column: default_branch_code
	DefaultBranchCode string `gorm:"column:default_branch_code" json:"defaultBranchCode"`
	//Email email dari user, column: email
	Email string `gorm:"column:email" json:"email"`
	//EmployeeNo reference ke data pegawai, column: emp_no
	EmployeeNo string `gorm:"column:emp_no" json:"employeeNo"`
	//ExpiredDate waktu expired dari user, column: expired_date
	ExpiredDate time.Time `gorm:"column:expired_date" json:"expiredDate"`
	//FailedLoginCount count berapa kali gagal login.counter. kalau melebihi max, maka di lock, column: failed_login_attemps
	FailedLoginCount int32 `gorm:"column:failed_login_attemps" json:"failedLoginCount"`
	//LocaleCode locale code, column: locale
	LocaleCode string `gorm:"column:locale_code" json:"localeCode"`
	//LockedFlag flag di lock atau tidak, column: is_locked
	LockedFlag string `gorm:"column:is_locked" json:"lockedFlag"`
	//ActiveFlag flag is user ative( Y or N)
	ActiveFlag string `gorm:"column:is_active" json:"activeFlag"`
	//Phone1 phone1, column: phone1
	Phone1 string `gorm:"column:phone1" json:"phone1"`
	//Phone2 phone secondary, column: phone2
	Phone2 string `gorm:"column:phone2" json:"phone2"`
	//Remark catatan, column: remark
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
	UpdatedAt *time.Time `gorm:"column:updatedAt" json:"updatedAt"`
	//ModifiedBy audit trail. latest user that modify data
	ModifiedBy *string `gorm:"column:modified_by" json:"modifiedBy"`
	//ModifiedIPAddress IP address from where user modify data(latest update)
	ModifiedIPAddress *string `gorm:"column:modified_by_ip" json:"modifiedIpAddress"`
}

//TableName table name for struct ApplicationUser
func (u ApplicationUser) TableName(db *gorm.DB) string {
	return tablenameApplicationUser
}

//sampleUsernameStruct sample user. Cache helper to get table name
var sampleUsernameStruct = ApplicationUser{}

//BeforeCreate before create task. to assign IP address and username on data
func (u ApplicationUser) BeforeCreate(scope *gorm.Scope) (err error) {
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
func (u ApplicationUser) BeforeUpdate(scope *gorm.Scope) (err error) {
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
