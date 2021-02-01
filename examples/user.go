package example

import "gorm.io/gorm"

//go:generate gormgen -structs User -output user_gen.go
type User struct {
	gorm.Model
	Name  string
	Age   int
	Email string
}
