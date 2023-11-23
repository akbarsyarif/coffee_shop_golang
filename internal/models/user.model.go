package models

import "time"

type UserModel struct {
	Id            string `db:"id" json:"id"`
	Full_name  string `db:"full_name" form:"full_name" json:"full_name"`
	Email  string `db:"email" form:"email" json:"email"`
	Pwd   string `db:"pwd" form:"pwd" json:"pwd"`
	Profile_pic   	  interface{} `db:"profile_pic" form:"profile_pic" json:"profile_pic"`
	Phone_number        interface{} `db:"phone_number" form:"phone_number" json:"phone_number"`
	Address      interface{} `db:"address" form:"address" json:"address"`
	Role_name         string `db:"role_name" form:"role_name" json:"role_name"`
	Isverified    bool `db:"isverified" form:"isverified" json:"isverified"`
	Created_at    *time.Time `db:"created_at" json:"created_at"`
	// Deleted_at    *time.Time `db:"deleted_at" json:"deleted_at"`
}