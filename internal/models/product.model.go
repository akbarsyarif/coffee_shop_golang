package models

import "time"

type ProductModel struct {
	Id            string `db:"id" json:"id"`
	Product_image interface{} `db:"product_image" form:"product_image" json:"product_image"`
	Product_name  string `db:"product_name" form:"product_name" json:"product_name"`
	Description   string `db:"description" form:"description" json:"description"`
	Rating   	  string `db:"rating" form:"rating" json:"rating"`
	Price         int `db:"price" form:"price" json:"price"`
	Category      string `db:"category" form:"category" json:"category"`
	Promo         string `db:"promo" form:"promo" json:"promo"`
	Created_at    *time.Time `db:"created_at" json:"created_at"`
	// Updated_at    *time.Time `db:"updated_at" json:"updated_at"`
	// Deleted_at    *time.Time `db:"deleted_at" json:"deleted_at"`
}