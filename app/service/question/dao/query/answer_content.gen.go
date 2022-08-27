// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package query

import (
	"context"
	"main/app/service/question/dao/model"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/schema"

	"gorm.io/gen"
	"gorm.io/gen/field"

	"gorm.io/plugin/dbresolver"
)

func newAnswerContent(db *gorm.DB) answerContent {
	_answerContent := answerContent{}

	_answerContent.answerContentDo.UseDB(db)
	_answerContent.answerContentDo.UseModel(&model.AnswerContent{})

	tableName := _answerContent.answerContentDo.TableName()
	_answerContent.ALL = field.NewField(tableName, "*")
	_answerContent.AnswerID = field.NewInt64(tableName, "answer_id")
	_answerContent.Content = field.NewString(tableName, "content")
	_answerContent.IPLoc = field.NewString(tableName, "ip_loc")
	_answerContent.Meta = field.NewString(tableName, "meta")
	_answerContent.CreateTime = field.NewTime(tableName, "create_time")
	_answerContent.UpdateTime = field.NewTime(tableName, "update_time")

	_answerContent.fillFieldMap()

	return _answerContent
}

type answerContent struct {
	answerContentDo answerContentDo

	ALL        field.Field
	AnswerID   field.Int64
	Content    field.String
	IPLoc      field.String
	Meta       field.String
	CreateTime field.Time
	UpdateTime field.Time

	fieldMap map[string]field.Expr
}

func (a answerContent) Table(newTableName string) *answerContent {
	a.answerContentDo.UseTable(newTableName)
	return a.updateTableName(newTableName)
}

func (a answerContent) As(alias string) *answerContent {
	a.answerContentDo.DO = *(a.answerContentDo.As(alias).(*gen.DO))
	return a.updateTableName(alias)
}

func (a *answerContent) updateTableName(table string) *answerContent {
	a.ALL = field.NewField(table, "*")
	a.AnswerID = field.NewInt64(table, "answer_id")
	a.Content = field.NewString(table, "content")
	a.IPLoc = field.NewString(table, "ip_loc")
	a.Meta = field.NewString(table, "meta")
	a.CreateTime = field.NewTime(table, "create_time")
	a.UpdateTime = field.NewTime(table, "update_time")

	a.fillFieldMap()

	return a
}

func (a *answerContent) WithContext(ctx context.Context) *answerContentDo {
	return a.answerContentDo.WithContext(ctx)
}

func (a answerContent) TableName() string { return a.answerContentDo.TableName() }

func (a answerContent) Alias() string { return a.answerContentDo.Alias() }

func (a *answerContent) GetFieldByName(fieldName string) (field.OrderExpr, bool) {
	_f, ok := a.fieldMap[fieldName]
	if !ok || _f == nil {
		return nil, false
	}
	_oe, ok := _f.(field.OrderExpr)
	return _oe, ok
}

func (a *answerContent) fillFieldMap() {
	a.fieldMap = make(map[string]field.Expr, 6)
	a.fieldMap["answer_id"] = a.AnswerID
	a.fieldMap["content"] = a.Content
	a.fieldMap["ip_loc"] = a.IPLoc
	a.fieldMap["meta"] = a.Meta
	a.fieldMap["create_time"] = a.CreateTime
	a.fieldMap["update_time"] = a.UpdateTime
}

func (a answerContent) clone(db *gorm.DB) answerContent {
	a.answerContentDo.ReplaceDB(db)
	return a
}

type answerContentDo struct{ gen.DO }

func (a answerContentDo) Debug() *answerContentDo {
	return a.withDO(a.DO.Debug())
}

func (a answerContentDo) WithContext(ctx context.Context) *answerContentDo {
	return a.withDO(a.DO.WithContext(ctx))
}

func (a answerContentDo) ReadDB() *answerContentDo {
	return a.Clauses(dbresolver.Read)
}

func (a answerContentDo) WriteDB() *answerContentDo {
	return a.Clauses(dbresolver.Write)
}

func (a answerContentDo) Clauses(conds ...clause.Expression) *answerContentDo {
	return a.withDO(a.DO.Clauses(conds...))
}

func (a answerContentDo) Returning(value interface{}, columns ...string) *answerContentDo {
	return a.withDO(a.DO.Returning(value, columns...))
}

func (a answerContentDo) Not(conds ...gen.Condition) *answerContentDo {
	return a.withDO(a.DO.Not(conds...))
}

func (a answerContentDo) Or(conds ...gen.Condition) *answerContentDo {
	return a.withDO(a.DO.Or(conds...))
}

func (a answerContentDo) Select(conds ...field.Expr) *answerContentDo {
	return a.withDO(a.DO.Select(conds...))
}

func (a answerContentDo) Where(conds ...gen.Condition) *answerContentDo {
	return a.withDO(a.DO.Where(conds...))
}

func (a answerContentDo) Exists(subquery interface{ UnderlyingDB() *gorm.DB }) *answerContentDo {
	return a.Where(field.CompareSubQuery(field.ExistsOp, nil, subquery.UnderlyingDB()))
}

func (a answerContentDo) Order(conds ...field.Expr) *answerContentDo {
	return a.withDO(a.DO.Order(conds...))
}

func (a answerContentDo) Distinct(cols ...field.Expr) *answerContentDo {
	return a.withDO(a.DO.Distinct(cols...))
}

func (a answerContentDo) Omit(cols ...field.Expr) *answerContentDo {
	return a.withDO(a.DO.Omit(cols...))
}

func (a answerContentDo) Join(table schema.Tabler, on ...field.Expr) *answerContentDo {
	return a.withDO(a.DO.Join(table, on...))
}

func (a answerContentDo) LeftJoin(table schema.Tabler, on ...field.Expr) *answerContentDo {
	return a.withDO(a.DO.LeftJoin(table, on...))
}

func (a answerContentDo) RightJoin(table schema.Tabler, on ...field.Expr) *answerContentDo {
	return a.withDO(a.DO.RightJoin(table, on...))
}

func (a answerContentDo) Group(cols ...field.Expr) *answerContentDo {
	return a.withDO(a.DO.Group(cols...))
}

func (a answerContentDo) Having(conds ...gen.Condition) *answerContentDo {
	return a.withDO(a.DO.Having(conds...))
}

func (a answerContentDo) Limit(limit int) *answerContentDo {
	return a.withDO(a.DO.Limit(limit))
}

func (a answerContentDo) Offset(offset int) *answerContentDo {
	return a.withDO(a.DO.Offset(offset))
}

func (a answerContentDo) Scopes(funcs ...func(gen.Dao) gen.Dao) *answerContentDo {
	return a.withDO(a.DO.Scopes(funcs...))
}

func (a answerContentDo) Unscoped() *answerContentDo {
	return a.withDO(a.DO.Unscoped())
}

func (a answerContentDo) Create(values ...*model.AnswerContent) error {
	if len(values) == 0 {
		return nil
	}
	return a.DO.Create(values)
}

func (a answerContentDo) CreateInBatches(values []*model.AnswerContent, batchSize int) error {
	return a.DO.CreateInBatches(values, batchSize)
}

// Save : !!! underlying implementation is different with GORM
// The method is equivalent to executing the statement: db.Clauses(clause.OnConflict{UpdateAll: true}).Create(values)
func (a answerContentDo) Save(values ...*model.AnswerContent) error {
	if len(values) == 0 {
		return nil
	}
	return a.DO.Save(values)
}

func (a answerContentDo) First() (*model.AnswerContent, error) {
	if result, err := a.DO.First(); err != nil {
		return nil, err
	} else {
		return result.(*model.AnswerContent), nil
	}
}

func (a answerContentDo) Take() (*model.AnswerContent, error) {
	if result, err := a.DO.Take(); err != nil {
		return nil, err
	} else {
		return result.(*model.AnswerContent), nil
	}
}

func (a answerContentDo) Last() (*model.AnswerContent, error) {
	if result, err := a.DO.Last(); err != nil {
		return nil, err
	} else {
		return result.(*model.AnswerContent), nil
	}
}

func (a answerContentDo) Find() ([]*model.AnswerContent, error) {
	result, err := a.DO.Find()
	return result.([]*model.AnswerContent), err
}

func (a answerContentDo) FindInBatch(batchSize int, fc func(tx gen.Dao, batch int) error) (results []*model.AnswerContent, err error) {
	buf := make([]*model.AnswerContent, 0, batchSize)
	err = a.DO.FindInBatches(&buf, batchSize, func(tx gen.Dao, batch int) error {
		defer func() { results = append(results, buf...) }()
		return fc(tx, batch)
	})
	return results, err
}

func (a answerContentDo) FindInBatches(result *[]*model.AnswerContent, batchSize int, fc func(tx gen.Dao, batch int) error) error {
	return a.DO.FindInBatches(result, batchSize, fc)
}

func (a answerContentDo) Attrs(attrs ...field.AssignExpr) *answerContentDo {
	return a.withDO(a.DO.Attrs(attrs...))
}

func (a answerContentDo) Assign(attrs ...field.AssignExpr) *answerContentDo {
	return a.withDO(a.DO.Assign(attrs...))
}

func (a answerContentDo) Joins(fields ...field.RelationField) *answerContentDo {
	for _, _f := range fields {
		a = *a.withDO(a.DO.Joins(_f))
	}
	return &a
}

func (a answerContentDo) Preload(fields ...field.RelationField) *answerContentDo {
	for _, _f := range fields {
		a = *a.withDO(a.DO.Preload(_f))
	}
	return &a
}

func (a answerContentDo) FirstOrInit() (*model.AnswerContent, error) {
	if result, err := a.DO.FirstOrInit(); err != nil {
		return nil, err
	} else {
		return result.(*model.AnswerContent), nil
	}
}

func (a answerContentDo) FirstOrCreate() (*model.AnswerContent, error) {
	if result, err := a.DO.FirstOrCreate(); err != nil {
		return nil, err
	} else {
		return result.(*model.AnswerContent), nil
	}
}

func (a answerContentDo) FindByPage(offset int, limit int) (result []*model.AnswerContent, count int64, err error) {
	result, err = a.Offset(offset).Limit(limit).Find()
	if err != nil {
		return
	}

	if size := len(result); 0 < limit && 0 < size && size < limit {
		count = int64(size + offset)
		return
	}

	count, err = a.Offset(-1).Limit(-1).Count()
	return
}

func (a answerContentDo) ScanByPage(result interface{}, offset int, limit int) (count int64, err error) {
	count, err = a.Count()
	if err != nil {
		return
	}

	err = a.Offset(offset).Limit(limit).Scan(result)
	return
}

func (a answerContentDo) Scan(result interface{}) (err error) {
	return a.DO.Scan(result)
}

func (a *answerContentDo) withDO(do gen.Dao) *answerContentDo {
	a.DO = *do.(*gen.DO)
	return a
}
