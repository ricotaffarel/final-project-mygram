package models

import (
	"github.com/asaskevich/govalidator"
	"gorm.io/gorm"
)

type Comment struct {
	GormModel
	Messege string `json:"messege" form:"messege" valid:"required~Messege of your messege is required"`
	UserID  uint   `json:"user_id" form:"user_id" valid:"required~UserID of your messege is required"`
	PhotoID uint   `json:"photo_id" form:"photo_id" valid:"required~PhotoID of your messege is required"`
	// User    *User
	Photo   *Photo
}

func (p *Comment) BeforeCreate(tx *gorm.DB) (err error) {
	_, errCreate := govalidator.ValidateStruct(p)

	if errCreate != nil {
		err = errCreate
		return
	}

	err = nil
	return
}

func (p *Comment) BeforeUpdate(tx *gorm.DB) (err error) {
	_, errCreate := govalidator.ValidateStruct(p)

	if errCreate != nil {
		err = errCreate
		return
	}

	err = nil
	return
}
