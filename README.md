# gormgen

[![Build Status](https://travis-ci.org/MohamedBassem/gormgen.svg?branch=master)](https://travis-ci.org/MohamedBassem/gormgen)

gormgen is a code generation tool to generate a better API to query and update [gorm](https://github.com/jinzhu/gorm) structs without having to deal with `interface{}`s or with database column names.

**Note** : gormgen is still is still in early development phase. It may contain bugs and the API is not yet stable. Your suggestions for improving gormgen are welcome through issues/PRs.

## Why to use gormgen

```go

// Querying

// The gorm way:
users := []User{}
err := db.Where("age > ?", 20).Order("age ASC").Limit(10).Find(&users).Error

// gormgen way
users, err := (&UserQueryBuilder{}).
  WhereAge(gormgen.GreaterThanPredicate, 20).
  OrderByAge(true).
  Limit(10).
  QueryAll(db)


// Creating Object
user := &User{
  Name: "Bla",
  Age: 20,
}

// The gorm way
err := db.Create(user).Error

// The gormgen way
err := user.Save(db)
```

- No more ugly `interface{}`s when doing in the `Where` function. Using gormgen, the passed values will be type checked.
- No more ugly strings for column names for `Where` and `Order` functions. By this, you won't need to convert the field name to the column name yourself, gormgen will do it for you. Also, you won't forget to change a column name when you change the field name because your code won't compile until you fix it everywhere.
- A more intuitive way to return the results instead of passing them as a param.
- It doesn't alter your struct, so it's still compatible with gorm and you can still use the gorm way whenever you want (or for missing features in gormgen).

## How it works

If you have the following :

```go
//go:generate gormgen -structs User -output user_gen.go
type User struct {
	ID   uint   `gorm:"primary_key"`
	Name string
	Age  int
}
```

Run `go generate` and gormgen will generate for you :

```go
func (t *User) Save(db *gorm.DB) error {/* … */}
func (t *User) Delete(db *gorm.DB) error {/* … */}
type UserQueryBuilder struct {/* … */}
func (qb *UserQueryBuilder) Count(db *gorm.DB) (int, error) {/* … */}
func (qb *UserQueryBuilder) First(db *gorm.DB) (*User, error) {/* … */} // Sorted by primary key
func (qb *UserQueryBuilder) QueryOne(db *gorm.DB) (*User, error) {/* … */} // Sorted by the order specified
func (qb *UserQueryBuilder) QueryAll(db *gorm.DB) ([]User, error) {/* … */}
func (qb *UserQueryBuilder) Limit(limit int) *UserQueryBuilder {/* … */}
func (qb *UserQueryBuilder) Offset(offset int) *UserQueryBuilder {/* … */}
func (qb *UserQueryBuilder) WhereID(p gormgen.Predicate, value uint) *UserQueryBuilder {/* … */}
func (qb *UserQueryBuilder) OrderByID(asc bool) *UserQueryBuilder {/* … */}
func (qb *UserQueryBuilder) WhereName(p gormgen.Predicate, value string) *UserQueryBuilder {/* … */}
func (qb *UserQueryBuilder) OrderByName(asc bool) *UserQueryBuilder {/* … */}
func (qb *UserQueryBuilder) WhereAge(p gormgen.Predicate, value int) *UserQueryBuilder {/* … */}
func (qb *UserQueryBuilder) OrderByAge(asc bool) *UserQueryBuilder {/* … */}
```

For the actual generated code, check the examples folder.

## How to use it

- `go get -u github.com/MohamedBassem/gormgen/...`
- Add the `//go:generate` comment mentioned above anywhere in your code.
- Add `go generate` to your build steps.
- **The generated code will depend on gorm and gormgen, so make sure to vendor both of them.**

## Not yet supported features

- [X] Inferring database column name from gorm convention or gorm struct tag.
- [X] Ignoring fields with `gorm:"-"`.
- [X] Support for anonymous structs (IMPORTANT for gorm.Model).
- [ ] Support for type aliases.
- [ ] Support for detecting and querying with primary key.
- [ ] Support for the embedded struct tag.

## Contributing

Your contributions and ideas are welcomed through issues and pull requests.

*Note for development* : Make sure to have `gormgen` in your path to be able to run the tests. Also, always run the tests with `make test` to regenerate the test structs.

## Note

The parser of this package is heavily inspired from the source code of `https://godoc.org/golang.org/x/tools/cmd/stringer`. That's where I learned how to parse and type check a go package.
