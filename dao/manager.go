package dao

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/SoftwareUndagi/golibs/common"
	"github.com/SoftwareUndagi/golibs/coremodel"
	"github.com/SoftwareUndagi/golibs/dao/query"
	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"
)

//DefaultDaoManager default dai manager. for default on app
var DefaultDaoManager = NewManager()

//GetDaoManager get manager current
func GetDaoManager() Manager {
	return DefaultDaoManager
}

//RegisterCoreModel register model on common-lib
func RegisterCoreModel() {
	DefaultDaoManager.RegisterModel(
		GeneratorDefinition{
			Generator:      func() interface{} { return &coremodel.LookupHeader{} },
			SliceGenerator: func() interface{} { return &[]coremodel.LookupHeader{} }},
		GeneratorDefinition{
			Generator:      func() interface{} { return &coremodel.LookupDetail{} },
			SliceGenerator: func() interface{} { return &[]coremodel.LookupDetail{} }},
		GeneratorDefinition{
			Generator:      func() interface{} { return &coremodel.ApplicationGroup{} },
			SliceGenerator: func() interface{} { return &[]coremodel.ApplicationGroup{} }},
		GeneratorDefinition{
			Generator:      func() interface{} { return &coremodel.ApplicationMenu{} },
			SliceGenerator: func() interface{} { return &[]coremodel.ApplicationMenu{} }},
		GeneratorDefinition{
			Generator:      func() interface{} { return &coremodel.ApplicationUser{} },
			SliceGenerator: func() interface{} { return &[]coremodel.ApplicationUser{} }},
		GeneratorDefinition{
			Generator:      func() interface{} { return &coremodel.ApplicationGroupToMenuAssignment{} },
			SliceGenerator: func() interface{} { return &[]coremodel.ApplicationGroupToMenuAssignment{} }},
		GeneratorDefinition{
			Generator:      func() interface{} { return &coremodel.ApplicationUserGroup{} },
			SliceGenerator: func() interface{} { return &[]coremodel.ApplicationUserGroup{} }},
		GeneratorDefinition{
			Generator:      func() interface{} { return &coremodel.ApplicationUserToGroupAssignment{} },
			SliceGenerator: func() interface{} { return &[]coremodel.ApplicationUserToGroupAssignment{} }})
}

//OperatorCode code for query operator. for parsing json param query
type OperatorCode string

const (
	//OprEq code for eq( = )
	OprEq OperatorCode = "$eq"
	//OprGt code for gt( > )
	OprGt OperatorCode = "$gt"
	//OprGte code for gte( >= )
	OprGte OperatorCode = "$gte"
	//OprIn code for in
	OprIn OperatorCode = "$in"
	//OprIsNot code for is not
	OprIsNot OperatorCode = "$isNot"
	//OprIs  code for is
	OprIs OperatorCode = "$is"
	//OprLike  code for like
	OprLike OperatorCode = "$like"
	//OprNotLike parse parse operaator not like
	OprNotLike OperatorCode = "$notLike"
	//OprLt  code for <( less then )
	OprLt OperatorCode = "$lt"
	//OprLte  code for <=( less then equals)
	OprLte OperatorCode = "$lte"
	//OprNe  code for != ( not equals)
	OprNe OperatorCode = "$ne"
	//OprOr code for or
	OprOr OperatorCode = "$or"
)

//OperatorCodeMap key = OperatorCode , value = string
var OperatorCodeMap = map[OperatorCode]string{
	OprEq:      string(OprEq),
	OprNe:      string(OprNe),
	OprLte:     string(OprLte),
	OprLt:      string(OprLt),
	OprLike:    string(OprLike),
	OprNotLike: string(OprNotLike),
	OprIsNot:   string(OprIsNot),
	OprIs:      string(OprIs),
	OprIn:      string(OprIn),
	OprGte:     string(OprGte),
	OprGt:      string(OprGt),
	OprOr:      string(OprOr)}

//StartWith check apah string start with code
func (p OperatorCode) StartWith(testedQuery string) (result bool) {
	return strings.HasPrefix(testedQuery, OperatorCodeMap[p])
}

//usageCountQueryByDialect wraper,as of db dialect. key = table name
type usageCountQueryByDialect map[string][]string

//usageCountQuery usae count query,key  = dialect of database
type usageCountQuery map[string]usageCountQueryByDialect

//Manager dao manager interface
type Manager interface {
	//RegisterModel register model to manager
	RegisterModel(instanceGenerators ...GeneratorDefinition)
	//FindByID find data by id
	FindByID(modelName string, IDstring string, DB *gorm.DB, logEntry *logrus.Entry) (result interface{}, err error)
	//List query for list of data
	List(modelName string, q query.Q, pageSize int32, page int32, DB *gorm.DB, baseLogEntry *logrus.Entry) (listData interface{}, count int64, err error)

	//GetColumnName membaca nama column actual. dari catalog membaca nama actual column
	//syaratnya model sudah di register. kalau model belum di register ini akan throw error
	GetColumnName(modelName string, name string) (dbColumnName string, err error)

	//GenerateDaoWithWhere generate where. so reference is ready to use for query ( i.e. for count, list etc)
	GenerateDaoWithWhere(modelName string, q query.Q, DB *gorm.DB, baseLogEntry *logrus.Entry) (dbWithWhere *gorm.DB, err error)
	//SimpleManager generate simple dao manager
	SimpleManager(DB *gorm.DB, baseLogEntry *logrus.Entry) SimpleDaoManager
	//GetUsageCountQueries get usage count queries. this query will be use to determine is data allowed to be delete or not<br/>
	// example :
	// select   'UserRole' modelName ,  user_id   dataId , count(*) usageCount from sec_user_role where user_id in ( :ids )
	// pattern of queries :
	// result :
	// a. modelName : name of model that using refered data
	// b. dataId : id of data to count. aliased as dataId
	// c. usageCount int64 : count of usage of data
	// where
	// there should be string :ids. this pattern will be replaced with id of data to be count the usage count
	GetUsageCountQueries(dialect string, tableName string) (queries []string)

	//RegisterUsageCountQuery register usage count query to manager
	RegisterUsageCountQuery(dialect string, tableName string, query string)
}

//StructGeneratorFunction generator struct
type StructGeneratorFunction func() (instance interface{})

//StructSliceGeneratorFunction generator slice for struct
type StructSliceGeneratorFunction func() (instance interface{})

//GeneratorDefinition wrapper generator
type GeneratorDefinition struct {
	//Generator generator single instance
	Generator StructGeneratorFunction
	//SliceGenerator generator slice ( pointer)
	SliceGenerator StructSliceGeneratorFunction
}

//assignIDToStructFunction assign id ke dalam struct
type assignIDToStructFunction func(target interface{}, IDasString string) (err error)

type modelWorker struct {
	//columns  data columns. flat
	columns []extractGormResult
	//modelType tipe dari model
	modelType reflect.Type
	//generator generator instance struct
	generator StructGeneratorFunction
	//sliceGenerator generator slice
	sliceGenerator StructSliceGeneratorFunction
	//primaryKey data primary key
	primaryKey extractGormResult
	//columnNameIndexByStructFieldName nama column(database) di index dengan pendekatan :
	// 1. nama struct ( case sensitive)
	// 2. case in(insensitive)
	// 3. json tag( kalau ada)
	columnNameIndexByStructFieldName map[string]string
	//foreign keys
	foreignKeys []extractGormForeignKeyResult
}

//analizeModel extract model. dengan relection untuk memaca data
func analizeModel(structType reflect.Type) (result modelWorker) {
	primaryKeyColumn, columns, foreignKeys := analizeModelColumns(structType)
	columnNameIndexByStructFieldName := make(map[string]string)
	for _, col := range columns {
		columnNameIndexByStructFieldName[col.name] = col.columnName
		columnNameIndexByStructFieldName[strings.ToLower(col.name)] = col.columnName
		if len(col.jsonName) > 0 {
			columnNameIndexByStructFieldName[col.jsonName] = col.columnName
			columnNameIndexByStructFieldName[strings.ToLower(col.jsonName)] = col.columnName
		}
	}
	result = modelWorker{
		primaryKey:                       primaryKeyColumn,
		columns:                          columns,
		modelType:                        structType,
		foreignKeys:                      foreignKeys,
		columnNameIndexByStructFieldName: columnNameIndexByStructFieldName}
	return
}

//Manager dao utils untuk helper
type managerImplementation struct {
	//structGeneratorMap struct generator di index dengan model name
	structGeneratorMap map[string]modelWorker

	//usageCountQueries usage count queries
	usageCountQueries usageCountQuery
	//defaultUsageCountQueries if for example, postgres sql query was not defined, the default query will be returned for query of usage count on data
	defaultUsageCountQueries usageCountQueryByDialect
}

//RegisterModel register model
func (p *managerImplementation) RegisterModel(instanceGenerators ...GeneratorDefinition) {
	for _, smpl := range instanceGenerators {
		p.registerModelWorker(smpl.Generator, smpl.SliceGenerator)
	}
}

//GetUsageCountQueries get usage count queries. this query will be use to determine is data allowed to be delete or not
func (p *managerImplementation) GetUsageCountQueries(dialect string, tableName string) (queries []string) {
	if s, ok := p.usageCountQueries[dialect]; ok {
		return s[tableName]
	}
	return p.defaultUsageCountQueries[tableName]
}

//RegisterUsageCountQuery register usage count query to manager
func (p *managerImplementation) RegisterUsageCountQuery(dialect string, tableName string, query string) {
	if _, ok := p.usageCountQueries[dialect]; !ok {
		p.usageCountQueries[dialect] = make(usageCountQueryByDialect)
	}
	if _, ok2 := p.usageCountQueries[dialect][tableName]; !ok2 {
		p.usageCountQueries[dialect][tableName] = []string{query}
	} else {
		p.usageCountQueries[dialect][tableName] = append(p.usageCountQueries[dialect][tableName], query)
	}

}

//GetColumnName membaca nama column actual. dari catalog membaca nama actual column
//syaratnya model sudah di register. kalau model belum di register ini akan throw error
func (p *managerImplementation) GetColumnName(modelName string, name string) (dbColumnName string, err error) {
	dbColumnName = name
	rslt, ok := p.structGeneratorMap[modelName]
	if !ok {
		err = fmt.Errorf("Model %s not registered on Dao maanger, could not get actual colunm name for field %s", modelName, name)
		return
	}
	colActual, ok := rslt.columnNameIndexByStructFieldName[name]
	if ok {
		dbColumnName = colActual
	}
	return
}

//FindByID find data by id
func (p *managerImplementation) FindByID(modelName string, IDstring string, DB *gorm.DB, baseLogEntry *logrus.Entry) (result interface{}, err error) {
	var mdlInstc modelWorker
	var mdlFound = false

	logEntry := baseLogEntry.WithField("modelName", modelName)
	if mdlInstc, mdlFound = p.structGeneratorMap[modelName]; !mdlFound {
		msg := fmt.Sprintf("Model %s was not registered to manager. could not find by id for this model ", modelName)
		logEntry.Errorf(msg)
		err = fmt.Errorf(msg)
		return
	}
	if len(mdlInstc.primaryKey.columnName) == 0 {
		msgNoPK := fmt.Sprintf("Model %s does not have primary key. please check your model definition ", modelName)
		logEntry.Errorf(msgNoPK)
		err = fmt.Errorf(msgNoPK)
		return

	}
	instanceReturn := mdlInstc.generator()
	swapDB := DB.Where(fmt.Sprintf("%s = ?", mdlInstc.primaryKey.columnName), IDstring)
	rslt := swapDB.First(instanceReturn)
	if rslt.Error != nil {
		logEntry.WithError(rslt.Error).Errorf("Fail to query[%s] , erorr:%s", modelName, rslt.Error.Error())
		return nil, rslt.Error
	}
	if rslt.RowsAffected > 0 {
		return instanceReturn, nil
	}
	return nil, nil
}

//registerModelWorker actual worker untuk
func (p *managerImplementation) registerModelWorker(instanceGenerator StructGeneratorFunction, slicePointerGenerator StructSliceGeneratorFunction, modelAliases ...string) {
	sampleModel := instanceGenerator()
	structType := common.GetReflectTypeOfStructObject(sampleModel)
	fullName := structType.Name()
	var simpleName = fullName
	if strings.Contains(fullName, ".") {
		spl := strings.Split(simpleName, ".")
		simpleName = spl[len(spl)-1]
	}
	analizeResult := analizeModel(structType)
	analizeResult.generator = instanceGenerator
	analizeResult.sliceGenerator = slicePointerGenerator
	p.structGeneratorMap[fullName] = analizeResult
	p.structGeneratorMap[simpleName] = analizeResult
	for _, modelAlias := range modelAliases {
		p.structGeneratorMap[modelAlias] = analizeResult
	}

}

//SimpleManager generate simple dao manager
func (p *managerImplementation) SimpleManager(DB *gorm.DB, logEntry *logrus.Entry) SimpleDaoManager {
	return NewSimpleManager(p, DB, logEntry)
}

//List query for list of data
func (p *managerImplementation) List(modelName string, q query.Q, pageSize int32, page int32, DB *gorm.DB, baseLogEntry *logrus.Entry) (listData interface{}, count int64, err error) {
	var mdlInstc modelWorker
	var mdlFound = false
	logEntry := baseLogEntry.WithField("modelName", modelName).WithField("method", "List")
	if mdlInstc, mdlFound = p.structGeneratorMap[modelName]; !mdlFound {
		msg := fmt.Sprintf("Model %s was not registered to manager. could not find by id for this model ", modelName)
		logEntry.Errorf(msg)
		err = fmt.Errorf(msg)
		return
	}
	dbWithWhere, errToGenWhere := p.GenerateDaoWithWhere(modelName, q, DB, baseLogEntry)
	if errToGenWhere != nil {
		err = errToGenWhere
		return
	}
	listData = mdlInstc.sliceGenerator()
	dbWithWhere.Find(listData).Count(&count)

	return
}

//GenerateDaoWithWhere generate where. so reference is ready to use for query ( i.e. for count, list etc)
func (p *managerImplementation) GenerateDaoWithWhere(modelName string, q query.Q, DB *gorm.DB, baseLogEntry *logrus.Entry) (dbWithWhere *gorm.DB, err error) {

	var mdlFound = false
	logEntry := baseLogEntry.WithField("modelName", modelName).WithField("method", "GenerateDaoWithWhere")
	if _, mdlFound = p.structGeneratorMap[modelName]; !mdlFound {
		msg := fmt.Sprintf("Model %s was not registered to manager. could not find by id for this model ", modelName)
		logEntry.Errorf(msg)
		err = fmt.Errorf(msg)
		return
	}
	whereSQL, parameters, errGenQuery := q.GenerateSQL(modelName, p.GetColumnName)
	if errGenQuery != nil {
		logEntry.WithError(errGenQuery).Errorf("Failed to generate where sql, error: %s", errGenQuery.Error())
		err = errGenQuery
		return
	}
	dbWithWhere = DB.Where(whereSQL, parameters...)
	return
}

func (p *managerImplementation) ParseJSONQuery(JSONQuery string, startModel string) (parsedQuery query.Q) {
	JSONQuery = strings.TrimLeft(JSONQuery, " ")
	if OprOr.StartWith(JSONQuery) {

	} else {

	}
	return
}

//NewManager generate dao manager
func NewManager() Manager {
	return &managerImplementation{structGeneratorMap: make(map[string]modelWorker), usageCountQueries: make(usageCountQuery), defaultUsageCountQueries: make(usageCountQueryByDialect)}
}
