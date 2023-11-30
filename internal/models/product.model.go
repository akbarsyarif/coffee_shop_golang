package models

import "time"

type ProductModel struct {
	Id            string `db:"id" json:"id,omitempty" valid:"-"`
	Product_image interface{} `db:"product_image" form:"product_image" json:"product_image" valid:"-"`
	Product_name  string `db:"product_name" form:"product_name" json:"product_name,omitempty" valid:"-"`
	Description   string `db:"description" form:"description" json:"description,omitempty" valid:"-"`
	Rating   	  string `db:"rating" form:"rating" json:"rating,omitempty" valid:"-"`
	Price         int `db:"price" form:"price" json:"price,omitempty" valid:"-"`
	Category      string `db:"category" form:"category" json:"category" valid:"-"`
	Promo         string `db:"promo" form:"promo" json:"promo" valid:"-"`
	Created_at    *time.Time `db:"created_at" json:"created_at,omitempty" valid:"-"`
	// Updated_at    *time.Time `db:"updated_at" json:"updated_at"`
	// Deleted_at    *time.Time `db:"deleted_at" json:"deleted_at"`
}

type ProductParams struct {
	Product_name     string `form:"search" json:"search" valid:"optional"`
	Max_price    string `form:"max_price" json:"max_price" valid:"optional"`
	Min_price    string `form:"min_price" json:"min_price" valid:"optional"`
	Category string `form:"category" json:"category" valid:"optional"`
	Name        string `form:"product_name" json:"product_name" valid:"optional, in(asc|desc)"`
	Price       string `form:"price" json:"price" valid:"optional, in(asc|desc)"`
	Created_at    string `form:"created_at" json:"created_at" valid:"optional, in(asc|desc)"`
	Page        string `form:"page" json:"page" valid:"optional"`
}