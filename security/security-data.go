package security

//SimpleUserData simple user data
type SimpleUserData struct {
	//ID internal id of user. if primary key is int32
	ID int32
	//UUID uuid dari user
	UUID string
	//Username username. ini berisi sama dengan UserUUID kalau misal anonymous
	Username string
	//RealName dari column real_name
	RealName string
	//Email dari column: email
	Email string
	//UserUUID dari column uuid
	UserUUID string
	//Phone kalau user dengan phone auth , ini akan di ambil dari phone1
	Phone string
	//UserRoles user roles. this to use block path with specified roles. for example user administration only allowed for admin user
	UserRoles []string
}
