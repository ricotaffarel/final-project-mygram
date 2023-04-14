package models

import (
	"final-project-mygram/helpers"

	"github.com/asaskevich/govalidator"
	"gorm.io/gorm"
)

type User struct {
	GormModel
	Username string `gorm:"not null;uniqueIndex" json:"username" form:"username" valid:"required~Your full username is required"`
	Email    string `gorm:"not null;uniqueIndex" json:"email" form:"email" valid:"required~Your email is required, email~Invalid email format"`
	Password string `gorm:"not null" json:"password" form:"password" valid:"required~Your password is required,minstringlength(6)~Password has to have a minimun length of 6 characters"`
	Age      int    `gorm:"not null" json:"age" form:"age" valid:"required~Your age is required,range(8|100)~Age have minimum 8"`
	Role     string `gorm:"not null" json:"role" form:"role" valid:"required~Your role is required"`
}

type FormatterUser struct {
	GormModel
	Username string `gorm:"not null;uniqueIndex" json:"username" form:"username" valid:"required~Your full username is required"`
	Email    string `gorm:"not null;uniqueIndex" json:"email" form:"email" valid:"required~Your email is required, email~Invalid email format"`
	Age      int    `gorm:"not null" json:"age" form:"age" valid:"required~Your age is required,range(8|100)~Age have minimum 8"`
	Role     string `gorm:"not null" json:"role" form:"role" valid:"required~Your role is required"`
}

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	_, errCreate := govalidator.ValidateStruct(u)

	if errCreate != nil {
		err = errCreate
		return
	}

	u.Password = helpers.HassPass(u.Password)
	err = nil
	return
}
