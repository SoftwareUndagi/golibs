package security

//SimpleUserData simple user data
type SimpleUserData struct {
	//ID internal id of user. if primary key is int32
	ID int32 `json:"id"`
	//UUID uuid dari user
	UUID string `json:"uid"`
	//Username username. ini berisi sama dengan UserUUID kalau misal anonymous
	Username string `json:"username"`
	//RealName dari column real_name
	RealName string `json:"realName"`
	//Email dari column: email
	Email string `json:"email"`
	//SovereignAuthID id of firebase or else
	SovereignAuthID string `json:"sovereignAuthId"`
	//Phone1 kalau user dengan phone auth , ini akan di ambil dari phone1
	Phone1 string `json:"phone1"`
	//Phone2 secondary phone
	Phone2 string `json:"phone2"`
	//UserRoles user roles. this to use block path with specified roles. for example user administration only allowed for admin user
	UserRoles []string `json:"userRoles"`
}
