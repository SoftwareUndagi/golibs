package dao

import (
	"fmt"
	"regexp"
	"strings"
)

//regexNamedParameter  regex finder named parameter, :param1 etc
var regexNamedParameter = regexp.MustCompile(`(?m)(\040\072[a-zA-Z0-9]{1,})`)

//ConvertNamedParameterToPlainParameter convert from named parameter to plain query . for example :
// where 1=1
//		and type = :employeeType
//		and (address like :pattern or office_address like :pattern)
//		and age > :ageParam
//
// with param : {"type": "permanent" , "ageParam": 35 , "pattern": "%denpasar%"}
// query will convert to :
// where 1=1
//		and type = ?
//		and (address like ? or office_address like ?)
//		and age > ?
// and flatedQueryParameter = ["permanent" ,  "%denpasar%" ,"%denpasar%" , 35 ]
func ConvertNamedParameterToPlainParameter(queryWithNamedParameter string, queryParameter map[string]interface{}) (convertedQuery string, flatedQueryParameter []interface{}) {
	matches := regexNamedParameter.FindAllString(queryWithNamedParameter, -1)
	if len(matches) == 0 {
		convertedQuery = queryWithNamedParameter
		return
	}
	replaceableString := make([]string, 0)
	replaceableKey := make(map[string]bool)
	for _, mtch := range matches {
		paramCurrent := strings.Trim(mtch, " ")
		plain := paramCurrent[1:]
		if val, ok := queryParameter[plain]; ok {
			flatedQueryParameter = append(flatedQueryParameter, val)
			replaceableKey[plain] = true
			replaceableString = append(replaceableString, mtch)
		}
	}
	var replcRegex string
	for key := range replaceableKey {
		if len(replcRegex) > 0 {
			replcRegex = replcRegex + "|"
		}
		replcRegex = replcRegex + fmt.Sprintf("(\\040\\072%v)", key)
	}
	rplRegx := regexp.MustCompile(replcRegex)
	convertedQuery = rplRegx.ReplaceAllString(queryWithNamedParameter, " ? ")
	return
}
