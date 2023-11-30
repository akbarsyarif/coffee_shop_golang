package models

import (
	"time"
)

type UserModel struct {
	Id            string `db:"id" json:"id,omitempty" valid:"-"`
	Full_name  string `db:"full_name" form:"full_name" json:"full_name,omitempty" valid:"required,minstringlength(6)"`
	Email  string `db:"email" form:"email" json:"email,omitempty" valid:"required,email"`
	Pwd   string `db:"pwd" form:"pwd" json:"pwd,omitempty" valid:"required,minstringlength(4)"`
	Profile_pic   	  interface{} `db:"profile_pic" form:"profile_pic" json:"profile_pic" valid:"-"`
	Phone_number        string `db:"phone_number" form:"phone_number" json:"phone_number,omitempty" valid:"-"`
	Address      string `db:"address" uri:"address" form:"address" json:"address,omitempty" valid:"-"`
	Role_name         string `db:"role_name" form:"role_name" json:"role_name,omitempty" valid:"-"`
	Isverified    bool `db:"isverified" form:"isverified" json:"isverified,omitempty" valid:"-"`
	Created_at    *time.Time `db:"created_at" json:"created_at,omitempty" valid:"-"`
	// Deleted_at    *time.Time `db:"deleted_at" json:"deleted_at"`
}

type UpdateUserModel struct {
	Id            string `db:"id" json:"id" valid:"-"`
	Full_name  string `db:"full_name" form:"full_name" json:"full_name" valid:"optional"`
	Email  string `db:"email" form:"email" json:"email" valid:"-"`
	Profile_pic   	  interface{} `db:"profile_pic" uri:"profile_pic" form:"profile_pic" json:"profile_pic" valid:"-"`
	Phone_number        string `db:"phone_number" form:"phone_number" json:"phone_number" valid:"-"`
	Address      string `db:"address" uri:"address" form:"address" json:"address" valid:"-"`
	Role_name         string `db:"role_name" form:"role_name" json:"role_name" valid:"-"`
	Isverified    bool `db:"isverified" form:"isverified" json:"isverified" valid:"-"`
	Created_at    *time.Time `db:"created_at" json:"created_at" valid:"-"`
}

type GetUserInfoModel struct {
	Id       int    `db:"id" json:"id" valid:"-"`
	Full_name string `db:"full_name" form:"full_name" json:"full_name" valid:"-"`
	Email    string `db:"email" form:"email" json:"email" valid:"required,email"`
	Pwd string `db:"pwd" form:"pwd" json:"pwd" valid:"required"`
	Role_name     string `db:"role_name" json:"role_name" valid:"-"`
	Isverified    bool `db:"isverified" form:"isverified" json:"isverified" valid:"-"`

}

type BlackListToken struct {
	Blacklist_token  string  `db:"blacklist_token" json:"blacklist_token" valid:"-"`
}