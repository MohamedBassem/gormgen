package gormgen

import "github.com/jinzhu/gorm"

func (t *BasicModel) Save(db *gorm.DB) error {
	return db.Save(t).Error
}

func (t *BasicModel) Delete(db *gorm.DB) error {
	return db.Delete(t).Error
}

type BasicModelQueryBuilder struct {
	order []string
	where []struct {
		prefix string
		value  interface{}
	}
	limit  int
	offset int
}

func (qb *BasicModelQueryBuilder) buildQuery(db *gorm.DB) *gorm.DB {
	ret := db
	for _, where := range qb.where {
		ret = ret.Where(where)
	}
	for _, order := range qb.order {
		ret = ret.Order(order)
	}
	ret = ret.Limit(qb.limit).Offset(qb.offset)
	return ret
}

func (qb *BasicModelQueryBuilder) Count(db *gorm.DB) (int, error) {
	var c int
	res := qb.buildQuery(db).Model(&BasicModel{}).Count(&c)
	if res.RecordNotFound() {
		c = 0
	}
	return c, res.Error
}

func (qb *BasicModelQueryBuilder) First(db *gorm.DB) (*BasicModel, error) {
	ret := &BasicModel{}
	res := qb.buildQuery(db).First(ret)
	if res.RecordNotFound() {
		ret = nil
	}
	return ret, res.Error
}

func (qb *BasicModelQueryBuilder) QueryOne(db *gorm.DB) (*BasicModel, error) {
	qb.limit = 1
	ret, err := qb.QueryAll(db)
	if len(ret) > 0 {
		return &ret[0], err
	} else {
		return nil, err
	}
}

func (qb *BasicModelQueryBuilder) QueryAll(db *gorm.DB) ([]BasicModel, error) {
	ret := []BasicModel{}
	err := qb.buildQuery(db).Find(&ret).Error
	return ret, err
}

func (qb *BasicModelQueryBuilder) Limit(limit int) *BasicModelQueryBuilder {
	qb.limit = limit
	return qb
}

func (qb *BasicModelQueryBuilder) Offset(offset int) *BasicModelQueryBuilder {
	qb.offset = offset
	return qb
}

func (qb *BasicModelQueryBuilder) WhereID(p Predict, value int) *BasicModelQueryBuilder {
	qb.where = append(qb.where, struct {
		prefix string
		value  interface{}
	}{
		"id " + p.String(),
		value,
	})
	return qb
}

func (qb *BasicModelQueryBuilder) OrderByID(asc bool) *BasicModelQueryBuilder {
	order := "DESC"
	if asc {
		order = "ASC"
	}

	qb.order = append(qb.order, "id "+order)
	return qb
}

func (qb *BasicModelQueryBuilder) WhereName(p Predict, value string) *BasicModelQueryBuilder {
	qb.where = append(qb.where, struct {
		prefix string
		value  interface{}
	}{
		"name " + p.String(),
		value,
	})
	return qb
}

func (qb *BasicModelQueryBuilder) OrderByName(asc bool) *BasicModelQueryBuilder {
	order := "DESC"
	if asc {
		order = "ASC"
	}

	qb.order = append(qb.order, "name "+order)
	return qb
}

func (qb *BasicModelQueryBuilder) WhereAge(p Predict, value int) *BasicModelQueryBuilder {
	qb.where = append(qb.where, struct {
		prefix string
		value  interface{}
	}{
		"age " + p.String(),
		value,
	})
	return qb
}

func (qb *BasicModelQueryBuilder) OrderByAge(asc bool) *BasicModelQueryBuilder {
	order := "DESC"
	if asc {
		order = "ASC"
	}

	qb.order = append(qb.order, "age "+order)
	return qb
}
