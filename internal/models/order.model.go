package models

import "time"

type OrderModel struct {
	Id         string      `db:"id" json:"id"`
	User_id         string      `db:"user_id" form:"user_id" json:"user_id"`
	Full_name  string      `db:"full_name" form:"full_name" json:"full_name"`
	Status     string      `db:"status_name" form:"status_name" json:"status_name"`
	Shipping   string      `db:"shipping_name" form:"shipping_name" json:"shipping_name"`
	Total      int         `db:"total" form:"total" json:"total"`
	Created_at *time.Time  `db:"created_at" json:"created_at"`
	Product []OrderProductModel `form:"product" json:"product"`
	// Deleted_at    *time.Time `db:"deleted_at" json:"deleted_at"`
}

type OrderProductModel struct {
	// Id         string      `db:"id" json:"id"`
	Order_id  string      `db:"order_id" form:"order_id" json:"order_id"`
	Full_name  string      `db:"full_name" form:"full_name" json:"full_name"`
	Product_image interface{} `db:"product_image" form:"product_image" json:"product_image"`
	Product_name  string `db:"product_name" form:"product_name" json:"product_name"`
	Quantity int `db:"quantity" form:"quantity" json:"quantity"`
	Size     string      `db:"size_name" form:"size_name" json:"size_name"`
	WithIce     bool     `db:"with_ice" form:"with_ice" json:"with_ice"`
	Status     string      `db:"status_name" form:"status_name" json:"status_name"`
	Sub_total      int      `db:"sub_total" form:"sub_total" json:"sub_total"`
	Shipping   string      `db:"shipping_name" form:"shipping_name" json:"shipping_name"`
	Created_at *time.Time  `db:"created_at" json:"created_at"`
	// Deleted_at    *time.Time `db:"deleted_at" json:"deleted_at"`
}

