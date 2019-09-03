package coredao

import (
	"fmt"
	"reflect"
	"testing"
)

type sampleSetStruct struct {
	//ID id dari data
	ID int64 `gorm:"column:id;AUTO_INCREMENT"`
	//Name name of
	Name string `gorm:"column:name"`
	//Email email of user
	Email string `gorm:"column:email"`
}

func assignIds(prms []interface{}) {
	for index, p := range prms {
		swap := reflect.ValueOf(p)
		s := swap.Elem()
		s.FieldByName("ID").SetInt(int64(index + 1))
	}
}

func Test_assign_struct_value(t *testing.T) {
	var smplData []interface{}
	smplData = append(smplData, &sampleSetStruct{Name: "gede", Email: "gede.sutarsa@gmail.com"})
	smplData = append(smplData, &sampleSetStruct{Name: "Yuli", Email: "yuli.dana@gmail.com"})
	assignIds(smplData)
	for _, r := range smplData {
		fmt.Printf("Name : %s , Id : %d", r.(*sampleSetStruct).Name, r.(*sampleSetStruct).ID)
	}
	fmt.Printf("selesai")
}
