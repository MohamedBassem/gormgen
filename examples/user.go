package example

import "github.com/jinzhu/gorm"

//go:generate gormgen -structs User -output user_gen.go
type User struct {
	gorm.Model
	Name  string
	Age   int
	Email string
}
